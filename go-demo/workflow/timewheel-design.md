# 时间轮架构设计方案

## 概述

本文档描述了一个用于工作流引擎的高性能多重时间轮定时任务调度系统。该系统采用层级时间轮架构，支持大规模定时任务的高效调度，满足工作流引擎中任务超时、延迟执行、周期调度等场景需求。

---

## 一、时间轮基础原理

### 1.1 什么是时间轮

时间轮（Timing Wheel）是一种高效的定时器实现方式，类似于时钟的表盘：

```
        12 (槽位0)
    11      1
  10          2
 9      ●      3  <- 指针
  8          4
    7      5
        6
```

**核心思想**：
- 将时间划分为固定大小的槽位（Slot/Bucket）
- 定时任务根据延迟时间放入对应槽位
- 指针按固定间隔（Tick）推进，执行到期任务
- 相同延迟时间的任务存储在同一槽位（链表/数组）

**时间复杂度**：
- 添加任务：O(1)
- 删除任务：O(1)（需要维护任务引用）
- 执行任务：O(1)均摊复杂度

### 1.2 为什么需要多重时间轮

单层时间轮存在局限：
- **时间范围有限**：槽位数 × tick间隔 = 最大延迟时间
- **精度与范围矛盾**：高精度需要小tick，长时间需要多槽位
- **空间浪费**：支持1小时延迟，秒级精度需要3600个槽位

**多重时间轮方案**：
```
毫秒轮(1ms×1000) -> 秒轮(1s×60) -> 分轮(1min×60) -> 时轮(1h×24) -> 天轮(1d×30)
     0-999ms         0-59s          0-59min       0-23h         0-29d
```

类似于时钟的：毫秒针 → 秒针 → 分针 → 时针 → 日历

---

## 二、核心数据结构设计

### 2.1 时间轮定义

```go
// TimeWheel 单层时间轮
type TimeWheel struct {
    // 基础属性
    tick          time.Duration    // 每个槽位代表的时间跨度
    wheelSize     int              // 槽位数量
    interval      time.Duration    // 整个轮的时间范围 = tick * wheelSize
    currentTick   int64            // 当前指针位置
    startTime     time.Time        // 时间轮启动时间
    
    // 槽位存储
    slots         []TaskList       // 每个槽位的任务列表
    
    // 层级关系
    level         int              // 当前层级（0=最高精度）
    lowerWheel    *TimeWheel       // 下一级（更高精度）时间轮
    upperWheel    *TimeWheel       // 上一级（更低精度）时间轮
    
    // 并发控制
    mu            sync.RWMutex     // 保护槽位操作
    
    // 控制通道
    ticker        *time.Ticker     // 驱动时间轮转动
    stopChan      chan struct{}    // 停止信号
    
    // 统计信息
    taskCount     atomic.Int64     // 当前任务总数
    executedCount atomic.Int64     // 已执行任务数
}

// TaskList 任务列表（使用链表实现）
type TaskList struct {
    head *TaskNode
    tail *TaskNode
    size int
    mu   sync.Mutex
}

// TaskNode 任务节点
type TaskNode struct {
    task *Task
    next *TaskNode
    prev *TaskNode
}
```

### 2.2 任务定义

```go
// Task 定时任务
type Task struct {
    // 任务标识
    ID            string           // 唯一标识
    Name          string           // 任务名称
    
    // 时间属性
    Delay         time.Duration    // 延迟时间
    ExecuteAt     time.Time        // 精确执行时间
    CreatedAt     time.Time        // 创建时间
    
    // 执行属性
    Callback      TaskCallback     // 执行回调
    Context       context.Context  // 执行上下文
    Retries       int              // 重试次数
    MaxRetries    int              // 最大重试次数
    
    // 循环任务属性
    Type          TaskType         // 一次性/周期性
    Interval      time.Duration    // 周期间隔
    CronExpr      string           // Cron表达式（可选）
    RepeatCount   int              // 重复次数（-1表示无限）
    ExecutedCount int              // 已执行次数
    
    // 状态管理
    Status        TaskStatus       // 任务状态
    Round         int              // 圈数（多圈时间轮用）
    SlotIndex     int              // 所在槽位索引
    Level         int              // 所在时间轮层级
    
    // 回调结果
    Result        interface{}      // 执行结果
    Error         error            // 执行错误
    
    // 引用管理
    cancelFunc    context.CancelFunc
    node          *TaskNode        // 在链表中的节点引用
}

// TaskCallback 任务回调函数
type TaskCallback func(ctx context.Context, task *Task) error

// TaskType 任务类型
type TaskType int

const (
    TaskTypeOnce     TaskType = iota  // 一次性任务
    TaskTypePeriodic                  // 周期性任务
    TaskTypeCron                      // Cron表达式任务
)

// TaskStatus 任务状态
type TaskStatus int

const (
    TaskStatusPending   TaskStatus = iota  // 等待执行
    TaskStatusRunning                      // 执行中
    TaskStatusCompleted                    // 已完成
    TaskStatusFailed                       // 执行失败
    TaskStatusCanceled                     // 已取消
)
```

### 2.3 多重时间轮管理器

