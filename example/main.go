package main

import (
	"fmt"
	"github.com/sinicablyat/gas"
	"syscall/js"
)

// This all seems very weired, BUT! this will auto generate from nice-look .gas components.
// HelloWorld component will look like this:
//
//	`
//		<h1 id=hello-world>
//			Hello, from Gas!
//		</h1
//	`
func main() {
	app, err :=
		gas.New(
			"app",
			gas.NewComponent(gas.NilData, gas.NilData).
				AddInfo("h1", "#hello-world", gas.NilClasses, gas.NilAttrs).
				AddChildes(gas.SendComponents([]interface{}{
					"Hello, from Gas!",
				})))
	if err != nil {
		panic(err)
	}

	js.Global().Get("document").Call("getElementsByTagName", "body").Index(0).Set("innerHTML",
		fmt.Sprintf("%s",
			app.App.
				Childes(app.App)[0].(*gas.Component).
					Childes(*app.App.Childes(app.App)[0].(*gas.Component))[0].(string)))
}
