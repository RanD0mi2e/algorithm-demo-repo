// 6. 优化类型
type Status = 'pending' | 'fulfilled' | 'rejected';
type Executor<T> = (resolve: Resolve<T>, reject: Reject) => void;
type Resolve<T> = (result?: T | Thenable<T>) => void;
type Reject = (reason?: any) => void;
type Thenable<T> = {
  then<TResult1 = T, TResult2 = never>(
    onFulfilled?: ((result: T) => TResult1 | Thenable<TResult1>) | undefined | null,
    onRejected?: ((reason: any) => TResult2 | Thenable<TResult2>) | undefined | null
  ): any;
};

class MyPromise<T> {
  status: Status = 'pending';
  private fulfilledValue?: T | Thenable<T>;
  private rejectionReason?: any;
  private onFulfilledCbList: ((result?: T) => any)[] = [];
  private onRejectedCbList: ((reason?: any) => any)[] = [];

  constructor(executor?: Executor<T>) {
    try {
      executor?.(this.resolve, this.reject);
    } catch (error) {
      this.reject(error as any);
    }
  }

  private resolve: Resolve<T> = (result?: T | Thenable<T>) => {
    if (this.status !== 'pending') return;

    this.fulfilledValue = result;
    this.status = 'fulfilled';
    this.onFulfilledCbList.forEach((cb) => cb());
  };

  private reject: Reject = (reason?: any) => {
    if (this.status !== 'pending') return;

    this.rejectionReason = reason;
    this.status = 'rejected';
    this.onRejectedCbList.forEach((cb) => cb());
  };

  then<TResult1 = T, TResult2 = never>(
    onFulfilled?: ((result: T) => TResult1 | Thenable<TResult1>) | undefined | null,
    onRejected?: ((result: any) => TResult2 | Thenable<TResult2>) | undefined | null
  ): MyPromise<TResult1 | TResult2> {
    const promise2 = new MyPromise<TResult1 | TResult2>((resolve, reject) => {
      let hasOnFulfilled = typeof onFulfilled === 'function';
      let hasOnRejected = typeof onRejected === 'function';
      let isFulfilled = this.status === 'fulfilled';
      let isRejected = this.status === 'rejected';
      let isPending = this.status === 'pending';

      if (isPending) {
        // 处理待处理promise的成功回调
        this.onFulfilledCbList.push(() => {
          queueMicrotask(() => {
            try {
              if (hasOnFulfilled) {
                const x = onFulfilled!(this.fulfilledValue as T);
                resolvePromise2(promise2, resolve, reject, x);
              } else {
                resolve(this.fulfilledValue as TResult1 | TResult2 | Thenable<TResult1 | TResult2> | undefined); // 如果onFulfilled不是函数，则使用promise2传递调用then的promise
              }
            } catch (error) {
              reject(error);
            }
          });
        });
        // 处理待处理promise的失败回调
        this.onRejectedCbList.push(() => {
          queueMicrotask(() => {
            try {
              if (hasOnRejected) {
                const x = onRejected!(this.rejectionReason);
                resolvePromise2(promise2, resolve, reject, x);
              } else {
                reject(this.rejectionReason); // 如果onRejected不是函数，则使用promise2传递调用then的promise
              }
            } catch (error) {
              reject(error);
            }
          });
        });
      } else if (isFulfilled) {
        // 处理已成功promise的回调
        queueMicrotask(() => {
          try {
            if (hasOnFulfilled) {
              const x = onFulfilled!(this.fulfilledValue as T);
              resolvePromise2(promise2, resolve, reject, x);
            } else {
              resolve(this.fulfilledValue as TResult1 | TResult2 | Thenable<TResult1 | TResult2> | undefined); // 如果onFulfilled不是函数，则使用promise2传递调用then的promise
            }
          } catch (error) {
            reject(error);
          }
        });
      } else if (isRejected) {
        // 处理已失败promise的回调
        queueMicrotask(() => {
          try {
            if (hasOnRejected) {
              const x = onRejected!(this.rejectionReason);
              resolvePromise2(promise2, resolve, reject, x);
            } else {
              reject(this.rejectionReason); // 如果onRejected不是函数，则使用promise2传递调用then的promise
            }
          } catch (error) {
            reject(error);
          }
        });
      }
    });

    return promise2;
  }
}