```go
// TimingWheelManager 多重时间轮管理器
type TimingWheelManager struct {
    // 时间轮层级
    wheels        []*TimeWheel     // 多级时间轮数组
    wheelConfigs  []WheelConfig    // 各层级配置
    
    // 任务索引（用于快速查找和取消）
    taskIndex     sync.Map         // map[string]*Task
    
    // 工作协程池
    workerPool    *WorkerPool      // 任务执行协程池
    
    // 延迟队列（降级方案）
    delayQueue    *DelayQueue      // 超长延迟任务队列
    
    // 持久化
    storage       Storage          // 持久化存储接口
    enablePersist bool             // 是否启用持久化
    
    // 监控统计
    metrics       *Metrics         // 监控指标
    
    // 控制
    running       atomic.Bool      // 运行状态
    wg            sync.WaitGroup   // 等待所有goroutine结束
}

// WheelConfig 时间轮配置
type WheelConfig struct {
    Tick      time.Duration   // 槽位时间跨度
    SlotCount int             // 槽位数量
    Level     int             // 层级
}

// 默认五级时间轮配置
var DefaultWheelConfigs = []WheelConfig{
    {Tick: 10 * time.Millisecond, SlotCount: 100, Level: 0},  // 0-1秒，10ms精度
    {Tick: 1 * time.Second, SlotCount: 60, Level: 1},         // 0-60秒
    {Tick: 1 * time.Minute, SlotCount: 60, Level: 2},         // 0-60分钟
    {Tick: 1 * time.Hour, SlotCount: 24, Level: 3},           // 0-24小时
    {Tick: 24 * time.Hour, SlotCount: 30, Level: 4},          // 0-30天
}
```

---

## 三、核心算法实现

### 3.1 任务添加算法

```go
// AddTask 添加定时任务
func (m *TimingWheelManager) AddTask(task *Task) error {
    // 1. 参数校验
    if task.Delay < 0 {
        return errors.New("invalid delay time")
    }
    
    // 2. 计算执行时间
    now := time.Now()
    task.ExecuteAt = now.Add(task.Delay)
    task.CreatedAt = now
    task.Status = TaskStatusPending
    
    // 3. 任务分配到合适的时间轮
    wheel, slotIndex, round := m.findWheel(task.Delay)
    if wheel == nil {
        // 延迟时间超出时间轮范围，使用延迟队列
        return m.delayQueue.Push(task)
    }
    
    // 4. 设置任务属性
    task.Level = wheel.level
    task.SlotIndex = slotIndex
    task.Round = round
    
    // 5. 添加到槽位
    wheel.addToSlot(slotIndex, task)
    
    // 6. 建立索引
    m.taskIndex.Store(task.ID, task)
    
    // 7. 持久化（可选）
    if m.enablePersist {
        m.storage.SaveTask(task)
    }
    
    // 8. 统计
    m.metrics.TaskAdded.Inc()
    
    return nil
}

// findWheel 查找合适的时间轮和槽位
func (m *TimingWheelManager) findWheel(delay time.Duration) (*TimeWheel, int, int) {
    // 从最高精度轮开始查找
    for _, wheel := range m.wheels {
        if delay < wheel.interval {
            // 计算相对当前tick的偏移
            offset := int64(delay / wheel.tick)
            
            // 计算槽位索引和圈数
            slotIndex := (wheel.currentTick + offset) % int64(wheel.wheelSize)
            round := int(offset / int64(wheel.wheelSize))
            
            return wheel, int(slotIndex), round
        }
    }
    
    // 超出最大时间轮范围
    return nil, 0, 0
}

// addToSlot 添加任务到槽位
func (w *TimeWheel) addToSlot(slotIndex int, task *Task) {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    slot := &w.slots[slotIndex]
    
    // 创建任务节点
    node := &TaskNode{task: task}
    task.node = node // 保存引用，便于删除
    
    // 添加到链表尾部
    slot.mu.Lock()
    defer slot.mu.Unlock()
    
    if slot.tail == nil {
        slot.head = node
        slot.tail = node
    } else {
        slot.tail.next = node
        node.prev = slot.tail
        slot.tail = node
    }
    
    slot.size++
    w.taskCount.Add(1)
}
```

### 3.2 时间轮推进算法

```go
// Start 启动时间轮
func (w *TimeWheel) Start() {
    w.ticker = time.NewTicker(w.tick)
    w.startTime = time.Now()
    
    go func() {
        for {
            select {
            case <-w.ticker.C:
                w.advance()
            case <-w.stopChan:
                return
            }
        }
    }()
}

// advance 推进时间轮
func (w *TimeWheel) advance() {
    w.mu.Lock()
    
    // 获取当前槽位
    slotIndex := int(w.currentTick % int64(w.wheelSize))
    slot := &w.slots[slotIndex]
    
    // 推进指针
    w.currentTick++
    
    // 检查是否完整转了一圈
    if w.currentTick%int64(w.wheelSize) == 0 && w.upperWheel != nil {
        // 触发上级时间轮降级
        w.upperWheel.cascade()
    }
    
    w.mu.Unlock()
    
    // 执行当前槽位的任务
    w.executeSlot(slot)
}

// executeSlot 执行槽位中的任务
func (w *TimeWheel) executeSlot(slot *TaskList) {
    slot.mu.Lock()
    
    // 取出所有任务
    current := slot.head
    slot.head = nil
    slot.tail = nil
    size := slot.size
    slot.size = 0
    
    slot.mu.Unlock()
    
    // 遍历任务链表
    for current != nil {
        task := current.task
        next := current.next
        
        if task.Round > 0 {
            // 还需要继续等待
            task.Round--
            w.addToSlot(int(task.SlotIndex), task)
        } else {
            // 到期任务
            if task.Level > 0 && w.lowerWheel != nil {
                // 降级到下一级时间轮
                w.cascadeTask(task)
            } else {
                // 最高精度层级，执行任务
                w.executeTask(task)
            }
        }
        
        current = next
    }
    
    w.taskCount.Add(-int64(size))
}

// cascadeTask 任务降级到下级时间轮
func (w *TimeWheel) cascadeTask(task *Task) {
    if w.lowerWheel == nil {
        w.executeTask(task)
        return
    }
    
    // 计算在下级时间轮中的位置
    remainDelay := task.ExecuteAt.Sub(time.Now())
    if remainDelay <= 0 {
        w.executeTask(task)
        return
    }
    
    offset := int64(remainDelay / w.lowerWheel.tick)
    slotIndex := (w.lowerWheel.currentTick + offset) % int64(w.lowerWheel.wheelSize)
    round := int(offset / int64(w.lowerWheel.wheelSize))
    
    task.Level = w.lowerWheel.level
    task.SlotIndex = int(slotIndex)
    task.Round = round
    
    w.lowerWheel.addToSlot(int(slotIndex), task)
}

// executeTask 执行任务
func (w *TimeWheel) executeTask(task *Task) {
    task.Status = TaskStatusRunning
    
    // 提交到工作协程池执行
    w.workerPool.Submit(func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Task panic: %v, stack: %s", r, debug.Stack())
                task.Status = TaskStatusFailed
                task.Error = fmt.Errorf("panic: %v", r)
            }
        }()
        
        // 执行任务回调
        ctx, cancel := context.WithTimeout(task.Context, 5*time.Minute)
        defer cancel()
        
        err := task.Callback(ctx, task)
        
        if err != nil {
            task.Error = err
            task.Status = TaskStatusFailed
            
            // 重试逻辑
            if task.Retries < task.MaxRetries {
                task.Retries++
                // 重新添加到时间轮
                task.Delay = time.Duration(task.Retries) * time.Second * 2 // 指数退避
                w.manager.AddTask(task)
                return
            }
        } else {
            task.Status = TaskStatusCompleted
        }
        
        // 周期性任务重新调度
        if task.Type == TaskTypePeriodic && task.Status == TaskStatusCompleted {
            task.ExecutedCount++
            
            if task.RepeatCount == -1 || task.ExecutedCount < task.RepeatCount {
                task.Delay = task.Interval
                task.Status = TaskStatusPending
                task.Retries = 0
                w.manager.AddTask(task)
            }
        }
        
        // 统计
        w.executedCount.Add(1)
    })
}
```

