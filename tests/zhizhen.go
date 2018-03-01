
package main

import (
	"fmt"
	"time"
)

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
	addtimes :=time.Now().Unix()
	fmt.Println(addtimes)
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, 1)
	fmt.Println(tm2)
}