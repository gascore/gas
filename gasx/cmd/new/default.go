package new

var defaultMainGas = `<template>
    <div>
		 <header class="navbar bdb(1px,solid,#dedede) mb(12px)">
            <section class="navbar-section">
                <span class="navbar-brand mr-2">
                    Hello world!
                </span>
            </section>
        </header>

        <main>
            <e run="components.Hello()" />
        </main>

		<footer>
            Created by
            <a href="https://your.site/" target="_blank">
                Your name
            </a>
            with
            <a href="https://github.com/gascore/gas" target="_blank">
                gas
            </a>,
            <a href="https://github.com/gascore/gas/gasx" target="_blank">
                gasx
            </a>
            and love
        </footer>
	</div>
</template>

<script>
package main

import (
	"github.com/gascore/gas"
	"github.com/gascore/gas/web"

	"your_project_path/components"
)

func main() {
    app, err :=
        gas.New(
            web.GetBackEnd(),
            "app",
            &gas.Component{},
            mainT,)
    must(err)

    err = gas.Init(app)
    must(err)
    web.KeepAlive()
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}

var mainT gas.GetComponentChildes
</script>` + mainCss

var defaultHelloGas = `<template>
	<div id="home">
		<h1>
			{{ this.Get("hello") }}
		</h1>
	</div>
</template>

<script>
package components

import (
	"github.com/gascore/gas"
)

func Hello() *gas.Component {
	return gas.NC(
		&gas.C{
			Data: map[string]interface{} {
				"hello": "Hello world!",
			},
		},
		helloT,)
}

var helloT gas.GetComponentChildes
</script>
`