### 3.3 任务取消算法

```go
// CancelTask 取消任务
func (m *TimingWheelManager) CancelTask(taskID string) error {
    // 1. 从索引中查找任务
    value, ok := m.taskIndex.Load(taskID)
    if !ok {
        return errors.New("task not found")
    }
    
    task := value.(*Task)
    
    // 2. 标记为取消状态
    if task.Status == TaskStatusRunning {
        // 任务正在执行，取消上下文
        if task.cancelFunc != nil {
            task.cancelFunc()
        }
        return errors.New("task is running, context canceled")
    }
    
    task.Status = TaskStatusCanceled
    
    // 3. 从时间轮中移除
    wheel := m.wheels[task.Level]
    wheel.removeFromSlot(task.SlotIndex, task)
    
    // 4. 从索引中删除
    m.taskIndex.Delete(taskID)
    
    // 5. 持久化删除
    if m.enablePersist {
        m.storage.DeleteTask(taskID)
    }
    
    // 6. 统计
    m.metrics.TaskCanceled.Inc()
    
    return nil
}

// removeFromSlot 从槽位中移除任务
func (w *TimeWheel) removeFromSlot(slotIndex int, task *Task) {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    slot := &w.slots[slotIndex]
    slot.mu.Lock()
    defer slot.mu.Unlock()
    
    node := task.node
    if node == nil {
        return
    }
    
    // 从双向链表中删除
    if node.prev != nil {
        node.prev.next = node.next
    } else {
        slot.head = node.next
    }
    
    if node.next != nil {
        node.next.prev = node.prev
    } else {
        slot.tail = node.prev
    }
    
    slot.size--
    w.taskCount.Add(-1)
}
```

---

## 四、工作协程池设计

```go
// WorkerPool 工作协程池
type WorkerPool struct {
    workerCount int
    taskQueue   chan func()
    wg          sync.WaitGroup
    stopChan    chan struct{}
}

// NewWorkerPool 创建协程池
func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
    pool := &WorkerPool{
        workerCount: workerCount,
        taskQueue:   make(chan func(), queueSize),
        stopChan:    make(chan struct{}),
    }
    
    // 启动工作协程
    for i := 0; i < workerCount; i++ {
        pool.wg.Add(1)
        go pool.worker()
    }
    
    return pool
}

// worker 工作协程
func (p *WorkerPool) worker() {
    defer p.wg.Done()
    
    for {
        select {
        case task := <-p.taskQueue:
            if task != nil {
                task()
            }
        case <-p.stopChan:
            return
        }
    }
}

// Submit 提交任务
func (p *WorkerPool) Submit(task func()) error {
    select {
    case p.taskQueue <- task:
        return nil
    case <-time.After(5 * time.Second):
        return errors.New("worker pool is full")
    }
}

// Stop 停止协程池
func (p *WorkerPool) Stop() {
    close(p.stopChan)
    p.wg.Wait()
}
```

---

## 五、延迟队列（降级方案）

对于超出时间轮范围的超长延迟任务，使用优先队列：

