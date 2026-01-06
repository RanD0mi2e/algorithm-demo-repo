import { MyPromise } from './myPromise'

export function test() {
  let p = new MyPromise((resolve, reject) => {
    setTimeout(() => {
      resolve(1)
    }, 1000);
  }).then((val: any) => {
    return new MyPromise((resolve, reject) => {
      console.log('val', val)
      setTimeout(() => {
        resolve(val + 1)
      }, 2000);
    })
  }).then((val2) => {
    console.log('val2', val2)
  })
}