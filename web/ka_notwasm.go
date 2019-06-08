// +build !wasm

package web

import "fmt"

// KeepAlive keep alive application
func KeepAlive() {
	ch := make(chan int, 5)
	ch <- 1
	fmt.Println(<-ch)
}
