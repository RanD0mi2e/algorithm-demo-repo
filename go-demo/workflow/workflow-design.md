# 工作流引擎设计方案

## 概述

本文档描述了一个用于处理单据流转的工作流引擎的完整设计方案。该引擎支持动态流程定义、多种审批模式、条件分支、并行处理等企业级功能。

---

## 一、系统架构概览

### 1.1 架构图

```
┌─────────────────────────────────────────────────────────────────────────┐
│                            API Gateway                                   │
│                    (认证、限流、路由、日志)                                │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Workflow API Layer                               │
│              (REST API / gRPC / GraphQL 接口层)                          │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
        ┌───────────────────────────┼───────────────────────────┐
        ▼                           ▼                           ▼
┌───────────────┐          ┌───────────────┐          ┌───────────────┐
│  流程定义服务  │          │  流程执行引擎  │          │   任务管理服务  │
│ (Definition)  │          │  (Execution)  │          │    (Task)     │
└───────────────┘          └───────────────┘          └───────────────┘
        │                           │                           │
        └───────────────────────────┼───────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          Core Domain Layer                               │
│    ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐              │
│    │ 状态机   │  │ 条件引擎  │  │ 审批策略  │  │ 事件总线  │              │
│    └──────────┘  └──────────┘  └──────────┘  └──────────┘              │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Infrastructure Layer                              │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐    │
│  │ MySQL  │ │ Redis  │ │  MQ    │ │ 分布式锁│ │ 日志   │ │ 监控   │    │
│  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘ └────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.2 核心组件说明

| 组件 | 职责 | 技术选型 |
|------|------|----------|
| API Gateway | 统一入口、认证鉴权、限流熔断 | Kong/Nginx |
| 流程定义服务 | 流程模板的CRUD、版本管理、验证 | Go Service |
| 流程执行引擎 | 流程实例的创建、状态流转、并发控制 | Go Service |
| 任务管理服务 | 任务分配、查询、超时监控 | Go Service |
| 事件总线 | 解耦各模块、异步通知 | Kafka/RabbitMQ |
| 分布式锁 | 并发控制、幂等保证 | Redis/etcd |

---

## 二、核心数据模型

### 2.1 工作流定义

```go
// 工作流定义
type Workflow struct {
    ID          string       // 工作流唯一标识
    Name        string       // 工作流名称
    Version     string       // 版本号
    Description string       // 描述
    Nodes       []Node       // 流程节点
    Transitions []Transition // 节点间的转换关系
    Variables   []Variable   // 流程变量定义
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Status      string       // Active/Inactive/Archived
}

// 流程变量定义
type Variable struct {
    Name         string      // 变量名
    Type         string      // 类型：string/number/boolean/object/array
    DefaultValue interface{} // 默认值
    Required     bool        // 是否必填
    Description  string      // 描述
}

// 节点类型
type NodeType string

const (
    StartNode     NodeType = "start"      // 开始节点
    ApprovalNode  NodeType = "approval"   // 审批节点
    ConditionNode NodeType = "condition"  // 条件节点
    ParallelNode  NodeType = "parallel"   // 并行网关
    EndNode       NodeType = "end"        // 结束节点
)

// 流程节点
type Node struct {
    ID          string
    Name        string
    Type        NodeType
    Config      NodeConfig    // 节点配置
    Assignee    AssigneeRule  // 审批人规则
    Actions     []Action      // 可执行的操作
    TimeoutRule *TimeoutRule  // 超时规则
}

// 节点配置
type NodeConfig struct {
    RequireComment bool              // 是否必须填写意见
    AllowDelegate  bool              // 是否允许委托
    AllowReject    bool              // 是否允许拒绝
    CustomFields   map[string]string // 自定义字段
}

// 审批人规则
type AssigneeRule struct {
    Type       string   // User/Role/Department/Expression/Dynamic
    Values     []string // 具体值（用户ID、角色ID等）
    Expression string   // 表达式（用于动态计算）
    Mode       string   // Single/Multi/Sequential/Parallel
}

// 状态转换
type Transition struct {
    ID        string
    Name      string
    From      string // 源节点ID
    To        string // 目标节点ID
    Condition string // 转换条件表达式
    Priority  int    // 优先级（多个条件时使用）
}

// 操作定义
type Action struct {
    ID    string
    Name  string // 通过/拒绝/转交/加签/退回等
    Type  string
    Color string // UI显示颜色
}

// 超时规则
type TimeoutRule struct {
    Duration time.Duration // 超时时长
    Action   string        // Auto/Escalate/Notify
    Target   string        // 升级目标（上级/特定人员）
}
```

### 2.2 流程实例

```go
// 流程实例
type WorkflowInstance struct {
    ID           string
    WorkflowID   string                 // 关联的工作流定义
    WorkflowName string
    DocumentID   string                 // 关联的业务单据ID
    DocumentType string                 // 单据类型
    CurrentNode  string                 // 当前所在节点
    Status       InstanceStatus         // 实例状态
    Variables    map[string]interface{} // 流程变量
    Initiator    string                 // 发起人
    CreatedAt    time.Time
    UpdatedAt    time.Time
    CompletedAt  *time.Time
}

// 实例状态
type InstanceStatus string

const (
    StatusRunning   InstanceStatus = "running"   // 运行中
    StatusCompleted InstanceStatus = "completed" // 已完成
    StatusRejected  InstanceStatus = "rejected"  // 已拒绝
    StatusTerminated InstanceStatus = "terminated" // 已终止
    StatusSuspended InstanceStatus = "suspended" // 已挂起
)
```

### 2.3 任务

```go
// 任务
type Task struct {
    ID          string
    InstanceID  string     // 所属流程实例
    NodeID      string     // 所属节点
    NodeName    string
    Assignee    string     // 当前处理人
    Status      TaskStatus
    Action      string     // 执行的操作
    Comment     string     // 审批意见
    Variables   map[string]interface{} // 任务变量
    CreatedAt   time.Time
    CompletedAt *time.Time
    DueDate     *time.Time // 期望完成时间
}

// 任务状态
type TaskStatus string

