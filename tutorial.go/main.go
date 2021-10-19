package main

import "fmt"

func main() {
	//fmt.Println("Hello world!")
	//fizzbuzz(30)
	p := Person{name:"person", age:16}
	fmt.Println(p.ToString())

	h := Human{name:"human", age:16}
	fmt.Println(h.ToString())
}

type Person struct {
  name string
  age int
}

type Human struct {
  name string
  age int
}

func (p Person, h Human) ToString() string {
  return fmt.Sprintf("name: %s, age: %d", p.name, p.age)
}

func fizzbuzz(x int) {
	var str string
	for i:=0; i<x; i++ {
		str = ""
		if i%3==0 {str += "Fizz"}
		if i%5==0 {str += "Buzz"}
		if str == "" {
			fmt.Printf("%d :\n", i)
		} else {
			fmt.Printf("%d : %s\n", i, str)
		}
	}
}
