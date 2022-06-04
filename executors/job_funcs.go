package executors

import "fmt"

func DoAdd(x, y float64) int {
	fmt.Printf("DoAdd: %d + %d\n", int(x), int(y))
	return int(x) + int(y)
}

func DoPrint(v ...interface{}) {
	fmt.Println(v...)
}