const (
    TaskPending   TaskStatus = "pending"   // 待处理
    TaskCompleted TaskStatus = "completed" // 已完成
    TaskRejected  TaskStatus = "rejected"  // 已拒绝
    TaskCanceled  TaskStatus = "canceled"  // 已取消
    TaskDelegated TaskStatus = "delegated" // 已委托
)
```

### 2.4 历史记录

```go
// 历史记录
type HistoryLog struct {
    ID          string
    InstanceID  string
    TaskID      string
    NodeID      string
    NodeName    string
    Action      string
    Operator    string
    OperatorName string
    Comment     string
    Variables   map[string]interface{}
    Timestamp   time.Time
}
```

### 2.5 流程令牌（用于并行网关）

```go
// ProcessToken 流程令牌，用于跟踪并行分支的执行状态
type ProcessToken struct {
    ID          string
    InstanceID  string    // 所属流程实例
    NodeID      string    // 当前所在节点
    ParentToken string    // 父令牌ID（用于嵌套并行）
    Status      string    // active/completed/canceled
    CreatedAt   time.Time
    CompletedAt *time.Time
}
```

### 2.6 任务委托记录

```go
// TaskDelegation 任务委托记录
type TaskDelegation struct {
    ID        string
    TaskID    string    // 被委托的任务
    FromUser  string    // 委托人
    ToUser    string    // 被委托人
    Comment   string    // 委托说明
    CreatedAt time.Time
}
```

### 2.7 流程关系（用于子流程）

```go
// ProcessRelation 流程关系，用于子流程关联
type ProcessRelation struct {
    ID               string
    ParentInstanceID string    // 父流程实例ID
    ChildInstanceID  string    // 子流程实例ID
    ParentNodeID     string    // 父流程中触发子流程的节点
    Type             string    // subprocess/call_activity
    Status           string    // running/completed/failed
    CreatedAt        time.Time
    CompletedAt      *time.Time
}
```

---

## 三、核心引擎接口

### 3.1 工作流引擎主接口

```go
// WorkflowEngine 工作流引擎接口
type WorkflowEngine interface {
    // 流程管理
    DeployWorkflow(workflow *Workflow) error
    GetWorkflow(workflowID string) (*Workflow, error)
    ListWorkflows() ([]*Workflow, error)
    
    // 流程实例管理
    StartProcess(workflowID string, documentID string, initiator string, variables map[string]interface{}) (*WorkflowInstance, error)
    GetInstance(instanceID string) (*WorkflowInstance, error)
    ListInstances(filter InstanceFilter) ([]*WorkflowInstance, error)
    
    // 任务管理
    CompleteTask(taskID string, userID string, action string, comment string, variables map[string]interface{}) error
    GetTask(taskID string) (*Task, error)
    GetPendingTasks(userID string) ([]*Task, error)
    DelegateTask(taskID string, fromUser string, toUser string, comment string) error
    
    // 流程控制
    Rollback(instanceID string, targetNodeID string, reason string) error
    Terminate(instanceID string, reason string) error
    Suspend(instanceID string, reason string) error
    Resume(instanceID string) error
    
    // 历史查询
    GetHistory(instanceID string) ([]*HistoryLog, error)
    GetTasksByInstance(instanceID string) ([]*Task, error)
}

// InstanceFilter 实例查询过滤器
type InstanceFilter struct {
    WorkflowID   string
    DocumentType string
    Status       InstanceStatus
    Initiator    string
    StartTime    *time.Time
    EndTime      *time.Time
    PageSize     int
    PageNum      int
}
```

### 3.2 仓储接口定义

```go
// WorkflowRepository 工作流定义仓储
type WorkflowRepository interface {
    Save(workflow *Workflow) error
    GetByID(id string) (*Workflow, error)
    GetByNameAndVersion(name, version string) (*Workflow, error)
    List(filter WorkflowFilter) ([]*Workflow, error)
    Update(workflow *Workflow) error
    Delete(id string) error
}

// InstanceRepository 流程实例仓储
type InstanceRepository interface {
    Save(instance *WorkflowInstance) error
    GetByID(id string) (*WorkflowInstance, error)
    Update(instance *WorkflowInstance) error
    List(filter InstanceFilter) ([]*WorkflowInstance, error)
    CountByStatus(status InstanceStatus) (int64, error)
    Ping() error // 健康检查
}

// TaskRepository 任务仓储
type TaskRepository interface {
    Save(task *Task) error
    GetByID(id string) (*Task, error)
    Update(task *Task) error
    FindByAssignee(userID string, status TaskStatus) ([]*Task, error)
    FindByRoles(roles []string, status TaskStatus) ([]*Task, error)
    FindByNodeAndInstance(nodeID, instanceID string) ([]*Task, error)
    FindOverdueTasks() ([]*Task, error)
    CountPendingByUser(userID string) (int64, error)
}

// HistoryRepository 历史记录仓储
type HistoryRepository interface {
    Add(log *HistoryLog) error
    GetByInstanceID(instanceID string) ([]*HistoryLog, error)
    GetByTaskID(taskID string) ([]*HistoryLog, error)
}

// TokenRepository 流程令牌仓储
type TokenRepository interface {
    Save(token *ProcessToken) error
    GetByID(id string) (*ProcessToken, error)
    FindByNodeID(nodeID string) ([]*ProcessToken, error)
    FindByInstanceID(instanceID string) ([]*ProcessToken, error)
    Update(token *ProcessToken) error
}

// RelationRepository 流程关系仓储
type RelationRepository interface {
    Save(relation *ProcessRelation) error
    GetByChildID(childInstanceID string) (*ProcessRelation, error)
    FindByParentID(parentInstanceID string) ([]*ProcessRelation, error)
}
```

### 3.3 外部服务接口

```go
// UserService 用户服务接口
type UserService interface {
    GetUser(userID string) (*User, error)
    GetUsersByRole(roleID string) []string
    GetUserRoles(userID string) []string
    GetUserManager(userID string) (string, error)
}

// OrganizationService 组织架构服务接口
type OrganizationService interface {
    GetUsersByDepartment(deptID string) []string
    GetDepartmentManager(deptID string) (string, error)
    GetUserDepartment(userID string) (string, error)
}

// Notifier 通知服务接口
type Notifier interface {
    Notify(userID string, message string) error
    NotifyBatch(userIDs []string, message string) error
    SendEmail(to, subject, body string) error
    SendSMS(phone, message string) error
}

// User 用户信息
type User struct {
    ID         string
    Name       string
    Email      string
    Phone      string
    Department string
    Roles      []string
}
```

### 3.4 分布式锁接口

```go
// DistributedLock 分布式锁接口
type DistributedLock interface {
    // Lock 获取锁，返回锁标识
    Lock(key string, ttl time.Duration) (string, error)
    // Unlock 释放锁
    Unlock(key string, token string) error
    // Renew 续期
    Renew(key string, token string, ttl time.Duration) error
}

// RedisDistributedLock Redis实现的分布式锁
type RedisDistributedLock struct {
    client *redis.Client
}

func (l *RedisDistributedLock) Lock(key string, ttl time.Duration) (string, error) {
    token := generateID()
    ok, err := l.client.SetNX(context.Background(), "lock:"+key, token, ttl).Result()
    if err != nil {
        return "", err
    }
    if !ok {
        return "", ErrLockFailed
    }
    return token, nil
}

func (l *RedisDistributedLock) Unlock(key string, token string) error {
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
    _, err := l.client.Eval(context.Background(), script, []string{"lock:" + key}, token).Result()
    return err
}
```

---

## 四、核心设计模式

### 4.1 状态机模式

每个流程实例都是一个状态机，节点代表状态，转换代表状态变化的触发条件。

```go
type StateMachine struct {
    instance    *WorkflowInstance
    workflow    *Workflow
    currentNode *Node
}

func (sm *StateMachine) Transit(action string, variables map[string]interface{}) error {
    // 1. 查找可用的转换
    transitions := sm.findAvailableTransitions(action)
    
    // 2. 评估转换条件
    for _, trans := range transitions {
        if sm.evaluateCondition(trans.Condition, variables) {
            // 3. 执行状态转换
            return sm.executeTransition(trans)
        }
    }
    
    return errors.New("no valid transition found")
}
```

### 4.2 责任链模式

处理多级审批、会签等场景。

```go
type ApprovalHandler interface {
    SetNext(handler ApprovalHandler)
    Handle(task *Task) error
}

type SequentialApprovalHandler struct {
    next ApprovalHandler
    assignees []string
}