```go
// DelayQueue 延迟队列（小顶堆实现）
type DelayQueue struct {
    heap     *TaskHeap
    mu       sync.Mutex
    notEmpty *sync.Cond
    stopChan chan struct{}
    manager  *TimingWheelManager
}

// TaskHeap 任务堆
type TaskHeap []*Task

func (h TaskHeap) Len() int           { return len(h) }
func (h TaskHeap) Less(i, j int) bool { return h[i].ExecuteAt.Before(h[j].ExecuteAt) }
func (h TaskHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TaskHeap) Push(x interface{}) {
    *h = append(*h, x.(*Task))
}

func (h *TaskHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

// Start 启动延迟队列
func (q *DelayQueue) Start() {
    go func() {
        for {
            select {
            case <-q.stopChan:
                return
            default:
                q.process()
            }
        }
    }()
}

// process 处理队列
func (q *DelayQueue) process() {
    q.mu.Lock()
    
    for q.heap.Len() == 0 {
        q.notEmpty.Wait()
    }
    
    // 查看堆顶任务
    task := (*q.heap)[0]
    now := time.Now()
    
    if task.ExecuteAt.After(now) {
        // 还未到期，等待
        waitDuration := task.ExecuteAt.Sub(now)
        q.mu.Unlock()
        
        timer := time.NewTimer(waitDuration)
        select {
        case <-timer.C:
        case <-q.stopChan:
            timer.Stop()
            return
        }
        
        return
    }
    
    // 任务到期，弹出
    heap.Pop(q.heap)
    q.mu.Unlock()
    
    // 重新添加到时间轮
    task.Delay = 0
    q.manager.AddTask(task)
}

// Push 添加任务到延迟队列
func (q *DelayQueue) Push(task *Task) error {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    heap.Push(q.heap, task)
    q.notEmpty.Signal()
    
    return nil
}
```

---

## 六、与工作流引擎集成

### 6.1 工作流超时检测

```go
// TimeoutMonitor 超时监控器
type TimeoutMonitor struct {
    timingWheel *TimingWheelManager
    engine      *ExecutionEngine
}

// MonitorTask 监控任务超时
func (tm *TimeoutMonitor) MonitorTask(task *WorkflowTask, timeout time.Duration) error {
    // 创建超时检查任务
    timeoutTask := &Task{
        ID:    fmt.Sprintf("timeout_%s", task.ID),
        Name:  "Task Timeout Check",
        Delay: timeout,
        Type:  TaskTypeOnce,
        Callback: func(ctx context.Context, t *Task) error {
            // 检查工作流任务是否仍在执行
            workflowTask, err := tm.engine.taskRepo.GetByID(task.ID)
            if err != nil {
                return err
            }
            
            if workflowTask.Status == TaskPending {
                // 任务超时，执行超时处理
                return tm.handleTimeout(workflowTask)
            }
            
            return nil
        },
    }
    
    return tm.timingWheel.AddTask(timeoutTask)
}

// handleTimeout 处理超时
func (tm *TimeoutMonitor) handleTimeout(task *WorkflowTask) error {
    instance, _ := tm.engine.instanceRepo.GetByID(task.InstanceID)
    workflow, _ := tm.engine.workflowRepo.GetByID(instance.WorkflowID)
    node := workflow.FindNode(task.NodeID)
    
    if node.TimeoutRule == nil {
        return nil
    }
    
    switch node.TimeoutRule.Action {
    case "Auto":
        // 自动通过
        return tm.engine.CompleteTask(task.ID, "system", "auto_approve", "超时自动通过", nil)
        
    case "Escalate":
        // 升级处理
        return tm.escalateTask(task, node.TimeoutRule.Target)
        
    case "Notify":
        // 发送通知
        return tm.sendNotification(task)
        
    default:
        return nil
    }
}
```

### 6.2 延迟任务执行

```go
// ScheduleDelayedTask 调度延迟任务
func (e *ExecutionEngine) ScheduleDelayedTask(
    workflowID string,
    documentID string,
    delay time.Duration,
    variables map[string]interface{},
) error {
    
    task := &Task{
        ID:    generateID(),
        Name:  "Delayed Workflow Start",
        Delay: delay,
        Type:  TaskTypeOnce,
        Callback: func(ctx context.Context, t *Task) error {
            // 启动工作流
            _, err := e.StartProcess(workflowID, documentID, "system", variables)
            return err
        },
    }
    
    return e.timingWheel.AddTask(task)
}
```

### 6.3 周期性任务

```go
// SchedulePeriodicTask 调度周期性任务
func (e *ExecutionEngine) SchedulePeriodicTask(
    name string,
    interval time.Duration,
    callback TaskCallback,
) error {
    
    task := &Task{
        ID:       generateID(),
        Name:     name,
        Delay:    interval,
        Type:     TaskTypePeriodic,
        Interval: interval,
        Callback: callback,
        RepeatCount: -1, // 无限循环
    }
    
    return e.timingWheel.AddTask(task)
}

// 使用示例：定期清理过期数据
func setupPeriodicCleanup(engine *ExecutionEngine) {
    engine.SchedulePeriodicTask(
        "cleanup_expired_instances",
        24*time.Hour,
        func(ctx context.Context, task *Task) error {
            // 清理30天前的已完成流程实例
            cutoffTime := time.Now().Add(-30 * 24 * time.Hour)
            return engine.instanceRepo.DeleteBefore(cutoffTime)
        },
    )
}
```

---

## 七、持久化设计

### 7.1 持久化接口

```go
// Storage 持久化接口
type Storage interface {
    // 任务持久化
    SaveTask(task *Task) error
    LoadTasks() ([]*Task, error)
    DeleteTask(taskID string) error
    
    // 时间轮状态持久化
    SaveWheelState(level int, currentTick int64) error
    LoadWheelState(level int) (int64, error)
}

// RedisStorage Redis实现
type RedisStorage struct {
    client *redis.Client
}

func (s *RedisStorage) SaveTask(task *Task) error {
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }
    
    // 使用Redis的ZADD存储任务，score为执行时间戳
    return s.client.ZAdd(
        context.Background(),
        "timing_wheel:tasks",
        redis.Z{
            Score:  float64(task.ExecuteAt.Unix()),
            Member: data,
        },
    ).Err()
}

func (s *RedisStorage) LoadTasks() ([]*Task, error) {
    // 加载所有未执行的任务
    now := time.Now().Unix()
    
    results, err := s.client.ZRangeByScore(
        context.Background(),
        "timing_wheel:tasks",
        &redis.ZRangeBy{
            Min: "-inf",
            Max: "+inf",
        },
    ).Result()
    
    if err != nil {
        return nil, err
    }
    
    tasks := make([]*Task, 0, len(results))
    for _, data := range results {
        var task Task
        if err := json.Unmarshal([]byte(data), &task); err != nil {
            continue
        }
        tasks = append(tasks, &task)
    }
    
    return tasks, nil
}
```

