package config

const defaultConfig = `{
  "ignore_external": true,
  "go_mod_support": false,
  "acss": {
    "breakPoints": {
      "lg": "@media(min-width:1200px)",
      "md": "@media(min-width:1000px)",
      "sm": "@media(min-width:750px)"
    },
    "out": "css/a.scss"
  },
  "compile": {
    "supportStyles": true,
    "stylesOut": "css/main.scss",
    "suffix": "_gas",
    "external_suffix": "_gas_e"
  },
  "watch": {
    "init_tasks": [
      {
        "command": "build",
        "nowait": false,
        "is_gas": true
      },
      {
        "command": "gasx serve",
        "nowait": true,
        "is_gas": false
      }
    ],
    "watchers": [
      {
        "name": "components",
        "recursive": true
      },
      {
        "name": "static",
        "recursive": true
      }
    ],
    "tasks": [
      {
        "command": "build",
        "nowait": false,
        "is_gas": true
      },
      {
        "command": "fuser -n tcp -k 8080; gasx serve",
        "nowait": true,
        "is_gas": false
      }
    ],
    "ignore_compiled": true
  },
  "build": {
    "platform": "gopherjs",
    "sass": "sass INPUT OUTPUT",
    "scss": "scss INPUT OUTPUT",
    "less": "less INPUT OUTPUT",
    "files_deps": [
      {
        "Path": "css/s.scss",
        "Src": "$GOPATH/src/github.com/gascore/gas/std/components/css/spectre.min.css"
      },
      {
        "Path": "css/se.scss",
        "Src": "$GOPATH/src/github.com/gascore/gas/std/components/css/spectre-exp.min.css"
      },
      {
        "Path": "css/si.scss",
        "Src": "$GOPATH/src/github.com/gascore/gas/std/components/css/spectre-icons.min.css"
      }
    ]
  },
  "serve": {
    "port": ":8080"
  },
  "dependencies": {
    "js_out": "dist/deps.js",
    "css_out": "dist/deps.css"
  }
}`