func (h *SequentialApprovalHandler) Handle(task *Task) error {
    for _, assignee := range h.assignees {
        // 创建任务并等待完成
        if err := h.createAndWaitTask(assignee, task); err != nil {
            return err
        }
    }
    
    if h.next != nil {
        return h.next.Handle(task)
    }
    return nil
}
```

### 4.3 策略模式

不同的审批模式（单人、会签、或签）使用不同的策略。

```go
type ApprovalStrategy interface {
    CreateTasks(node *Node, instance *WorkflowInstance) ([]*Task, error)
    IsComplete(tasks []*Task) bool
}

// 单人审批策略
type SingleApprovalStrategy struct{}

func (s *SingleApprovalStrategy) CreateTasks(node *Node, instance *WorkflowInstance) ([]*Task, error) {
    assignee := s.resolveAssignee(node, instance)
    return []*Task{s.createTask(assignee, node, instance)}, nil
}

// 会签策略（所有人都要同意）
type CountersignStrategy struct{}

func (s *CountersignStrategy) IsComplete(tasks []*Task) bool {
    for _, task := range tasks {
        if task.Status != TaskCompleted {
            return false
        }
    }
    return true
}

// 或签策略（任意一人同意即可）
type OrSignStrategy struct{}

func (s *OrSignStrategy) IsComplete(tasks []*Task) bool {
    for _, task := range tasks {
        if task.Status == TaskCompleted {
            return true
        }
    }
    return false
}
```

---

## 五、核心功能模块

### 5.1 流程定义管理

- **可视化设计器**：拖拽式流程设计
- **JSON/XML定义**：支持代码化定义
- **版本管理**：支持多版本并存，灰度发布
- **导入导出**：流程定义的备份和迁移

### 5.2 流程执行引擎

```go
type ExecutionEngine struct {
    workflowRepo   WorkflowRepository
    instanceRepo   InstanceRepository
    taskRepo       TaskRepository
    historyRepo    HistoryRepository
    tokenRepo      TokenRepository      // 流程令牌仓储（并行网关用）
    relationRepo   RelationRepository   // 流程关系仓储（子流程用）
    eventBus       EventBus
    conditionEval  ConditionEvaluator
    locker         DistributedLock      // 分布式锁
    idGenerator    IDGenerator          // ID生成器
}

// IDGenerator ID生成器接口
type IDGenerator interface {
    Generate() string
}

// 辅助函数
func generateID() string {
    return uuid.New().String()
}

func timePtr(t time.Time) *time.Time {
    return &t
}

func (e *ExecutionEngine) StartProcess(workflowID string, documentID string, initiator string, variables map[string]interface{}) (*WorkflowInstance, error) {
    // 幂等性检查：防止重复创建
    lockKey := fmt.Sprintf("start_process:%s:%s", workflowID, documentID)
    token, err := e.locker.Lock(lockKey, 30*time.Second)
    if err != nil {
        return nil, fmt.Errorf("failed to acquire lock: %w", err)
    }
    defer e.locker.Unlock(lockKey, token)
    
    // 1. 加载工作流定义
    workflow, err := e.workflowRepo.GetByID(workflowID)
    if err != nil {
        return nil, err
    }
    
    // 2. 创建流程实例
    instance := &WorkflowInstance{
        ID:           generateID(),
        WorkflowID:   workflowID,
        DocumentID:   documentID,
        Status:       StatusRunning,
        Initiator:    initiator,
        Variables:    variables,
        CreatedAt:    time.Now(),
    }
    
    // 3. 找到开始节点
    startNode := workflow.FindStartNode()
    instance.CurrentNode = startNode.ID
    
    // 4. 保存实例
    if err := e.instanceRepo.Save(instance); err != nil {
        return nil, err
    }
    
    // 5. 自动流转到第一个任务节点
    if err := e.moveToNextNode(instance, workflow); err != nil {
        return nil, err
    }
    
    // 6. 发布事件
    e.eventBus.Publish(ProcessStartedEvent{InstanceID: instance.ID})
    
    return instance, nil
}

func (e *ExecutionEngine) CompleteTask(taskID string, userID string, action string, comment string, variables map[string]interface{}) error {
    // 幂等性检查
    lockKey := fmt.Sprintf("complete_task:%s", taskID)
    token, err := e.locker.Lock(lockKey, 30*time.Second)
    if err != nil {
        return fmt.Errorf("failed to acquire lock: %w", err)
    }
    defer e.locker.Unlock(lockKey, token)
    
    // 1. 加载任务
    task, err := e.taskRepo.GetByID(taskID)
    if err != nil {
        return err
    }
    
    // 检查任务是否已完成（幂等）
    if task.Status == TaskCompleted {
        return nil // 已完成，直接返回
    }
    
    // 2. 验证权限
    if task.Assignee != userID {
        return errors.New("permission denied")
    }
    
    // 3. 更新任务状态
    task.Status = TaskCompleted
    task.Action = action
    task.Comment = comment
    task.CompletedAt = timePtr(time.Now())
    if err := e.taskRepo.Update(task); err != nil {
        return err
    }
    
    // 4. 记录历史
    e.historyRepo.Add(&HistoryLog{
        InstanceID: task.InstanceID,
        TaskID:     taskID,
        Action:     action,
        Operator:   userID,
        Comment:    comment,
        Timestamp:  time.Now(),
    })
    
    // 5. 合并变量
    instance, _ := e.instanceRepo.GetByID(task.InstanceID)
    for k, v := range variables {
        instance.Variables[k] = v
    }
    
    // 6. 检查是否可以流转到下一节点
    if e.canMoveToNext(task, instance) {
        workflow, _ := e.workflowRepo.GetByID(instance.WorkflowID)
        if err := e.moveToNextNode(instance, workflow); err != nil {
            return err
        }
    }
    
    // 7. 发布事件
    e.eventBus.Publish(TaskCompletedEvent{TaskID: taskID, Action: action})
    
    return nil
}
```

### 5.3 任务管理

```go
type TaskManager struct {
    taskRepo    TaskRepository
    userService UserService
    notifier    Notifier
}

// 获取用户待办任务
func (tm *TaskManager) GetPendingTasks(userID string) ([]*Task, error) {
    // 获取直接分配的任务
    directTasks, err := tm.taskRepo.FindByAssignee(userID, TaskPending)
    if err != nil {
        return nil, err
    }
    
    // 获取角色任务（用户拥有的角色对应的任务）
    roles := tm.userService.GetUserRoles(userID)
    roleTasks, err := tm.taskRepo.FindByRoles(roles, TaskPending)
    if err != nil {
        return nil, err
    }
    
    return append(directTasks, roleTasks...), nil
}

// 任务委托
func (tm *TaskManager) DelegateTask(taskID string, fromUser string, toUser string, comment string) error {
    task, err := tm.taskRepo.GetByID(taskID)
    if err != nil {
        return err
    }
    
    if task.Assignee != fromUser {
        return errors.New("permission denied")
    }
    
    // 创建委托记录
    delegation := &TaskDelegation{
        TaskID:      taskID,
        FromUser:    fromUser,
        ToUser:      toUser,
        Comment:     comment,
        CreatedAt:   time.Now(),
    }
    
    // 更新任务处理人
    task.Assignee = toUser
    task.Status = TaskDelegated
    
    tm.taskRepo.Update(task)
    tm.notifier.Notify(toUser, "您有新的委托任务")
    
    return nil
}
```

### 5.4 审批人解析

```go
type AssigneeResolver interface {
    Resolve(rule *AssigneeRule, instance *WorkflowInstance) ([]string, error)
}

