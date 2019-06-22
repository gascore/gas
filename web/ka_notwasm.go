// +build !wasm

package web

// KeepAlive keep alive application
func KeepAlive() {
	ch := make(chan int, 5)
	ch <- 1
}
