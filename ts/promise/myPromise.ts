const Status = {
  pending: 0,
  fulfilled: 1,
  rejected: 2,
} as const;

type Status = (typeof Status)[keyof typeof Status];

export class MyPromise<T = any> {
  private state: Status;
  private value: T | undefined;
  private reason: any;
  private onFulfilledCallback: Array<(value: T) => void> = [];
  private onRejectedCallbacks: Array<(reson: any) => void> = [];

  constructor(
    executor: (
      resolve: (value: T) => void,
      reject: (reason: any) => void
    ) => void
  ) {
    this.state = Status.pending;
    try {
      executor(this.resolve.bind(this), this.reject.bind(this));
    } catch (error) {
      this.reject(error);
    }
  }

  resolve(value: T): void {
    if (this.state === Status.pending) {
      this.state = Status.fulfilled;
      this.value = value;
      this.onFulfilledCallback.forEach((fn) => fn(value));
    }
  }

  reject(reason: any): void {
    if (this.state === Status.pending) {
      this.state = Status.rejected;
      this.reason = reason;
      this.onRejectedCallbacks.forEach((fn) => fn(reason));
    }
  }

  then<TResult1 = T, TResult2 = never>(
    onFulfilled?: (value: T) => TResult1 | PromiseLike<TResult1> | null,
    onRejected?: (reason: any) => TResult2 | PromiseLike<TResult2> | null
  ): MyPromise<TResult1 | TResult2> {
    const promise2 = new MyPromise<TResult1 | TResult2>((resolve, reject) => {
      if (this.state === Status.fulfilled) {
        queueMicrotask(() => {
          try {
            if (typeof onFulfilled === "function") {
              const x = onFulfilled(this.value!);
              resolvePromise(promise2, x, resolve, reject);
            } else {
              resolve(this.value as any);
            }
          } catch (error) {
            reject(error);
          }
        });
      } else if (this.state === Status.rejected) {
        queueMicrotask(() => {
          try {
            if (typeof onRejected === "function") {
              const x = onRejected(this.reason);
              resolvePromise(promise2, x, resolve, reject);
            } else {
              reject(this.reason);
            }
          } catch (error) {
            reject(error);
          }
        });
      } else {
        // pending
        this.onFulfilledCallback.push((value) => {
          queueMicrotask(() => {
            try {
              if (typeof onFulfilled === "function") {
                const x = onFulfilled(value);
                resolvePromise(promise2, x, resolve, reject);
              } else {
                resolve(value as any);
              }
            } catch (error) {
              reject(error);
            }
          });
        });

        this.onRejectedCallbacks.push((reason) => {
          queueMicrotask(() => {
            try {
              if (typeof onRejected === "function") {
                const x = onRejected(reason);
                resolvePromise(promise2, x, resolve, reject);
              } else {
                reject(reason);
              }
            } catch (error) {
              reject(reason);
            }
          });
        });
      }
    });

    return promise2;
  }
}

function resolvePromise<T>(
  promise2: MyPromise<T>,
  x: any,
  resolve: (value: any) => void,
  reject: (reason: any) => void
): void {
  if (promise2 === x) {
    return reject(new TypeError('Chaining cycle detected for promise'));
  }

  if (x instanceof MyPromise) {
    x.then(resolve, reject)
    return
  }

  if (x !== null && (typeof x === 'object' || typeof x ==='function')) { 
    let called = false

    try {
      const then = x.then
      if (typeof then === 'function') {
        then.call(x, 
          (y: any) => {
            if (called) return
            called = true
            resolvePromise(promise2, y, resolve, reject)
          },
          (r: any) => {
            if (called) return
            called = true
            reject(r)
          }
        )
      } else {
        resolve(x)
      }
    } catch (error) {
      if (called) return;
      called = true;
      reject(error);
    }
  } else {
    resolve(x)
  }
}