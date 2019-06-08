package compile

import (
	"strings"
)

func returnOutStatic(out string) string {
	return "`" + out + "`,"
}

func returnOutText(out string) string {
	if len(trim(out)) == 0 {
		return ""
	}

	i := strings.Index(out, "{{")

	switch {
	case i != -1:
		return parseText(out, i, "")
	case len(out) != 0:
		return returnOutStatic(out)
	default:
		return ""
	}
}

func returnOutComponent(out string) string {
	return "\ngas.NE(" + out + "),"
}

func returnOutFOR(header, out string) string {
	return "gas.NewFor(" + header + "\n return gas.NE(" + out + ")})"
}

func returnOutFORRawData(header, out string) string {
	return "gas.NewForByData(" + header + "\n return gas.NE(" + out + ")}),"
}

func returnOutMainComponent(out string) string {
	return "func(this *gas.Component) []interface{} {return gas.ToGetComponentList(gas.NE(" + out + "),)}"
}