### 7.2 崩溃恢复

```go
// Recover 恢复时间轮状态
func (m *TimingWheelManager) Recover() error {
    if !m.enablePersist {
        return nil
    }
    
    // 1. 恢复时间轮状态
    for _, wheel := range m.wheels {
        currentTick, err := m.storage.LoadWheelState(wheel.level)
        if err == nil {
            wheel.currentTick = currentTick
        }
    }
    
    // 2. 加载未完成的任务
    tasks, err := m.storage.LoadTasks()
    if err != nil {
        return err
    }
    
    // 3. 重新调度任务
    now := time.Now()
    for _, task := range tasks {
        if task.Status == TaskStatusPending {
            // 重新计算延迟时间
            if task.ExecuteAt.After(now) {
                task.Delay = task.ExecuteAt.Sub(now)
            } else {
                task.Delay = 0 // 立即执行
            }
            
            m.AddTask(task)
        }
    }
    
    return nil
}

// 定期持久化时间轮状态
func (m *TimingWheelManager) startStatePersistence() {
    if !m.enablePersist {
        return
    }
    
    ticker := time.NewTicker(10 * time.Second)
    go func() {
        for range ticker.C {
            for _, wheel := range m.wheels {
                m.storage.SaveWheelState(wheel.level, wheel.currentTick)
            }
        }
    }()
}
```

---

## 八、性能优化策略

### 8.1 内存优化

```go
// 1. 对象池复用
var taskNodePool = sync.Pool{
    New: func() interface{} {
        return &TaskNode{}
    },
}

func acquireTaskNode() *TaskNode {
    return taskNodePool.Get().(*TaskNode)
}

func releaseTaskNode(node *TaskNode) {
    node.task = nil
    node.next = nil
    node.prev = nil
    taskNodePool.Put(node)
}

// 2. 槽位预分配
func NewTimeWheel(config WheelConfig) *TimeWheel {
    slots := make([]TaskList, config.SlotCount)
    for i := range slots {
        slots[i] = TaskList{
            // 预分配一定容量，减少扩容
        }
    }
    
    return &TimeWheel{
        slots: slots,
        // ...
    }
}

// 3. 使用环形缓冲区存储任务
type RingBuffer struct {
    buffer []*Task
    head   int
    tail   int
    size   int
    cap    int
}
```

### 8.2 并发优化

```go
// 1. 分段锁：为每个槽位独立加锁
type ShardedTimeWheel struct {
    shards    []TimeWheelShard
    shardMask int
}

type TimeWheelShard struct {
    slots []TaskList
    mu    sync.RWMutex
}

func (w *ShardedTimeWheel) getShard(slotIndex int) *TimeWheelShard {
    return &w.shards[slotIndex&w.shardMask]
}

// 2. 无锁数据结构（适用于单生产者单消费者场景）
type LockFreeQueue struct {
    head unsafe.Pointer
    tail unsafe.Pointer
}

// 3. 批量操作减少锁竞争
func (w *TimeWheel) addBatch(tasks []*Task) {
    // 按槽位分组
    groups := make(map[int][]*Task)
    for _, task := range tasks {
        groups[task.SlotIndex] = append(groups[task.SlotIndex], task)
    }
    
    // 批量添加
    for slotIndex, taskGroup := range groups {
        w.mu.Lock()
        for _, task := range taskGroup {
            w.addToSlotUnsafe(slotIndex, task)
        }
        w.mu.Unlock()
    }
}
```

### 8.3 算法优化

```go
// 1. 快速取模（位运算代替除法）
// 要求槽位数量为2的幂
func (w *TimeWheel) getSlotIndex(tick int64) int {
    return int(tick & int64(w.wheelSize-1))
}

// 2. 延迟加载：槽位按需初始化
type LazyTimeWheel struct {
    slots sync.Map // map[int]*TaskList
}

func (w *LazyTimeWheel) getOrCreateSlot(index int) *TaskList {
    if slot, ok := w.slots.Load(index); ok {
        return slot.(*TaskList)
    }
    
    newSlot := &TaskList{}
    actual, _ := w.slots.LoadOrStore(index, newSlot)
    return actual.(*TaskList)
}

// 3. 稀疏时间轮：只存储有任务的槽位
type SparseTimeWheel struct {
    slots    map[int]*TaskList
    slotsMu  sync.RWMutex
}
```

---

## 九、监控与可观测性

### 9.1 监控指标

