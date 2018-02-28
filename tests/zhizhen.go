
package main

import "fmt"

func change(x *int){
	*x += 1
}

func changeWithoutPointer(y int){
	y += 1
	fmt.Printf("y without pointer is %d\n", y)
}

func main() {
	x := 1
	y := 1
	change(&x)
	changeWithoutPointer(y)
	fmt.Printf("x is %d, y is %d\n", x, y)
}