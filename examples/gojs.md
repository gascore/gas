# How to run examples with gopherjs backend

You need:

1. Install [gopher-js](https://github.com/gopherjs/gopherjs) `go get -u github.com/gopherjs/gopherjs`
2. Install [node.js](https://nodejs.org/en/) (you can use [nvm](https://github.com/creationix/nvm))
3. Install source-map-support package `npm install --global source-map-support`
4. cd to example
5. run `gopherjs build -m -o app.js`
6. Fix index.gojs.html: change `<script src="TODO/app.js"></script>` to `<script src="{your-example-name}/app.js"></script>`
7. Open index.gojs.html in your browser (just throw file to new tab)