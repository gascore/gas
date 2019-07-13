// +build wasm

package web

var signal = make(chan struct{})

func KeepAlive() {
	for {
		<-signal
	}
}