function resolvePromise2<T>(promise2: MyPromise<T>, resolve: Resolve<T>, reject: Reject, x: any): void {
  if (promise2 === x) {
    /**
     * 防止这种情况出现
     * const p2 = new MyPromise((resolve) => resolve()).then(() => p2);
     */
    throw new TypeError('Chaining cycle detected for promise');
  }

  if (x instanceof MyPromise) {
    // Promise对象
    x.then((y) => resolvePromise2(promise2, resolve, reject, y), reject); // 如果内层Promise是成功状态，则继续递归
  } else if (isThisType(x, 'Function') || isThisType(x, 'Object')) {
    let then: any = null;
    try {
      then = x.then;
    } catch (error) {
      reject(error);
    }
    if (isThisType(then, 'Function')) {
      // thenable对象
      let called = false;
      try {
        then!.call(
          x,
          (y: any) => {
            if (called) return;
            called = true;
            resolvePromise2(promise2, resolve, reject, y); // 如果内层Promise是成功状态，则继续递归
          }, // resolvePromise
          (z: any) => {
            if (called) return;
            called = true;
            reject(z);
          } // rejectPromise
        );
      } catch (error) {
        if (called) return;
        called = true;
        reject(error);
      }
    } else {
      resolve(x);
    }
  } else {
    // 普通值
    resolve(x);
  }
}

function isThisType(x: any, type: string): boolean {
  return Object.prototype.toString.call(x) === `[object ${type}]`;
}

const p1 = new MyPromise(); // pending

const p2 = new MyPromise((resolve) => {
  resolve();
}); // fulfilled

const p3 = new MyPromise<void>((_, reject) => {
  reject();
}); // rejected

const p4 = new MyPromise<string>((resolve, reject) => {
  setTimeout(() => {
    if (Math.random() > 0.5) {
      resolve('F');
    } else {
      reject('R');
    }
  }, 1000);
}); // pending -> fulfilled or rejected

const p5 = p2
  .then(() => {
    return 1;
  })
  .then(() => {
    return 2;
  }); // 2

const p6 = p4.then(
  (res) => {
    return 'p6-' + res;
  },
  (res) => {
    return 'p6-' + res;
  }
); // 'p6-F' or 'p6-R'

new MyPromise<string>((resolve) => resolve('micro')).then((val) => {
  console.log('then 回调（应在微任务阶段执行）', val);
});
console.log('同步日志（应先输出）');

const p7 = p2.then(() => new MyPromise((resolve) => resolve('inner promise'))); // inner promise

/**
 * trickyThenable 演示：
 * 1. 首先 resolve 一个“仍在等待的 thenable”（里面用 setTimeout 异步再 resolve 真值）
 * 2. 紧接着同步调用 reject
 *
 * - 启用 called：同步 reject 被忽略，等异步 resolve 后输出「最终结果： 最终成功」
 * - 注释 called：promise 仍处于 pending，因此同步 reject 会生效，输出「最终结果： 后来又失败」
 *
 * 为什么myPromise不需要做这种限制？
 * 因为对于Thenable对象，resolve或者reject的调用是由用户自己决定的，而myPromise中，是由myPromise类中的resolvePromise2决定的
 */
const pendingThenable: Thenable<string> = {
  then(innerResolve) {
    if (!innerResolve) return;
    setTimeout(() => innerResolve('成功'), 0);
  },
};
const trickyThenable: Thenable<string> = {
  then(resolve: any, reject: any) {
    // 先resolve了这个Thenable对象
    resolve(pendingThenable as any);
    // 后reject了这个Thenable对象（根据规范，应该忽略reject的调用）
    reject('失败');
  },
};
new MyPromise((resolve) => resolve('占位值'))
  .then(() => trickyThenable)
  .then(
    (res) => console.log('最终结果：', res),
    (err) => console.error('最终结果：', err)
  );

setTimeout(() => {
  console.log(p1, p2, p3, p4, p5, p6, p7);
}, 1000);

export {};
