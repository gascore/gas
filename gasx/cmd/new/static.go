package new

var indexHtml = `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Gas framework example</title>
    <link rel="stylesheet" href="main.css" />
    <link rel="stylesheet" href="a.css" />
    <link rel="stylesheet" href="s.css" />
    <link rel="stylesheet" href="se.css" />
    <link rel="stylesheet" href="si.css" />
</head>
<body>
    <div id="app"></div>
    <script src="index.js"></script>
</body>
</html>`
var indexGoJsHtml = `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Gas framework example</title>
    <link rel="stylesheet" href="main.scss.css" />
    <link rel="stylesheet" href="a.scss.css" />
    <link rel="stylesheet" href="s.scss.css" />
    <link rel="stylesheet" href="se.scss.css" />
    <link rel="stylesheet" href="si.scss.css" />
</head>
<body>
<div id="app"></div>
<script src="./app.js"></script>
</body>
</html>`

var mainCss = `
<style>
a {
	color: #00CC99;
}
        
a:hover {
	color: #ccc;
}

#lorem-ipsum {
	padding: 8px 16px;
	border: 1px solid #dedede;
	border-radius: 4px;
	margin: 8px;
}

a, b, strong {
	margin: 0 4px;
}

header {
    padding: 0 8px;
}

footer {
    position: absolute;
    right: 0;
    bottom: 0;
    left: 0;
    padding: 6px;
    background-color: #efefef;
    text-align: center;
}
</style>`

var clearSh = `#!/usr/bin/env bash

rm *_gas.go **/*_gas.go components/**/*_gas.go static/app.js static/app.js.map a.css a.min.css css/a.scss css/a.css .gaslock
rm css/*.scss.css css/**/*.scss.css css/*.scss.css.map css/**/*.scss.css.map css/*.sass.css css/*.scss.css css/**/*.sass.css css/**/*.scss.css css/s.css css/se.css css/si.css css/s.scss css/se.scss css/si.scss css/tree.scss css/main.scss
rm -r dist/ .sass-cache/ css/.sass-cache/ css/**/.sass-cache/`

var configJSON = `{
  "ignore_external": true,
  "go_mod_support": false,
  "acss": {
    "breakPoints": {
      "lg": "@media(min-width:1200px)",
      "md": "@media(min-width:1000px)",
      "sm": "@media(min-width:750px)"
    },
    "exceptions": [
      "db(foo)"
    ],
    "custom": {
      "b": "1px solid #dedede"
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
        "name": "store",
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
    "port": ":8080",
    "dir": "dist"
  },
  "dependencies": {
    "js_out": "dist/deps.js",
    "css_out": "dist/deps.css",
    "deps": null
  }
}`

var gitIgnore = `
.gaslock
web_modules
dist
*_gas.go
.sass-cache

*scss.css
*sass.css

css/s.scss
css/s.css
css/si.scss
css/si.css
css/se.scss
css/se.css
css/*.css.map

# ignore app.js files
*app.js*

.idea/
# Created by https://www.gitignore.io/api/go,code,webstorm
# Edit at https://www.gitignore.io/?templates=go,code,webstorm

### Code ###
# Visual Studio Code - https://code.visualstudio.com/
.settings/
.vscode/
tsconfig.json
jsconfig.json

### Go ###
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, build with ` + "`" + `go test -c` + "`" + `
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

### Go Patch ###
/vendor/
/Godeps/

### WebStorm ###
# Covers JetBrains IDEs: IntelliJ, RubyMine, PhpStorm, AppCode, PyCharm, CLion, Android Studio and WebStorm
# Reference: https://intellij-support.jetbrains.com/hc/en-us/articles/206544839

# User-specific stuff
.idea/**/workspace.xml
.idea/**/tasks.xml
.idea/**/usage.statistics.xml
.idea/**/dictionaries
.idea/**/shelf

# Generated files
.idea/**/contentModel.xml

# Sensitive or high-churn files
.idea/**/dataSources/
.idea/**/dataSources.ids
.idea/**/dataSources.local.xml
.idea/**/sqlDataSources.xml
.idea/**/dynamic.xml
.idea/**/uiDesigner.xml
.idea/**/dbnavigator.xml

# Gradle
.idea/**/gradle.xml
.idea/**/libraries

# Gradle and Maven with auto-import
# When using Gradle or Maven with auto-import, you should exclude module files,
# since they will be recreated, and may cause churn.  Uncomment if using
# auto-import.
# .idea/modules.xml
# .idea/*.iml
# .idea/modules

# CMake
cmake-build-*/

# Mongo Explorer plugin
.idea/**/mongoSettings.xml

# File-based project format
*.iws

# IntelliJ
out/

# mpeltonen/sbt-idea plugin
.idea_modules/

# JIRA plugin
atlassian-ide-plugin.xml

# Cursive Clojure plugin
.idea/replstate.xml

# Crashlytics plugin (for Android Studio and IntelliJ)
com_crashlytics_export_strings.xml
crashlytics.properties
crashlytics-build.properties
fabric.properties

# Editor-based Rest Client
.idea/httpRequests

# Android studio 3.1+ serialized cache file
.idea/caches/build_file_checksums.ser

### WebStorm Patch ###
# Comment Reason: https://github.com/joeblau/gitignore.io/issues/186#issuecomment-215987721

# *.iml
# modules.xml
# .idea/misc.xml
# *.ipr

# Sonarlint plugin
.idea/sonarlint

# End of https://www.gitignore.io/api/go,code,webstorm`