```go
// Metrics 监控指标
type Metrics struct {
    // 任务指标
    TaskAdded     prometheus.Counter   // 添加任务数
    TaskExecuted  prometheus.Counter   // 执行任务数
    TaskFailed    prometheus.Counter   // 失败任务数
    TaskCanceled  prometheus.Counter   // 取消任务数
    
    // 性能指标
    TaskLatency   prometheus.Histogram // 任务延迟（实际执行时间-预期执行时间）
    ExecutionTime prometheus.Histogram // 任务执行耗时
    
    // 资源指标
    ActiveTasks   prometheus.Gauge     // 当前活跃任务数
    WheelTick     prometheus.Counter   // 时间轮tick次数
    QueueSize     prometheus.Gauge     // 工作队列大小
    
    // 层级指标
    TasksByLevel  *prometheus.GaugeVec // 各层级任务数
}

func (m *Metrics) Register() {
    prometheus.MustRegister(
        m.TaskAdded,
        m.TaskExecuted,
        m.TaskFailed,
        m.TaskLatency,
        m.ExecutionTime,
        m.ActiveTasks,
    )
}

// 使用示例
func (w *TimeWheel) executeTask(task *Task) {
    start := time.Now()
    
    // 记录延迟
    latency := start.Sub(task.ExecuteAt)
    w.metrics.TaskLatency.Observe(latency.Seconds())
    
    // 执行任务
    err := task.Callback(context.Background(), task)
    
    // 记录执行时间
    duration := time.Since(start)
    w.metrics.ExecutionTime.Observe(duration.Seconds())
    
    if err != nil {
        w.metrics.TaskFailed.Inc()
    } else {
        w.metrics.TaskExecuted.Inc()
    }
}
```

### 9.2 日志记录

```go
// Logger 结构化日志
type Logger struct {
    logger *zap.Logger
}

func (l *Logger) TaskAdded(task *Task) {
    l.logger.Info("task added",
        zap.String("task_id", task.ID),
        zap.String("task_name", task.Name),
        zap.Duration("delay", task.Delay),
        zap.Time("execute_at", task.ExecuteAt),
        zap.Int("level", task.Level),
    )
}

func (l *Logger) TaskExecuted(task *Task, err error) {
    if err != nil {
        l.logger.Error("task execution failed",
            zap.String("task_id", task.ID),
            zap.Error(err),
            zap.Int("retries", task.Retries),
        )
    } else {
        l.logger.Info("task executed successfully",
            zap.String("task_id", task.ID),
            zap.Duration("actual_delay", time.Since(task.CreatedAt)),
        )
    }
}

func (l *Logger) WheelAdvanced(level int, tick int64) {
    if tick%100 == 0 { // 每100次tick记录一次
        l.logger.Debug("wheel advanced",
            zap.Int("level", level),
            zap.Int64("tick", tick),
        )
    }
}
```

### 9.3 分布式追踪

```go
// 集成OpenTelemetry
func (w *TimeWheel) executeTaskWithTracing(task *Task) {
    ctx := task.Context
    
    // 创建span
    ctx, span := otel.Tracer("timing-wheel").Start(ctx, "execute_task")
    defer span.End()
    
    // 添加属性
    span.SetAttributes(
        attribute.String("task.id", task.ID),
        attribute.String("task.name", task.Name),
        attribute.Int("task.level", task.Level),
    )
    
    // 执行任务
    err := task.Callback(ctx, task)
    
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    } else {
        span.SetStatus(codes.Ok, "success")
    }
}
```

---

## 十、高可用设计

### 10.1 主从复制

```go
// ReplicatedTimingWheel 主从复制时间轮
type ReplicatedTimingWheel struct {
    local  *TimingWheelManager
    peers  []Peer
    role   Role // Master/Slave
    
    replicator *Replicator
}

type Replicator struct {
    logWriter  *wal.Writer  // 写入WAL日志
    logReader  *wal.Reader  // 读取WAL日志
    syncChan   chan LogEntry
}

type LogEntry struct {
    OpType    string // Add/Cancel/Execute
    Task      *Task
    Timestamp time.Time
}

// Replicate 主节点复制操作到从节点
func (r *Replicator) Replicate(entry LogEntry) error {
    // 1. 写入本地WAL
    if err := r.logWriter.Write(entry); err != nil {
        return err
    }
    
    // 2. 发送到从节点
    for _, peer := range r.peers {
        go peer.Send(entry)
    }
    
    return nil
}

// Follow 从节点跟随主节点
func (r *Replicator) Follow(master Peer) {
    stream := master.Subscribe()
    
    for entry := range stream {
        // 应用日志条目
        r.apply(entry)
    }
}
```

### 10.2 分布式一致性

```go
// DistributedTimingWheel 分布式时间轮
type DistributedTimingWheel struct {
    local      *TimingWheelManager
    coordinator Coordinator // Raft/Etcd
    nodeID     string
}

// AddTask 添加任务（需要达成一致）
func (d *DistributedTimingWheel) AddTask(task *Task) error {
    // 1. 提交到一致性协议
    proposal := &Proposal{
        Type: "AddTask",
        Data: task,
    }
    
    // 2. 等待多数节点确认
    result, err := d.coordinator.Propose(proposal)
    if err != nil {
        return err
    }
    
    // 3. 应用到本地时间轮
    if result.Committed {
        return d.local.AddTask(task)
    }
    
    return errors.New("proposal rejected")
}
```

### 10.3 故障检测与切换

```go
// HealthChecker 健康检查
type HealthChecker struct {
    timingWheel *TimingWheelManager
    peers       []Peer
    
    heartbeatInterval time.Duration
    failureThreshold  int
}

func (h *HealthChecker) Start() {
    ticker := time.NewTicker(h.heartbeatInterval)
    
    go func() {
        for range ticker.C {
            for _, peer := range h.peers {
                if !peer.IsHealthy() {
                    h.handleFailure(peer)
                }
            }
        }
    }()
}

func (h *HealthChecker) handleFailure(peer Peer) {
    peer.FailureCount++
    
    if peer.FailureCount >= h.failureThreshold {
        // 触发故障转移
        h.triggerFailover(peer)
    }
}

func (h *HealthChecker) triggerFailover(failed Peer) {
    // 1. 标记节点为不可用
    failed.Status = NodeStatusDown
    
    // 2. 选举新的主节点（如果失败的是主节点）
    if failed.Role == RoleMaster {
        h.electNewMaster()
    }
    
    // 3. 重新分配任务
    h.redistributeTasks(failed)
}
```

