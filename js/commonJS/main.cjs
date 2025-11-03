const counter1 = require('./counter.cjs')
const counter2 = require('./counter.cjs')
counter1.increment()
console.log(counter1.count) // 仍然为0
console.log(counter2.count) // 仍然为0

// 造成这个问题的原因：
// require导入是值传递,而increment修改的是counter模块内count的值，所以counter1，counter2内部的count还是最初导入时的0值