type DefaultAssigneeResolver struct {
    userService UserService
    orgService  OrganizationService
    exprEngine  ExpressionEngine
}

func (r *DefaultAssigneeResolver) Resolve(rule *AssigneeRule, instance *WorkflowInstance) ([]string, error) {
    switch rule.Type {
    case "User":
        return rule.Values, nil
        
    case "Role":
        var users []string
        for _, roleID := range rule.Values {
            roleUsers := r.userService.GetUsersByRole(roleID)
            users = append(users, roleUsers...)
        }
        return users, nil
        
    case "Department":
        var users []string
        for _, deptID := range rule.Values {
            deptUsers := r.orgService.GetUsersByDepartment(deptID)
            users = append(users, deptUsers...)
        }
        return users, nil
        
    case "Expression":
        // 动态表达式，例如：$initiator.manager
        result := r.exprEngine.Evaluate(rule.Expression, instance.Variables)
        return result.([]string), nil
        
    case "Dynamic":
        // 从单据中获取审批人
        // 例如：单据的"申请人上级"字段
        return r.resolveFromDocument(rule.Expression, instance.DocumentID)
        
    default:
        return nil, fmt.Errorf("unknown assignee type: %s", rule.Type)
    }
}
```

### 5.5 条件评估引擎

```go
type ConditionEvaluator interface {
    Evaluate(condition string, variables map[string]interface{}) (bool, error)
}

type ExpressionEvaluator struct {
    // 可使用 github.com/antonmedv/expr 等表达式库
}

func (e *ExpressionEvaluator) Evaluate(condition string, variables map[string]interface{}) (bool, error) {
    if condition == "" {
        return true, nil // 无条件默认为true
    }
    
    // 示例条件：
    // - "amount > 10000"
    // - "department == 'IT' && amount < 5000"
    // - "approvalCount >= requiredCount"
    
    program, err := expr.Compile(condition, expr.Env(variables))
    if err != nil {
        return false, err
    }
    
    result, err := expr.Run(program, variables)
    if err != nil {
        return false, err
    }
    
    return result.(bool), nil
}
```

---

## 六、特殊场景处理

### 6.1 条件分支

根据流程变量自动选择路由：

```go
// 条件节点处理
func (e *ExecutionEngine) handleConditionNode(node *Node, instance *WorkflowInstance, workflow *Workflow) error {
    // 获取所有出口转换
    transitions := workflow.GetTransitionsFrom(node.ID)
    
    // 按优先级排序
    sort.Slice(transitions, func(i, j int) bool {
        return transitions[i].Priority > transitions[j].Priority
    })
    
    // 评估条件，选择第一个满足的分支
    for _, trans := range transitions {
        matched, err := e.conditionEval.Evaluate(trans.Condition, instance.Variables)
        if err != nil {
            continue
        }
        if matched {
            instance.CurrentNode = trans.To
            return e.instanceRepo.Update(instance)
        }
    }
    
    return errors.New("no matching condition branch")
}
```

### 6.2 并行网关（会签）

```go
// 并行网关 - 分支
func (e *ExecutionEngine) handleParallelGatewaySplit(node *Node, instance *WorkflowInstance) error {
    transitions := e.workflow.GetTransitionsFrom(node.ID)
    
    // 为每个分支创建子流程token
    for _, trans := range transitions {
        token := &ProcessToken{
            InstanceID: instance.ID,
            NodeID:     trans.To,
            Status:     "active",
        }
        e.tokenRepo.Save(token)
        
        // 触发分支节点的任务创建
        e.processNode(trans.To, instance)
    }
    
    return nil
}

// 并行网关 - 合并
func (e *ExecutionEngine) handleParallelGatewayJoin(node *Node, instance *WorkflowInstance) error {
    // 检查所有前驱节点的token是否都已完成
    predecessors := e.workflow.GetPredecessors(node.ID)
    
    allCompleted := true
    for _, predID := range predecessors {
        tokens := e.tokenRepo.FindByNodeID(predID)
        for _, token := range tokens {
            if token.Status != "completed" {
                allCompleted = false
                break
            }
        }
    }
    
    if allCompleted {
        // 所有分支完成，继续向下流转
        return e.moveToNextNode(instance, e.workflow)
    }
    
    // 等待其他分支完成
    return nil
}
```

### 6.3 子流程

```go
type SubProcessNode struct {
    Node
    SubWorkflowID string
}

func (e *ExecutionEngine) handleSubProcess(node *SubProcessNode, instance *WorkflowInstance) error {
    // 启动子流程
    subInstance, err := e.StartProcess(
        node.SubWorkflowID,
        instance.DocumentID,
        instance.Initiator,
        instance.Variables,
    )
    if err != nil {
        return err
    }
    
    // 建立父子关系
    relation := &ProcessRelation{
        ParentInstanceID: instance.ID,
        ChildInstanceID:  subInstance.ID,
        Type:            "subprocess",
    }
    e.relationRepo.Save(relation)
    
    // 监听子流程完成事件
    e.eventBus.Subscribe(ProcessCompletedEvent{}, func(event Event) {
        if event.InstanceID == subInstance.ID {
            // 子流程完成，继续父流程
            e.moveToNextNode(instance, e.workflow)
        }
    })
    
    return nil
}
```

### 6.4 超时处理

```go
type TimeoutMonitor struct {
    taskRepo    TaskRepository
    instanceRepo InstanceRepository
    engine      *ExecutionEngine
}

func (tm *TimeoutMonitor) Start() {
    ticker := time.NewTicker(1 * time.Minute)
    
    go func() {
        for range ticker.C {
            tm.checkTimeouts()
        }
    }()
}