---

## 十一、测试方案

### 11.1 单元测试

```go
func TestTimeWheelAddTask(t *testing.T) {
    tw := NewTimeWheel(WheelConfig{
        Tick:      100 * time.Millisecond,
        SlotCount: 10,
        Level:     0,
    })
    
    executed := false
    task := &Task{
        ID:    "test-1",
        Delay: 500 * time.Millisecond,
        Callback: func(ctx context.Context, t *Task) error {
            executed = true
            return nil
        },
    }
    
    err := tw.AddTask(task)
    assert.Nil(t, err)
    
    tw.Start()
    
    // 等待任务执行
    time.Sleep(600 * time.Millisecond)
    
    assert.True(t, executed)
}

func TestTimeWheelCancelTask(t *testing.T) {
    tw := NewTimeWheel(/* ... */)
    
    task := &Task{
        ID:    "test-cancel",
        Delay: 1 * time.Second,
        Callback: func(ctx context.Context, t *Task) error {
            t.Error("task should not execute")
            return nil
        },
    }
    
    tw.AddTask(task)
    time.Sleep(100 * time.Millisecond)
    
    err := tw.CancelTask("test-cancel")
    assert.Nil(t, err)
    
    time.Sleep(1 * time.Second)
    // 任务不应该执行
}
```

### 11.2 基准测试

```go
func BenchmarkTimeWheelAddTask(b *testing.B) {
    tw := NewTimingWheelManager(DefaultWheelConfigs)
    tw.Start()
    defer tw.Stop()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        task := &Task{
            ID:    fmt.Sprintf("task-%d", i),
            Delay: time.Duration(rand.Intn(1000)) * time.Millisecond,
            Callback: func(ctx context.Context, t *Task) error {
                return nil
            },
        }
        tw.AddTask(task)
    }
}

func BenchmarkTimeWheelExecute(b *testing.B) {
    tw := NewTimingWheelManager(DefaultWheelConfigs)
    
    // 预先添加大量任务
    for i := 0; i < 10000; i++ {
        task := &Task{
            ID:    fmt.Sprintf("task-%d", i),
            Delay: 100 * time.Millisecond,
            Callback: func(ctx context.Context, t *Task) error {
                time.Sleep(1 * time.Microsecond)
                return nil
            },
        }
        tw.AddTask(task)
    }
    
    b.ResetTimer()
    tw.Start()
    
    // 等待所有任务执行完成
    time.Sleep(200 * time.Millisecond)
}
```

### 11.3 压力测试

```go
func TestTimeWheelStress(t *testing.T) {
    tw := NewTimingWheelManager(DefaultWheelConfigs)
    tw.Start()
    defer tw.Stop()
    
    const (
        taskCount    = 100000
        concurrency  = 100
    )
    
    var wg sync.WaitGroup
    var executedCount atomic.Int64
    
    // 并发添加任务
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(base int) {
            defer wg.Done()
            
            for j := 0; j < taskCount/concurrency; j++ {
                task := &Task{
                    ID:    fmt.Sprintf("task-%d-%d", base, j),
                    Delay: time.Duration(rand.Intn(10000)) * time.Millisecond,
                    Callback: func(ctx context.Context, t *Task) error {
                        executedCount.Add(1)
                        return nil
                    },
                }
                
                if err := tw.AddTask(task); err != nil {
                    t.Errorf("failed to add task: %v", err)
                }
            }
        }(i)
    }
    
    wg.Wait()
    
    // 等待所有任务执行
    time.Sleep(15 * time.Second)
    
    executed := executedCount.Load()
    t.Logf("Added: %d, Executed: %d", taskCount, executed)
    
    // 允许5%的误差（可能有些任务还在队列中）
    assert.Greater(t, executed, int64(taskCount*0.95))
}
```

---

## 十二、使用示例

### 12.1 初始化

```go
func main() {
    // 1. 创建时间轮管理器
    manager := NewTimingWheelManager(TimingWheelConfig{
        WheelConfigs: DefaultWheelConfigs,
        WorkerCount:  100,
        QueueSize:    10000,
        EnablePersist: true,
        Storage:      NewRedisStorage(redisClient),
    })
    
    // 2. 启动时间轮
    if err := manager.Start(); err != nil {
        log.Fatal(err)
    }
    defer manager.Stop()
    
    // 3. 恢复状态（如果需要）
    if err := manager.Recover(); err != nil {
        log.Printf("recover failed: %v", err)
    }
    
    // 4. 设置监控
    setupMonitoring(manager)
    
    // 5. 启动API服务
    startAPIServer(manager)
}
```

### 12.2 添加一次性任务

```go
// 5秒后发送通知
task := &Task{
    ID:    "notify-001",
    Name:  "Send Notification",
    Delay: 5 * time.Second,
    Type:  TaskTypeOnce,
    Callback: func(ctx context.Context, t *Task) error {
        return notificationService.Send("user-123", "您的任务已完成")
    },
}

manager.AddTask(task)
```

### 12.3 添加周期性任务

```go
// 每10分钟执行一次数据同步
task := &Task{
    ID:       "sync-001",
    Name:     "Data Sync",
    Delay:    10 * time.Minute,
    Type:     TaskTypePeriodic,
    Interval: 10 * time.Minute,
    RepeatCount: -1, // 无限循环
    Callback: func(ctx context.Context, t *Task) error {
        return dataService.Sync()
    },
}

manager.AddTask(task)
```

### 12.4 取消任务

```go
// 取消任务
if err := manager.CancelTask("notify-001"); err != nil {
    log.Printf("cancel task failed: %v", err)
}
```

---

