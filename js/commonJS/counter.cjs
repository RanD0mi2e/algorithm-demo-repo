let count = 0
module.exports = {
  count,
  increment () {
    count++
    console.log('模块内部count: ', count)
  }
}