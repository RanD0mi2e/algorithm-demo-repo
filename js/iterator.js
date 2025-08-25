const person = {
  name: 'Alice',
  age: 28,
  city: 'Beijing',

  [Symbol.iterator]() {
    const keys = Object.keys(this)
    let index = 0

    return {
      next: function () {
        if (index < keys.length) {
          const key = keys[index++]
          if(key === Symbol.iterator.toString()) {
            return this.next()
          }
          return {value: person[key], done: false}
        }
        return {done: true}
      }
    }
  }
}

const person2 = {
  name: 'Niko',
  age: 29,
  city: 'Beijing',

  *[Symbol.iterator]() {
    for (const key in person2) {
      if (person2.hasOwnProperty(key) && typeof key !== 'function') {
        yield person2[key]
      }
    }
  }
}

for (const item of person) {
  console.log(item)
}
for (const item of person2) {
  console.log(item)
}