func (tm *TimeoutMonitor) checkTimeouts() {
    // 查找超时任务
    tasks := tm.taskRepo.FindOverdueTasks()
    
    for _, task := range tasks {
        instance, _ := tm.instanceRepo.GetByID(task.InstanceID)
        workflow, _ := tm.engine.workflowRepo.GetByID(instance.WorkflowID)
        node := workflow.FindNode(task.NodeID)
        
        if node.TimeoutRule == nil {
            continue
        }
        
        switch node.TimeoutRule.Action {
        case "Auto":
            // 自动通过
            tm.engine.CompleteTask(task.ID, "system", "auto_approve", "超时自动通过", nil)
            
        case "Escalate":
            // 升级到上级
            target := tm.resolveEscalationTarget(task, node.TimeoutRule.Target)
            tm.engine.taskRepo.Update(&Task{
                ID:       task.ID,
                Assignee: target,
            })
            
        case "Notify":
            // 发送催办通知
            tm.sendReminderNotification(task)
        }
    }
}
```

### 6.5 加签/减签

```go
// 加签（增加审批人）
func (e *ExecutionEngine) AddSign(taskID string, users []string, mode string) error {
    task, _ := e.taskRepo.GetByID(taskID)
    instance, _ := e.instanceRepo.GetByID(task.InstanceID)
    
    // 创建加签任务
    for _, userID := range users {
        signTask := &Task{
            ID:         generateID(),
            InstanceID: instance.ID,
            NodeID:     task.NodeID + "_sign",
            Assignee:   userID,
            Status:     TaskPending,
            CreatedAt:  time.Now(),
        }
        e.taskRepo.Save(signTask)
    }
    
    // 根据模式决定流转逻辑
    if mode == "before" {
        // 前加签：先处理加签任务，再处理原任务
        task.Status = TaskPending // 保持挂起
    } else if mode == "after" {
        // 后加签：先处理原任务，再处理加签任务
        // 原任务正常流转
    }
    
    return nil
}
```

### 6.6 退回/回退

```go
func (e *ExecutionEngine) Rollback(instanceID string, targetNodeID string, reason string) error {
    instance, err := e.instanceRepo.GetByID(instanceID)
    if err != nil {
        return err
    }
    
    // 取消当前节点的所有待办任务
    currentTasks := e.taskRepo.FindByNodeAndInstance(instance.CurrentNode, instanceID)
    for _, task := range currentTasks {
        task.Status = TaskCanceled
        e.taskRepo.Update(task)
    }
    
    // 回退到目标节点
    instance.CurrentNode = targetNodeID
    e.instanceRepo.Update(instance)
    
    // 重新创建目标节点的任务
    workflow, _ := e.workflowRepo.GetByID(instance.WorkflowID)
    targetNode := workflow.FindNode(targetNodeID)
    e.createTasksForNode(targetNode, instance)
    
    // 记录历史
    e.historyRepo.Add(&HistoryLog{
        InstanceID: instanceID,
        Action:     "rollback",
        Comment:    reason,
        Timestamp:  time.Now(),
    })
    
    return nil
}
```

---

## 七、数据库设计

### 7.1 表结构

```sql
-- 工作流定义表
CREATE TABLE workflows (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    version VARCHAR(20) NOT NULL,
    description TEXT,
    definition_json TEXT NOT NULL,  -- 完整的流程定义JSON
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_name_version (name, version)
);

-- 流程实例表
CREATE TABLE workflow_instances (
    id VARCHAR(50) PRIMARY KEY,
    workflow_id VARCHAR(50) NOT NULL,
    workflow_name VARCHAR(200),
    document_id VARCHAR(50) NOT NULL,
    document_type VARCHAR(50),
    current_node VARCHAR(50),
    status VARCHAR(20) NOT NULL,
    variables TEXT,  -- JSON格式
    initiator VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    INDEX idx_workflow_id (workflow_id),
    INDEX idx_document_id (document_id),
    INDEX idx_status (status),
    INDEX idx_initiator (initiator)
);

-- 任务表
CREATE TABLE tasks (
    id VARCHAR(50) PRIMARY KEY,
    instance_id VARCHAR(50) NOT NULL,
    node_id VARCHAR(50) NOT NULL,
    node_name VARCHAR(200),
    assignee VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    action VARCHAR(50),
    comment TEXT,
    variables TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    due_date TIMESTAMP NULL,
    INDEX idx_instance_id (instance_id),
    INDEX idx_assignee_status (assignee, status),
    INDEX idx_due_date (due_date)
);

-- 历史记录表
CREATE TABLE history_logs (
    id VARCHAR(50) PRIMARY KEY,
    instance_id VARCHAR(50) NOT NULL,
    task_id VARCHAR(50),
    node_id VARCHAR(50),
    node_name VARCHAR(200),
    action VARCHAR(50),
    operator VARCHAR(50),
    operator_name VARCHAR(100),
    comment TEXT,
    variables TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_instance_id (instance_id),
    INDEX idx_task_id (task_id)
);

-- 流程令牌表（用于并行网关）
CREATE TABLE process_tokens (
    id VARCHAR(50) PRIMARY KEY,
    instance_id VARCHAR(50) NOT NULL,
    node_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    INDEX idx_instance_node (instance_id, node_id)
);

-- 任务委托表
CREATE TABLE task_delegations (
    id VARCHAR(50) PRIMARY KEY,
    task_id VARCHAR(50) NOT NULL,
    from_user VARCHAR(50) NOT NULL,
    to_user VARCHAR(50) NOT NULL,
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_task_id (task_id)
);

-- 流程关系表（用于子流程）
CREATE TABLE process_relations (
    id VARCHAR(50) PRIMARY KEY,
    parent_instance_id VARCHAR(50) NOT NULL,
    child_instance_id VARCHAR(50) NOT NULL,
    relation_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_parent (parent_instance_id),
    INDEX idx_child (child_instance_id)
);
```

---

## 八、事件驱动架构

```go
type Event interface {
    GetType() string
    GetTimestamp() time.Time
}

type ProcessStartedEvent struct {
    InstanceID string
    WorkflowID string
    Initiator  string
    Timestamp  time.Time
}

type TaskCreatedEvent struct {
    TaskID     string
    InstanceID string
    Assignee   string
    Timestamp  time.Time
}

type TaskCompletedEvent struct {
    TaskID     string
    InstanceID string
    Action     string
    Operator   string
    Timestamp  time.Time
}

type ProcessCompletedEvent struct {
    InstanceID string
    Status     string
    Timestamp  time.Time
}

type EventBus interface {
    Publish(event Event)
    Subscribe(eventType string, handler EventHandler)
}

type EventHandler func(event Event)

// 事件处理器示例
func setupEventHandlers(bus EventBus, notifier Notifier) {
    // 任务创建时发送通知
    bus.Subscribe("TaskCreated", func(event Event) {
        e := event.(TaskCreatedEvent)
        notifier.Notify(e.Assignee, "您有新的待办任务")
    })
    
    // 流程完成时更新单据状态
    bus.Subscribe("ProcessCompleted", func(event Event) {
        e := event.(ProcessCompletedEvent)
        // 更新业务单据状态
        updateDocumentStatus(e.InstanceID, e.Status)
    })
    
    // 任务超时时发送提醒
    bus.Subscribe("TaskOverdue", func(event Event) {
        e := event.(TaskOverdueEvent)
        notifier.Notify(e.Assignee, "您有任务即将超时")
    })
}
```

---

## 九、API接口设计

```go
// REST API路由
type WorkflowAPI struct {
    engine *ExecutionEngine
}

// POST /api/workflows - 部署工作流
func (api *WorkflowAPI) DeployWorkflow(c *gin.Context) {
    var workflow Workflow
    if err := c.BindJSON(&workflow); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    if err := api.engine.DeployWorkflow(&workflow); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, workflow)
}

// POST /api/processes - 启动流程
func (api *WorkflowAPI) StartProcess(c *gin.Context) {
    var req struct {
        WorkflowID string                 `json:"workflow_id"`
        DocumentID string                 `json:"document_id"`
        Variables  map[string]interface{} `json:"variables"`
    }
    
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    userID := c.GetString("user_id") // 从JWT中获取
    
    instance, err := api.engine.StartProcess(req.WorkflowID, req.DocumentID, userID, req.Variables)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, instance)
}

// POST /api/tasks/:id/complete - 完成任务
func (api *WorkflowAPI) CompleteTask(c *gin.Context) {
    taskID := c.Param("id")
    userID := c.GetString("user_id")
    
    var req struct {
        Action    string                 `json:"action"`
        Comment   string                 `json:"comment"`
        Variables map[string]interface{} `json:"variables"`
    }
    
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    if err := api.engine.CompleteTask(taskID, userID, req.Action, req.Comment, req.Variables); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"message": "success"})
}

