package main

import "fmt"

type  vertex struct {
	x int
	y int
}
func main(){
	i, j :=42,24990
	p := &i
	fmt.Println(*p)
	*p =21
	fmt.Println(i)
	p =&j
	*p =*p/37
	fmt.Println(j)
 	v :=vertex{1,2}
 	v.x=4
 	fmt.Println(v.x)
 	fmt.Println(v)
	}
