// case1 原型链继承
// 优点：子类可以复用父类原型上的变量、函数
// 缺点：多个子类共享同一父类，如果父类的值变化，影响所有子类
function case1 () {
  function Animal(sound) {
    this.sound = sound
  }

  Animal.prototype.getSound = function () {
    console.log(this.sound)
  }

  function Dog(name) {
    this.name = name
  }

  Dog.prototype = new Animal('Woof')

  const myDog = new Dog('Buddy')
  console.log(myDog.getSound())
}

// case2 构造函数继承
// 优点：可以从子类初始化父类，给父类传递参数；子类修改父类引用属性不会引起全局变化
// 缺点：父子类不在同一条原型链内，子类无法共享父类的属性、方法，而是在初始化时拷贝一份父类数据在子类内部
function case2 () {
  function Animal (sound) {
    this.sound = sound
    this.color = ['red', 'blue']
  }

  Animal.prototype.getSound = function() {
    console.log(this.sound)
  }

  function Dog (name) {
    this.name = name
    Animal.call(this, 'woof')
  }

  const myDog = new Dog('Buddy')
  myDog.color.push('green')
  console.log(myDog.name)
  console.log(myDog.sound)
  console.log(myDog.color)
  console.log(myDog.getSound()) // error: getSound is not a function 
}

// case3: 组合继承（即原型链继承+构造函数继承）
// 优点：综合了两种继承方式的优点，既可以初始化时在子类构造函数给父类构造函数传参，也可以子类调用父类函数
// 缺点：构造函数被创建了两次(子类初始化时调用一次，绑定父类原型链时执行一次)
function case3() {
  function Animal (name) {
    this.name = name
    this.color = ['red', 'blue']
  }
  Animal.prototype.getName = function() {
    console.log(this.name)
  }
  Animal.prototype.getColor = function () {
    console.log(this.color)
  }

  function Dog (name, age) {
    // 构造函数继承
    Animal.call(this, name)
    this.age = age
  }

  // 原型链继承
  Dog.prototype = new Animal()
  Dog.prototype.constructor = Dog
  
  const myDog = new Dog('Buddy', 8)
  myDog.getName()
  myDog.getColor()
  myDog.color.push('green')
  myDog.getColor() 
}

// case4: 寄生组合继承
// js实现OOP继承的最优解的，也是ES6后类语法糖的底层
function case4() {
  function InheritPrototype (parent, children) {
    const proto = Object.create(parent.prototype)
    children.prototype = proto
    proto.constructor = children
  }


  function Animal(sound) {
    this.sound = sound
    this.color = ['red', 'blue']
  }
  Animal.prototype.getSound = function () {
    console.log(this.sound)
  }
  Animal.prototype.getColor = function () {
    console.log(this.color)
  }

  function Dog (name, sound) {
    this.name = name
    Animal.call(this, sound)
  }

  InheritPrototype(Animal, Dog)
  Dog.prototype.Woof = function () {
    console.log('Wooffffff~~~')
  }

  const myDog = new Dog('Buddy', 'gaaaaaa!')
  myDog.Woof()
  myDog.getSound()
  myDog.getColor()
  myDog.color.push('green')
  myDog.getColor()
}

// case5: es2015后的现代语法
function case5 () {
  class Animal {
    constructor(sound){
      this.sound = sound
      this.color = ['red', 'blue']
    }

    getSound() {
      console.log(this.sound)
    }

    getColor() {
      console.log(this.color)
    }

    addColor(color) {
      this.color.push(color)
    }

    removeColor() {
      this.color.pop()
    }
  }

  class Dog extends Animal {
    constructor(name, sound) {
      super(sound)
      this.name = name
    }

    getName() {
      console.log(this.name)
    }
  }

  const myDog = new Dog('Buddy', 'Woof~')
  myDog.getName()
  myDog.getSound()
  myDog.getColor()
  myDog.addColor('green')
  myDog.getColor()
}

// case1()
// case2()
// case3()
// case4()
case5()