// GET /api/tasks/pending - 获取待办任务
func (api *WorkflowAPI) GetPendingTasks(c *gin.Context) {
    userID := c.GetString("user_id")
    
    tasks, err := api.engine.GetPendingTasks(userID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, tasks)
}

// GET /api/processes/:id - 查询流程实例
func (api *WorkflowAPI) GetInstance(c *gin.Context) {
    instanceID := c.Param("id")
    
    instance, err := api.engine.GetInstance(instanceID)
    if err != nil {
        c.JSON(404, gin.H{"error": "instance not found"})
        return
    }
    
    c.JSON(200, instance)
}

// GET /api/processes/:id/history - 查询流程历史
func (api *WorkflowAPI) GetHistory(c *gin.Context) {
    instanceID := c.Param("id")
    
    history, err := api.engine.GetHistory(instanceID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, history)
}

// POST /api/processes/:id/rollback - 流程回退
func (api *WorkflowAPI) Rollback(c *gin.Context) {
    instanceID := c.Param("id")
    
    var req struct {
        TargetNodeID string `json:"target_node_id"`
        Reason       string `json:"reason"`
    }
    
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    if err := api.engine.Rollback(instanceID, req.TargetNodeID, req.Reason); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"message": "success"})
}
```

---

## 十、性能优化

### 10.1 缓存策略

```go
type CachedWorkflowEngine struct {
    engine *ExecutionEngine
    cache  *redis.Client
}

func (c *CachedWorkflowEngine) GetWorkflow(workflowID string) (*Workflow, error) {
    // 先从缓存读取
    cacheKey := fmt.Sprintf("workflow:%s", workflowID)
    cached, err := c.cache.Get(cacheKey).Result()
    if err == nil {
        var workflow Workflow
        json.Unmarshal([]byte(cached), &workflow)
        return &workflow, nil
    }
    
    // 缓存未命中，从数据库读取
    workflow, err := c.engine.GetWorkflow(workflowID)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存
    data, _ := json.Marshal(workflow)
    c.cache.Set(cacheKey, data, 1*time.Hour)
    
    return workflow, nil
}
```

### 10.2 异步处理

```go
type AsyncExecutionEngine struct {
    engine    *ExecutionEngine
    taskQueue chan Task
}

func (a *AsyncExecutionEngine) CompleteTask(taskID string, userID string, action string, comment string, variables map[string]interface{}) error {
    // 快速响应用户
    task := Task{
        ID:        taskID,
        // ... 其他字段
    }
    
    // 异步处理
    a.taskQueue <- task
    
    return nil
}

func (a *AsyncExecutionEngine) startWorkers(workerCount int) {
    for i := 0; i < workerCount; i++ {
        go func() {
            for task := range a.taskQueue {
                a.engine.processTask(&task)
            }
        }()
    }
}
```

### 10.3 批量查询优化

```go
func (r *TaskRepository) FindByAssignees(userIDs []string, status TaskStatus) ([]*Task, error) {
    query := `
        SELECT * FROM tasks 
        WHERE assignee IN (?) AND status = ?
        ORDER BY created_at DESC
    `
    
    var tasks []*Task
    err := r.db.Select(&tasks, query, userIDs, status)
    return tasks, err
}
```

---

## 十一、监控与运维

### 11.1 监控指标

- **流程指标**：启动数、完成数、拒绝数、平均耗时
- **任务指标**：待办数、超时数、完成率
- **性能指标**：响应时间、吞吐量、并发数
- **异常指标**：错误率、超时率

### 11.2 日志记录

```go
type WorkflowLogger struct {
    logger *zap.Logger
}

func (l *WorkflowLogger) LogProcessStart(instance *WorkflowInstance) {
    l.logger.Info("Process started",
        zap.String("instance_id", instance.ID),
        zap.String("workflow_id", instance.WorkflowID),
        zap.String("initiator", instance.Initiator),
    )
}

func (l *WorkflowLogger) LogTaskComplete(task *Task) {
    l.logger.Info("Task completed",
        zap.String("task_id", task.ID),
        zap.String("action", task.Action),
        zap.Duration("duration", time.Since(task.CreatedAt)),
    )
}
```

### 11.3 健康检查

```go
func (api *WorkflowAPI) HealthCheck(c *gin.Context) {
    // 检查数据库连接
    if err := api.engine.instanceRepo.Ping(); err != nil {
        c.JSON(500, gin.H{"status": "unhealthy", "error": "database connection failed"})
        return
    }
    
    // 检查待处理任务队列
    queueSize := len(api.engine.taskQueue)
    if queueSize > 10000 {
        c.JSON(500, gin.H{"status": "unhealthy", "error": "task queue overflow"})
        return
    }
    
    c.JSON(200, gin.H{
        "status": "healthy",
        "queue_size": queueSize,
    })
}
```

---

## 十二、扩展性设计

### 12.1 插件系统

```go
type Plugin interface {
    Name() string
    OnProcessStart(instance *WorkflowInstance) error
    OnTaskCreate(task *Task) error
    OnTaskComplete(task *Task) error
    OnProcessComplete(instance *WorkflowInstance) error
}

type PluginManager struct {
    plugins []Plugin
}

func (pm *PluginManager) Register(plugin Plugin) {
    pm.plugins = append(pm.plugins, plugin)
}

func (pm *PluginManager) ExecuteHook(hookName string, data interface{}) {
    for _, plugin := range pm.plugins {
        switch hookName {
        case "OnProcessStart":
            plugin.OnProcessStart(data.(*WorkflowInstance))
        case "OnTaskComplete":
            plugin.OnTaskComplete(data.(*Task))
        }
    }
}
```

### 12.2 自定义节点类型

```go
type CustomNodeHandler interface {
    Handle(node *Node, instance *WorkflowInstance) error
}

type NodeHandlerRegistry struct {
    handlers map[NodeType]CustomNodeHandler
}

func (r *NodeHandlerRegistry) Register(nodeType NodeType, handler CustomNodeHandler) {
    r.handlers[nodeType] = handler
}

// 使用示例：添加"发送邮件"节点
type EmailNodeHandler struct {
    emailService EmailService
}

func (h *EmailNodeHandler) Handle(node *Node, instance *WorkflowInstance) error {
    // 从节点配置中获取收件人、主题等
    to := node.Config.CustomFields["to"]
    subject := node.Config.CustomFields["subject"]
    
    return h.emailService.Send(to, subject, "流程通知")
}
```

---

## 十三、安全性

### 13.1 权限控制

```go
type PermissionChecker interface {
    CanStartProcess(userID string, workflowID string) bool
    CanCompleteTask(userID string, taskID string) bool
    CanViewInstance(userID string, instanceID string) bool
    CanRollback(userID string, instanceID string) bool
}

type RBACPermissionChecker struct {
    roleService RoleService
}

func (c *RBACPermissionChecker) CanCompleteTask(userID string, taskID string) bool {
    task, _ := getTask(taskID)
    
    // 检查是否是任务处理人
    if task.Assignee == userID {
        return true
    }
    
    // 检查是否有管理员权限
    if c.roleService.HasRole(userID, "admin") {
        return true
    }
    
    return false
}
```

### 13.2 数据加密

```go
// 敏感变量加密存储
type EncryptedVariableStore struct {
    cipher cipher.AEAD
}