## 十三、注意事项与最佳实践

### 13.1 配置建议

```go
// 1. 根据业务场景选择合适的层级配置
// 高频短延迟场景（如游戏）
var GameWheelConfigs = []WheelConfig{
    {Tick: 10 * time.Millisecond, SlotCount: 100, Level: 0},   // 0-1s
    {Tick: 100 * time.Millisecond, SlotCount: 100, Level: 1},  // 0-10s
    {Tick: 1 * time.Second, SlotCount: 300, Level: 2},         // 0-5min
}

// 低频长延迟场景（如订单超时）
var OrderWheelConfigs = []WheelConfig{
    {Tick: 1 * time.Second, SlotCount: 60, Level: 0},     // 0-1min
    {Tick: 1 * time.Minute, SlotCount: 60, Level: 1},     // 0-1hour
    {Tick: 1 * time.Hour, SlotCount: 24, Level: 2},       // 0-1day
    {Tick: 1 * time.Day, SlotCount: 30, Level: 3},        // 0-30days
}

// 2. 工作协程池大小
// CPU密集型：核心数 * 1-2
// IO密集型：核心数 * 2-4
workerCount := runtime.NumCPU() * 2

// 3. 队列大小：根据任务量和内存限制
queueSize := 10000
```

### 13.2 避免的问题

```go
// 1. 避免在回调中执行耗时操作
// ❌ 错误示例
task.Callback = func(ctx context.Context, t *Task) error {
    time.Sleep(10 * time.Second) // 阻塞工作协程
    return nil
}

// ✅ 正确示例
task.Callback = func(ctx context.Context, t *Task) error {
    go func() {
        // 异步执行耗时操作
        heavyOperation()
    }()
    return nil
}

// 2. 避免任务回调中出现panic
task.Callback = func(ctx context.Context, t *Task) error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("task panic: %v", r)
        }
    }()
    
    // 业务逻辑
    return nil
}

// 3. 注意任务ID的唯一性
// 使用UUID或雪花算法生成唯一ID
task.ID = uuid.New().String()

// 4. 合理设置超时时间
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
task.Context = ctx
```

### 13.3 监控告警

```go
// 1. 关键指标告警
if manager.metrics.TaskFailed.Count() > 100 {
    alerting.Send("时间轮任务失败率过高")
}

if manager.metrics.ActiveTasks.Value() > 100000 {
    alerting.Send("时间轮任务堆积")
}

// 2. 延迟监控
if manager.metrics.TaskLatency.P99() > 1*time.Second {
    alerting.Send("任务执行延迟过高")
}

// 3. 资源监控
if manager.workerPool.QueueSize() > queueSize*0.8 {
    alerting.Send("工作队列接近饱和")
}
```

---

## 十四、与业界方案对比

| 特性             | Kafka时间轮 | Netty时间轮 | 本方案          |
|------------------|-------------|-------------|-----------------|
| 层级支持         | 单层        | 单层/多层   | 多层            |
| 持久化           | ✓           | ✗           | ✓               |
| 分布式           | ✓           | ✗           | ✓               |
| 任务取消         | ✓           | ✓           | ✓               |
| 周期任务         | ✗           | ✓           | ✓               |
| 精度             | 毫秒级      | 毫秒级      | 可配置（毫秒-天）|
| 最大延迟         | 较短        | 较短        | 30天+           |
| 使用场景         | 请求超时    | 网络超时    | 工作流调度      |

---

## 十五、总结

### 15.1 核心优势

1. **高性能**：O(1)的任务添加、删除、执行复杂度
2. **低内存**：相比传统定时器，内存占用大幅降低
3. **高精度**：支持毫秒级精度，同时支持长时间延迟
4. **可扩展**：支持水平扩展和分布式部署
5. **高可靠**：支持持久化和崩溃恢复

### 15.2 适用场景

- 工作流引擎的任务超时检测
- 延迟任务调度（如延迟消息、延迟通知）
- 周期性任务执行（如定时清理、定时同步）
- 订单超时自动取消
- 游戏Buff定时器
- 网络请求超时控制

### 15.3 实施路线图

**阶段一**（1-2周）：
- 实现单层时间轮基础框架
- 实现任务添加、执行、取消
- 编写单元测试

**阶段二**（2-3周）：
- 实现多层时间轮
- 实现任务降级和升级
- 实现工作协程池

**阶段三**（1-2周）：
- 实现持久化和恢复
- 实现监控指标
- 性能优化

**阶段四**（2-3周）：
- 实现分布式支持
- 实现高可用
- 压力测试和调优

**阶段五**（1周）：
- 与工作流引擎集成
- 完善文档和示例
- 生产环境试运行

---

## 附录

### A. 参考资料

1. **论文**：
   - "Hashed and Hierarchical Timing Wheels" - George Varghese, Tony Lauck (1987)
   - "Timer Management in a High-Performance Protocol Implementation" - C. Maeda, B. Bershad (1993)

2. **开源项目**：
   - Apache Kafka: `kafka.utils.timer.TimingWheel`
   - Netty: `io.netty.util.HashedWheelTimer`
   - Go语言: `github.com/RussellLuo/timingwheel`

3. **博客文章**：
   - Kafka如何实现延时队列
   - Netty时间轮原理分析
   - 分层时间轮设计与实现

### B. 示例代码仓库

完整实现代码：`https://github.com/your-org/timing-wheel`

```bash
git clone https://github.com/your-org/timing-wheel.git
cd timing-wheel
go test -v ./...
go run examples/basic/main.go
```

### C. API文档

详细API文档：`https://pkg.go.dev/your-org/timing-wheel`

---

**文档版本**：v1.0  
**最后更新**：2025-12-31  
**维护者**：工作流引擎团队
