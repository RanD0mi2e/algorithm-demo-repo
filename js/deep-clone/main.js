function deepClone(source, hash = new WeakMap()) {
  // 值为空
  if (source === null || source === undefined) {
    return source
  }
  // 原始值
  if (typeof source !== 'object') {
    return source
  }
  // 处理循环引用
  if (hash.has(source)) {
    return hash.get(source)
  }
  // 处理正则
  if (source instanceof RegExp) {
    return new RegExp(source.source, source.flags)
  }
  // Data
  if (source instanceof Date) {
    return new Date(source.getTime())
  }
  // Function
  if (typeof source === 'function') {
    return function(...args) {
      return source.apply(this, args)
    }
  }
  // Map
  if (source instanceof Map) {
    const cloneMap = new Map()
    hash.set(source, cloneMap)
    for (const [key, value] of source) {
      cloneMap.set(deepClone(key, hash), deepClone(value, hash))
    }
    return cloneMap
  }
  // Set
  if (source instanceof Set) {
    const cloneSet = new Set()
    hash.set(source, cloneSet)
    for (const value of source) {
      cloneSet.add(deepClone(value, hash))
    }
    return cloneSet
  }
  // ArrayBuffer
  if (source instanceof ArrayBuffer) {
    const cloned = new ArrayBuffer(source.byteLength)
    new Uint8Array(cloned).set(new Uint8Array(source))
    return cloned
  }
  // TypedArray
  if (ArrayBuffer.isView(source) && !(source instanceof DataView)) {
    const Cstr = source.constructor
    const cloned = new Cstr(source.length)
    cloned.set(source)
    return cloned
  }
  // DataView
  if (source instanceof DataView) {
    const clonedBuffer = deepClone(source.buffer, hash)
    return new DataView(clonedBuffer, source.byteOffset, source.byteLength)
  }

  let clonedObj
  // Array
  if (Array.isArray(source)) {
    clonedObj = []
  } else {
    // 保持拷贝后对象和源对象的原型链一致
    clonedObj = Object.create(Object.getPrototypeOf(source))
  }
  hash.set(source, clonedObj)
  const keys = Object.keys(source)
  for (const key of keys) {
    clonedObj[key] = deepClone(source[key], hash)
  }

  return clonedObj
}

const origin = {
  data: new Date(),
  regex: /test/i,
  map: new Map([['key1', 'val1'], ['key2', 'val2']]),
  set: new Set([1,2,3]),
  func: function(x) {console.log(x*2)},
  arr: [1,2,{a:1}],
  buffer: new ArrayBuffer(1),
  typedArr: new Int16Array([1,2,3,4]),
  dataView: new DataView(new ArrayBuffer(2))
}

const clone = deepClone(origin)

console.log('clone', clone)
console.log('event?', clone === origin)