func (s *EncryptedVariableStore) Set(key string, value interface{}) error {
    data, _ := json.Marshal(value)
    encrypted := s.encrypt(data)
    return s.store(key, encrypted)
}

func (s *EncryptedVariableStore) Get(key string) (interface{}, error) {
    encrypted, _ := s.load(key)
    decrypted := s.decrypt(encrypted)
    
    var value interface{}
    json.Unmarshal(decrypted, &value)
    return value, nil
}
```

---

## 十四、错误码定义

### 14.1 错误码规范

```go
// 错误码格式：WXXYYYY
// WX: 模块代码（10-流程定义, 20-流程实例, 30-任务, 40-系统）
// YYYY: 具体错误编号

type ErrorCode int

const (
    // 通用错误 (40xxxx)
    ErrSuccess           ErrorCode = 0
    ErrUnknown           ErrorCode = 400000
    ErrInvalidParam      ErrorCode = 400001
    ErrUnauthorized      ErrorCode = 400002
    ErrForbidden         ErrorCode = 400003
    ErrNotFound          ErrorCode = 400004
    ErrTimeout           ErrorCode = 400005
    ErrRateLimited       ErrorCode = 400006
    ErrDatabaseError     ErrorCode = 400007
    ErrLockFailed        ErrorCode = 400008
    
    // 流程定义错误 (10xxxx)
    ErrWorkflowNotFound      ErrorCode = 100001
    ErrWorkflowInactive      ErrorCode = 100002
    ErrWorkflowVersionExists ErrorCode = 100003
    ErrInvalidWorkflowDef    ErrorCode = 100004
    ErrNodeNotFound          ErrorCode = 100005
    ErrTransitionNotFound    ErrorCode = 100006
    ErrInvalidNodeType       ErrorCode = 100007
    ErrCyclicWorkflow        ErrorCode = 100008  // 循环依赖
    ErrNoStartNode           ErrorCode = 100009
    ErrNoEndNode             ErrorCode = 100010
    ErrDisconnectedNode      ErrorCode = 100011  // 孤立节点
    
    // 流程实例错误 (20xxxx)
    ErrInstanceNotFound      ErrorCode = 200001
    ErrInstanceCompleted     ErrorCode = 200002
    ErrInstanceSuspended     ErrorCode = 200003
    ErrInstanceTerminated    ErrorCode = 200004
    ErrNoValidTransition     ErrorCode = 200005
    ErrConditionEvalFailed   ErrorCode = 200006
    ErrVariableMissing       ErrorCode = 200007
    ErrInvalidInstanceStatus ErrorCode = 200008
    ErrRollbackNotAllowed    ErrorCode = 200009
    
    // 任务错误 (30xxxx)
    ErrTaskNotFound        ErrorCode = 300001
    ErrTaskCompleted       ErrorCode = 300002
    ErrTaskCanceled        ErrorCode = 300003
    ErrTaskAssigneeMismatch ErrorCode = 300004
    ErrTaskDelegateNotAllowed ErrorCode = 300005
    ErrCommentRequired     ErrorCode = 300006
    ErrInvalidAction       ErrorCode = 300007
    ErrTaskOverdue         ErrorCode = 300008
)

// WorkflowError 工作流错误
type WorkflowError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Detail  string    `json:"detail,omitempty"`
}

func (e *WorkflowError) Error() string {
    return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
}

// 预定义错误
var (
    ErrLockFailedError = &WorkflowError{Code: ErrLockFailed, Message: "获取分布式锁失败"}
    ErrWorkflowNotFoundError = &WorkflowError{Code: ErrWorkflowNotFound, Message: "工作流不存在"}
    ErrTaskNotFoundError = &WorkflowError{Code: ErrTaskNotFound, Message: "任务不存在"}
)
```

### 14.2 错误处理最佳实践

```go
// 统一错误响应
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    TraceID string      `json:"trace_id,omitempty"`
}

// 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            var wfErr *WorkflowError
            if errors.As(err, &wfErr) {
                c.JSON(http.StatusOK, APIResponse{
                    Code:    int(wfErr.Code),
                    Message: wfErr.Message,
                    TraceID: c.GetString("trace_id"),
                })
                return
            }
            
            // 未知错误
            c.JSON(http.StatusInternalServerError, APIResponse{
                Code:    int(ErrUnknown),
                Message: "服务器内部错误",
                TraceID: c.GetString("trace_id"),
            })
        }
    }
}
```

---

## 十五、配置管理

### 15.1 配置结构

```go
// Config 工作流引擎配置
type Config struct {
    // 服务配置
    Server ServerConfig `yaml:"server"`
    
    // 数据库配置
    Database DatabaseConfig `yaml:"database"`
    
    // Redis配置
    Redis RedisConfig `yaml:"redis"`
    
    // 消息队列配置
    MessageQueue MQConfig `yaml:"message_queue"`
    
    // 引擎配置
    Engine EngineConfig `yaml:"engine"`
    
    // 日志配置
    Log LogConfig `yaml:"log"`
}

type ServerConfig struct {
    Port         int           `yaml:"port" default:"8080"`
    ReadTimeout  time.Duration `yaml:"read_timeout" default:"30s"`
    WriteTimeout time.Duration `yaml:"write_timeout" default:"30s"`
    Mode         string        `yaml:"mode" default:"release"` // debug/release
}

type DatabaseConfig struct {
    Driver          string        `yaml:"driver" default:"mysql"`
    DSN             string        `yaml:"dsn"`
    MaxOpenConns    int           `yaml:"max_open_conns" default:"100"`
    MaxIdleConns    int           `yaml:"max_idle_conns" default:"10"`
    ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" default:"1h"`
}

type RedisConfig struct {
    Addr         string        `yaml:"addr" default:"localhost:6379"`
    Password     string        `yaml:"password"`
    DB           int           `yaml:"db" default:"0"`
    PoolSize     int           `yaml:"pool_size" default:"100"`
    MinIdleConns int           `yaml:"min_idle_conns" default:"10"`
    DialTimeout  time.Duration `yaml:"dial_timeout" default:"5s"`
}

type MQConfig struct {
    Type     string `yaml:"type" default:"kafka"` // kafka/rabbitmq
    Brokers  string `yaml:"brokers"`
    Topic    string `yaml:"topic" default:"workflow_events"`
    GroupID  string `yaml:"group_id" default:"workflow_engine"`
}

type EngineConfig struct {
    // 任务超时检查间隔
    TimeoutCheckInterval time.Duration `yaml:"timeout_check_interval" default:"1m"`
    
    // 任务默认超时时间
    DefaultTaskTimeout time.Duration `yaml:"default_task_timeout" default:"24h"`
    
    // 工作协程数
    WorkerCount int `yaml:"worker_count" default:"10"`
    
    // 任务队列大小
    TaskQueueSize int `yaml:"task_queue_size" default:"1000"`
    
    // 是否启用异步处理
    EnableAsync bool `yaml:"enable_async" default:"true"`
    
    // 分布式锁过期时间
    LockTTL time.Duration `yaml:"lock_ttl" default:"30s"`
    
    // 流程定义缓存时间
    WorkflowCacheTTL time.Duration `yaml:"workflow_cache_ttl" default:"1h"`
}

