package compile

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func getXAttr(attrs []html.Attribute, x string) string {
	for _, a := range attrs {
		if a.Key == x {
			return a.Val
		}
	}
	return ""
}

func parseAttrs(attrs []html.Attribute) (string, ForInfo, error) {
	// handlers, binds, directives, attrs
	body := make(map[string]string)

	var haveFor bool
	var forForOut string
	var forIsRawData bool
	for _, attr := range attrs {
		switch {
		case attr.Key == gDirs+"ref":
			break
		case strings.HasPrefix(attr.Key, gOn), strings.HasPrefix(attr.Key, gOnAlt):
			var current string
			if strings.HasPrefix(attr.Key, gOn) {
				current = gOn
			} else {
				current = gOnAlt
			}

			var warnErrorS, warnErrorE string
			if attr.Key[len(current)] == '!' {
				current = current + "!"
				warnErrorS = "this.WarnError("
				warnErrorE = ")"
			}

			if body["handlers"] == "" {
				body["handlers"] = "\nHandlers: map[string]gas.Handler{"
			}

			body["handlers"] = body["handlers"] + "\n\"" + strings.TrimPrefix(attr.Key, current) + "\": func(p *gas.Component, e gas.Object) { " + warnErrorS + attr.Val + warnErrorE + " },"

			break
		case strings.HasPrefix(attr.Key, gBind), strings.HasPrefix(attr.Key, gBindAlt):
			var current string
			if strings.HasPrefix(attr.Key, gBind) {
				current = gBind
			} else {
				current = gBindAlt
			}

			if body["binds"] == "" {
				body["binds"] = "\nBinds: map[string]gas.Bind{"
			}

			body["binds"] = fmt.Sprintf(`
%s
"%s": func() string {
	return %s
},
`, body["binds"], strings.TrimPrefix(attr.Key, current), attr.Val)
			break
		case attr.Key == gFor:
			if haveFor {
				continue
			}

			parsedAttr := strings.Split(attr.Val, " in ")
			if len(parsedAttr) < 2 {
				return "", ForInfo{}, errors.New("invalid FOR directive")
			}

			vars := parsedAttr[0]
			data := parsedAttr[1]

			var component string
			if len(parsedAttr) >= 3 {
				component = parsedAttr[2]
			} else {
				component = "this"
			}

			varsParsed := strings.Split(vars, ",")
			var varsI string
			var varsEl string
			if len(varsParsed) == 1 {
				varsI = "i"
				varsEl = trim(varsParsed[0])
			} else {
				varsI = trim(varsParsed[0])
				varsEl = trim(varsParsed[1])
			}

			// If before first data char is '@' => user want to use NewForByData
			if data[0] == '@' {
				data = data[1:]
				forIsRawData = true
				component = ""
			} else {
				component = ", " + component
			}

			haveFor = true
			forForOut = data + component + `, func(` + varsI + ` int, ` + varsEl + ` interface{}) interface{} {`

			break
		case strings.HasPrefix(attr.Key, gDirs): // because it checks last it works with directives only
			directiveName := strings.ToLower(strings.TrimPrefix(attr.Key, gDirs))

			var appendStr string
			switch directiveName {
			case "if":
				appendStr = fmt.Sprintf(`
If: func(p *gas.Component) bool {
	return %s
},
`, attr.Val)
				break
			case "else":
				appendStr = "\nElse: true, \n"
				break
			case "else-if":
				appendStr = fmt.Sprintf(`
If: func(p *gas.Component) bool {
	return %s
},
Else: true,
`, attr.Val)
				break
			case "show":
				appendStr = fmt.Sprintf(`
Show: func(p *gas.Component) bool {
	return %s
},
`, attr.Val)
				break
			case "model":
				var deep, modelData string
				if strings.Index(attr.Val, ".") != -1 || strings.Index(attr.Val, "[") != -1 { // have '.'
					var deepVal string
					dotIndex := strings.Index(attr.Val, ".")
					bracketIndex := strings.Index(attr.Val, "[")
					if bracketIndex == -1 || dotIndex < bracketIndex {
						modelData = attr.Val[:dotIndex]
						deepVal = attr.Val[dotIndex+1:]
					} else {
						modelData = attr.Val[:bracketIndex]
						deepVal = attr.Val[strings.Index(attr.Val[bracketIndex+1:], "]")+1:]
					}

					var deepData []ModelDirectiveDeepData
					deepData = getDeepData(deepData, deepVal)
					for _, deepEl := range deepData {
						deep = deep + fmt.Sprintf(`{Data:%s,Brackets:%t},`, deepEl.Data.(string), deepEl.Brackets)
					}

					deep = "\nDeep: []ModelDirectiveDeepData{" + deep + "},"
				} else {
					modelData = attr.Val
				}

				appendStr = fmt.Sprintf(`
Model: gas.ModelDirective{
	Data: "%s",
	Component: this,%s
},
`, modelData, deep)
			case "html":
				appendStr = fmt.Sprintf(`
HTML: gas.HTMLDirective {
	Render: func(p *gas.Component) string {
		return %s
	},
},
`, attr.Val)
				break
			}

			body["directives"] = fmt.Sprintf("%s\n%s", body["directives"], appendStr)

			break
		default:
			if body["attrs"] == "" {
				body["attrs"] = "\nAttrs: map[string]string{"
			}

			body["attrs"] = body["attrs"] + `"` + attr.Key + `"` + `: "` + attr.Val + `",`

			break
		}
	}

	var out string
	for key, value := range body {
		if key == "directives" {
			out = out + value
			continue
		}

		if len(value) > 1 {
			out = out + value + "\n},"
		}
	}

	ref := getXAttr(attrs, "g-ref")
	if ref != "" {
		out = out + "Ref: \"" + ref + "\","
	}

	return out, ForInfo{haveFor: haveFor, out: forForOut, rawData: forIsRawData}, nil
}

func getDeepData(arr []ModelDirectiveDeepData, s string) []ModelDirectiveDeepData {
	dotIndex := strings.Index(s, ".")
	bracketIndex := strings.Index(s, "[")

	if dotIndex == -1 && bracketIndex == -1 {
		return append(arr, ModelDirectiveDeepData{Data: `"` + s + `"`, Brackets: false})
	}

	if dotIndex > bracketIndex {
		a := s[:dotIndex]
		arr = append(arr, ModelDirectiveDeepData{Data: `"` + a + `"`})
		s = s[dotIndex+1:]
	} else {
		if bracketIndex != 0 { // '[' not first
			arr = append(arr, ModelDirectiveDeepData{Data: `"` + s[:bracketIndex] + `"`, Brackets: false})
			s = s[bracketIndex:]
		}

		endBracketIndex := strings.Index(s, "]")
		if endBracketIndex == -1 {
			panic("invalid bracket indexes")
		}
		a := s[1:endBracketIndex]

		arr = append(arr, ModelDirectiveDeepData{Data: a, Brackets: true})

		s = s[endBracketIndex+1:]
	}

	if strings.Index(s, ".") != -1 || strings.Index(s, "[") != -1 {
		return getDeepData(arr, s)
	}

	return arr
}

type ModelDirectiveDeepData struct{
	Data     interface{}
	Brackets bool
}
