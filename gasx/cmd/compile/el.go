package compile

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func executeEl(t *html.Node, isMain bool) (string, error) {
	var out string

	if t.Type == html.CommentNode {
		return "/*" + t.Data + "*/", nil
	}

	if t.Type != html.ElementNode {
		return returnOutText(t.Data), nil
	}

	if t.Data == "e" {
		runA := getXAttr(t.Attr, "run")
		if len(runA) == 0 {
			return "", fmt.Errorf("invalid external component: no 'run' attribute")
		}

		out = runA

		var body, slots, templates string

		for c := t.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				if c.Data == "template" {
					tName := getXAttr(c.Attr, "name")
					if len(tName) == 0 {
						log.Println("ERROR: invalid template syntax: 'name' attribute undefined")
						continue
					}

					var templateBody string
					for child := c.FirstChild; child != nil; child = child.NextSibling {
						childOut, err := executeEl(child, false)
						if err != nil {
							return "", err
						}

						templateBody = templateBody + childOut
					}

					tTypes := getXAttr(c.Attr, "types")
					if len(tName) != 0 {
						tTypes = "\n" + tTypes + "\n"
					}

					templates = templates + `"` + tName + `": func(values ...interface{}) []interface{} {` + tTypes + `return []interface{} {` + templateBody + `} },`
					continue
				}

				slotAttr := getXAttr(c.Attr, "slot")
				if len(slotAttr) != 0 {
					cOut, err := executeEl(c, false)
					if err != nil {
						return "", err
					}

					slots = slots + `"` + slotAttr + `":` + cOut
					continue
				}
			}

			cOut, err := executeEl(c, false)
			if err != nil {
				return "", err
			}

			body = body + cOut
		}

		var external string

		if len(body) != 0 {
			external = `Body: []interface{}{` + body + `},`
		}

		if len(slots) != 0 {
			external = external + "Slots: map[string]interface{}{" + slots + "},"
		}

		if len(templates) != 0 {
			external = external + "Templates: map[string]gas.Template{" + templates + "},"
		}

		if len(external) > 0 {
			external = "gas.External{" + external
			if out[len(out)-2] != '(' {
				external = ", " + external
			}

			out = out[:len(out)-1] + external + "})"
		}

		return out + ",", nil
	}

	if t.Data == "g-slot" {
		name := getXAttr(t.Attr, "name")
		if len(name) == 0 {
			panic("slot name undefined")
		}
		return `e.Slots["` + name + `"],`, nil
	}

	if t.Data == "g-body" {
		return `gas.NE(&gas.C{}, e.Body),`, nil
	}

	attrs, infoAboutFor, err := parseAttrs(t.Attr)
	if err != nil {
		return "", err
	}

	tBody := "\n&gas.Component{Tag:\"" + t.Data + "\"," + attrs + "},"

	for c := t.FirstChild; c != nil; c = c.NextSibling {
		cOut, err := executeEl(c, false)
		if err != nil {
			return "", err
		}

		out = out + cOut

		t.FirstChild = c.NextSibling
	}

	toOut := tBody + out
	var finallyOut string
	switch {
	case isMain:
		finallyOut = returnOutMainComponent(toOut)
	case infoAboutFor.haveFor:
		if infoAboutFor.rawData {
			finallyOut = returnOutFORRawData(infoAboutFor.out, toOut)
		} else {
			finallyOut = returnOutFOR(infoAboutFor.out, toOut)
		}
	default:
		finallyOut = returnOutComponent(toOut)
	}

	return finallyOut, nil
}

type ForInfo struct {
	haveFor bool
	out     string
	rawData bool
}

func parseText(in string, i int, buf string) string {
	if i > 0 && in[i-1] == '\\' {
		return buf + returnOutStatic(removeInArr([]byte(in), i-1))
	}

	closeI := strings.Index(in, "}}")
	if closeI == -1 {
		return buf + returnOutStatic(in)
	}

	s1 := in[:i]        // was edit in past
	s3 := in[closeI+2:] // will be edited

	s2 := in[i : closeI+2]     // editing now
	s2 = s2[2:len(s2)-2] + "," // remove first and last 2 chars "{{" and "}}"

	buf = buf + "`" + trim(s1) + "`, " + s2

	startI := strings.Index(s3, "{{")

	switch {
	case startI != -1:
		return parseText(s3, startI, buf)
	case startI == -1 && stringNotEmpty(s3):
		return buf + "`" + trim(s3) + "`," // hard trim for last element
	default:
		return buf
	}
}