type LogConfig struct {
    Level      string `yaml:"level" default:"info"`
    Format     string `yaml:"format" default:"json"` // json/text
    OutputPath string `yaml:"output_path" default:"stdout"`
}
```

### 15.2 配置文件示例

```yaml
# config.yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  mode: release

database:
  driver: mysql
  dsn: "user:password@tcp(localhost:3306)/workflow?charset=utf8mb4&parseTime=True&loc=Local"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 1h

redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 100

message_queue:
  type: kafka
  brokers: localhost:9092
  topic: workflow_events
  group_id: workflow_engine

engine:
  timeout_check_interval: 1m
  default_task_timeout: 24h
  worker_count: 10
  task_queue_size: 1000
  enable_async: true
  lock_ttl: 30s
  workflow_cache_ttl: 1h

log:
  level: info
  format: json
  output_path: /var/log/workflow/app.log
```

### 15.3 环境变量覆盖

```go
// 支持通过环境变量覆盖配置
// 格式：WORKFLOW_${SECTION}_${KEY}
// 例如：WORKFLOW_DATABASE_DSN, WORKFLOW_SERVER_PORT

func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.SetEnvPrefix("WORKFLOW")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

---

## 十六、流程定义验证

### 16.1 验证器接口

```go
// WorkflowValidator 流程定义验证器
type WorkflowValidator interface {
    Validate(workflow *Workflow) []ValidationError
}

// ValidationError 验证错误
type ValidationError struct {
    Field   string `json:"field"`
    Code    string `json:"code"`
    Message string `json:"message"`
}

// DefaultWorkflowValidator 默认验证器实现
type DefaultWorkflowValidator struct{}

func (v *DefaultWorkflowValidator) Validate(workflow *Workflow) []ValidationError {
    var errors []ValidationError
    
    // 1. 基础验证
    if workflow.Name == "" {
        errors = append(errors, ValidationError{
            Field:   "name",
            Code:    "required",
            Message: "工作流名称不能为空",
        })
    }
    
    // 2. 节点验证
    errors = append(errors, v.validateNodes(workflow)...)
    
    // 3. 转换验证
    errors = append(errors, v.validateTransitions(workflow)...)
    
    // 4. 结构验证
    errors = append(errors, v.validateStructure(workflow)...)
    
    return errors
}

func (v *DefaultWorkflowValidator) validateNodes(workflow *Workflow) []ValidationError {
    var errors []ValidationError
    
    hasStart := false
    hasEnd := false
    nodeIDs := make(map[string]bool)
    
    for i, node := range workflow.Nodes {
        // 检查节点ID唯一性
        if nodeIDs[node.ID] {
            errors = append(errors, ValidationError{
                Field:   fmt.Sprintf("nodes[%d].id", i),
                Code:    "duplicate",
                Message: fmt.Sprintf("节点ID重复: %s", node.ID),
            })
        }
        nodeIDs[node.ID] = true
        
        // 检查开始节点
        if node.Type == StartNode {
            if hasStart {
                errors = append(errors, ValidationError{
                    Field:   fmt.Sprintf("nodes[%d]", i),
                    Code:    "multiple_start",
                    Message: "流程只能有一个开始节点",
                })
            }
            hasStart = true
        }
        
        // 检查结束节点
        if node.Type == EndNode {
            hasEnd = true
        }
        
        // 审批节点必须有审批人规则
        if node.Type == ApprovalNode {
            if node.Assignee.Type == "" {
                errors = append(errors, ValidationError{
                    Field:   fmt.Sprintf("nodes[%d].assignee", i),
                    Code:    "required",
                    Message: fmt.Sprintf("审批节点[%s]必须设置审批人规则", node.Name),
                })
            }
        }
    }
    
    if !hasStart {
        errors = append(errors, ValidationError{
            Field:   "nodes",
            Code:    "no_start_node",
            Message: "流程必须有一个开始节点",
        })
    }
    
    if !hasEnd {
        errors = append(errors, ValidationError{
            Field:   "nodes",
            Code:    "no_end_node",
            Message: "流程必须有至少一个结束节点",
        })
    }
    
    return errors
}

func (v *DefaultWorkflowValidator) validateStructure(workflow *Workflow) []ValidationError {
    var errors []ValidationError
    
    // 构建邻接表
    graph := make(map[string][]string)
    inDegree := make(map[string]int)
    
    for _, node := range workflow.Nodes {
        graph[node.ID] = []string{}
        inDegree[node.ID] = 0
    }
    
    for _, trans := range workflow.Transitions {
        graph[trans.From] = append(graph[trans.From], trans.To)
        inDegree[trans.To]++
    }
    
    // 检查孤立节点（除开始节点外）
    for _, node := range workflow.Nodes {
        if node.Type != StartNode && inDegree[node.ID] == 0 {
            errors = append(errors, ValidationError{
                Field:   "nodes",
                Code:    "disconnected_node",
                Message: fmt.Sprintf("节点[%s]没有入口转换", node.Name),
            })
        }
    }
    
    // 检查循环依赖（使用DFS）
    if hasCycle(graph) {
        errors = append(errors, ValidationError{
            Field:   "transitions",
            Code:    "cyclic_dependency",
            Message: "流程存在循环依赖",
        })
    }
    
    return errors
}

// hasCycle 检测图中是否有循环
func hasCycle(graph map[string][]string) bool {
    visited := make(map[string]int) // 0: 未访问, 1: 访问中, 2: 已完成
    
    var dfs func(node string) bool
    dfs = func(node string) bool {
        visited[node] = 1
        
        for _, next := range graph[node] {
            if visited[next] == 1 {
                return true // 发现环
            }
            if visited[next] == 0 && dfs(next) {
                return true
            }
        }
        
        visited[node] = 2
        return false
    }
    
    for node := range graph {
        if visited[node] == 0 {
            if dfs(node) {
                return true
            }
        }
    }
    
    return false
}
```

---

## 十七、总结

这是一个完整的工作流引擎设计方案，涵盖了：

1. **系统架构**：清晰的分层架构、核心组件说明
2. **核心功能**：流程定义、实例管理、任务管理
3. **高级特性**：条件分支、并行网关、子流程、超时处理
4. **可靠性保障**：分布式锁、幂等性设计、事件驱动
5. **完善的接口**：仓储接口、外部服务接口、分布式锁接口
6. **错误处理**：完整的错误码定义、统一错误响应
7. **配置管理**：灵活的配置结构、环境变量支持
8. **流程验证**：完整的流程定义验证器

### 实施路线图

**第一阶段（基础功能）- 2-3周**
- 数据模型和数据库设计
- 基础CRUD接口
- 简单的串行审批流程

**第二阶段（核心功能）- 3-4周**
- 条件分支和路由
- 并行网关（会签/或签）
- 任务超时处理
- 事件通知

**第三阶段（高级功能）- 2-3周**
- 子流程支持
- 加签/减签
- 流程回退
- 任务委托

**第四阶段（可靠性）- 2周**
- 分布式锁
- 幂等性保证
- 异步处理
- 性能优化

**第五阶段（运维支持）- 1-2周**
- 监控告警
- 日志追踪
- 配置管理
- 文档完善

---

**文档版本**：v2.0  
**最后更新**：2025-12-31  
**维护者**：工作流引擎团队
