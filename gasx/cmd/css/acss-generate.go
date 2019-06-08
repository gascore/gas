package css

import (
	"fmt"
	"strings"
)

func GenerateStyleForClass(class string, custom map[string]string) string {
	cDots := strings.Contains(class, ":")
	cC := strings.Contains(class, "(")
	if !(cC || cDots) {
		return ""
	}

	bi := strings.Index(class, "(")
	value := class[bi+1 : strings.Index(class, ")")]
	class = class[:bi]

	if len(custom[value]) != 0 {
		value = custom[value]
	} else {
		value = strings.Replace(value, ",", " ", -1)
	}

	switch class {

	case "pb":
		return fmt.Sprintf(`padding-bottom: %s;`, value)

	case "trf":
		return fmt.Sprintf(`transform: %s;`, value)

	case "animn":
		return fmt.Sprintf(`animation-name: %s;`, value)

	case "bdr+":
		return fmt.Sprintf(`border-right: %s;`, value)

	case "bdc":
		return fmt.Sprintf(`border-color: %s;`, value)

	case "bdi":
		return fmt.Sprintf(`border-image:url(%s);`, value)

	case "fz":
		return fmt.Sprintf(`font-size: %s;`, value)

	case "fef":
		return fmt.Sprintf(`font-effect: %s;`, value)

	case "f":
		return fmt.Sprintf(`font: %s;`, value)

	case "bdw":
		return fmt.Sprintf(`border-width: %s;`, value)

	case "wos":
		return fmt.Sprintf(`word-spacing: %s;`, value)

	case "trsde":
		return fmt.Sprintf(`transition-delay: %s;`, value)

	case "bxz":
		return fmt.Sprintf(`box-sizing: %s;`, value)

	case "bga":
		return fmt.Sprintf(`background-attachment: %s;`, value)

	case "fza":
		return fmt.Sprintf(`font-size-adjust: %s;`, value)

	case "fem":
		return fmt.Sprintf(`font-emphasize: %s;`, value)

	case "animdel":
		return fmt.Sprintf(`animation-delay: %s;`, value)

	case "colmr":
		return fmt.Sprintf(`column-rule: %s;`, value)

	case "to":
		return fmt.Sprintf(`text-outline: %s;`, value)

	case "bdbk":
		return fmt.Sprintf(`border-break: %s;`, value)

	case "bgpy":
		return fmt.Sprintf(`background-position-y: %s;`, value)

	case "colmg":
		return fmt.Sprintf(`column-gap: %s;`, value)

	case "bdrst":
		return fmt.Sprintf(`border-right-style: %s;`, value)

	case "coi":
		return fmt.Sprintf(`counter-increment: %s;`, value)

	case "colmrs":
		return fmt.Sprintf(`column-rule-style: %s;`, value)

	case "h":
		return fmt.Sprintf(`height: %s;`, value)

	case "ols":
		return fmt.Sprintf(`outline-style: %s;`, value)

	case "femp":
		return fmt.Sprintf(`font-emphasize-position: %s;`, value)

	case "ord":
		return fmt.Sprintf(`order: %s;`, value)

	case "pos":
		return fmt.Sprintf(`position: %s;`, value)

	case "cl":
		return fmt.Sprintf(`clear: %s;`, value)

	case "wm":
		return fmt.Sprintf(`writing-mode: %s;`, value)

	case "pgba":
		return fmt.Sprintf(`page-break-after: %s;`, value)

	case "anim":
		return fmt.Sprintf(`animation: %s;`, value)

	case "animdir":
		return fmt.Sprintf(`animation-direction: %s;`, value)

	case "ol":
		return fmt.Sprintf(`outline: %s;`, value)

	case "tsh":
		return fmt.Sprintf(`text-shadow: %s;`, value)

	case "olc":
		return fmt.Sprintf(`outline-color: %s;`, value)

	case "fxsh":
		return fmt.Sprintf(`flex-shrink: %s;`, value)

	case "@media":
		return fmt.Sprintf(`@media %s 
`, value)

	case "fl":
		return fmt.Sprintf(`float: %s;`, value)

	case "bdts":
		return fmt.Sprintf(`border-top-style: %s;`, value)

	case "bdtw":
		return fmt.Sprintf(`border-top-width: %s;`, value)

	case "orp":
		return fmt.Sprintf(`orphans: %s;`, value)

	case "fx":
		return fmt.Sprintf(`flex: %s;`, value)

	case "bdbrrs":
		return fmt.Sprintf(`border-bottom-right-radius: %s;`, value)

	case "tbl":
		return fmt.Sprintf(`table-layout: %s;`, value)

	case "animps":
		return fmt.Sprintf(`animation-play-state: %s;`, value)

	case "bgp":
		return fmt.Sprintf(`background-position: %s;`, value)

	case "@import":
		return fmt.Sprintf(`@import url(%s);`, value)

	case "bgo":
		return fmt.Sprintf(`background-origin: %s;`, value)

	case "op":
		return fmt.Sprintf(`opacity: %s;`, value)

	case "p":
		return fmt.Sprintf(`padding: %s;`, value)

	case "bgc":
		return fmt.Sprintf(`background-color: %s;`, value)

	case "ac":
		return fmt.Sprintf(`align-content: %s;`, value)

	case "@i":
		return fmt.Sprintf(`@import url(%s);`, value)

	case "pr":
		return fmt.Sprintf(`padding-right: %s;`, value)

	case "pgbi":
		return fmt.Sprintf(`page-break-inside: %s;`, value)

	case "whs":
		return fmt.Sprintf(`white-space: %s;`, value)

	case "bdb":
		return fmt.Sprintf(`border-bottom: %s;`, value)

	case "cnt":
		return fmt.Sprintf(`content:'%s';`, value)

	case "ovy":
		return fmt.Sprintf(`overflow-y: %s;`, value)

	case "c":
		return fmt.Sprintf(`color: %s;`, value)

	case "colmw":
		return fmt.Sprintf(`column-width: %s;`, value)

	case "bdrs":
		return fmt.Sprintf(`border-radius: %s;`, value)

	case "rsz":
		return fmt.Sprintf(`resize: %s;`, value)

	case "@m":
		return fmt.Sprintf(`@media %s 
`, value)

	case "bdti":
		return fmt.Sprintf(`border-top-image:url(%s);`, value)

	case "tov":
		return fmt.Sprintf(`text-overflow: %s;`, value)

	case "mah":
		return fmt.Sprintf(`max-height: %s;`, value)

	case "bdcl":
		return fmt.Sprintf(`border-collapse: %s;`, value)

	case "cur":
		return fmt.Sprintf(`cursor: %s;`, value)

	case "ori":
		return fmt.Sprintf(`orientation: %s;`, value)

	case "tw":
		return fmt.Sprintf(`text-wrap: %s;`, value)

	case "fxd":
		return fmt.Sprintf(`flex-direction: %s;`, value)

	case "ec":
		return fmt.Sprintf(`empty-cells: %s;`, value)

	case "cor":
		return fmt.Sprintf(`counter-reset: %s;`, value)

	case "bdrw":
		return fmt.Sprintf(`border-right-width: %s;`, value)

	case "miw":
		return fmt.Sprintf(`min-width: %s;`, value)

	case "list":
		return fmt.Sprintf(`list-style-type: %s;`, value)

	case "bdbli":
		return fmt.Sprintf(`border-bottom-left-image:url(%s);`, value)

	case "pt":
		return fmt.Sprintf(`padding-top: %s;`, value)

	case "bdtlrs":
		return fmt.Sprintf(`border-top-left-radius: %s;`, value)

	case "colmc":
		return fmt.Sprintf(`column-count: %s;`, value)

	case "pl":
		return fmt.Sprintf(`padding-left: %s;`, value)

	case "bd":
		return fmt.Sprintf(`border: %s;`, value)

	case "fs":
		return fmt.Sprintf(`font-style: %s;`, value)

	case "v":
		return fmt.Sprintf(`visibility: %s;`, value)

	case "w":
		return fmt.Sprintf(`width: %s;`, value)

	case "z":
		return fmt.Sprintf(`z-index: %s;`, value)

	case "bgpx":
		return fmt.Sprintf(`background-position-x: %s;`, value)

	case "trsdu":
		return fmt.Sprintf(`transition-duration: %s;`, value)

	case "colmrc":
		return fmt.Sprintf(`column-rule-color: %s;`, value)

	case "bd+":
		return fmt.Sprintf(`border: %s;`, value)

	case "bgi":
		return fmt.Sprintf(`background-image:url(%s);`, value)

	case "lisi":
		return fmt.Sprintf(`list-style-image: %s;`, value)

	case "fsm":
		return fmt.Sprintf(`font-smooth: %s;`, value)

	case "colmrw":
		return fmt.Sprintf(`column-rule-width: %s;`, value)

	case "bdbi":
		return fmt.Sprintf(`border-bottom-image:url(%s);`, value)

	case "ti":
		return fmt.Sprintf(`text-indent: %s;`, value)

	case "tr":
		return fmt.Sprintf(`text-replace: %s;`, value)

	case "as":
		return fmt.Sprintf(`align-self: %s;`, value)

	case "bdt":
		return fmt.Sprintf(`border-top: %s;`, value)

	case "bg":
		return fmt.Sprintf(`background: %s;`, value)

	case "m":
		return fmt.Sprintf(`margin: %s;`, value)

	case "bdrc":
		return fmt.Sprintf(`border-right-color: %s;`, value)

	case "bdri":
		return fmt.Sprintf(`border-right-image:url(%s);`, value)

	case "bdls":
		return fmt.Sprintf(`border-left-style: %s;`, value)

	case "ct":
		return fmt.Sprintf(`content: %s;`, value)

	case "ml":
		return fmt.Sprintf(`margin-left: %s;`, value)

	case "fxb":
		return fmt.Sprintf(`flex-basis: %s;`, value)

	case "animfm":
		return fmt.Sprintf(`animation-fill-mode: %s;`, value)

	case "trfs":
		return fmt.Sprintf(`transform-style: %s;`, value)

	case "mih":
		return fmt.Sprintf(`min-height: %s;`, value)

	case "wob":
		return fmt.Sprintf(`word-break: %s;`, value)

	case "ovx":
		return fmt.Sprintf(`overflow-x: %s;`, value)

	case "bdbw":
		return fmt.Sprintf(`border-bottom-width: %s;`, value)

	case "bdbs":
		return fmt.Sprintf(`border-bottom-style: %s;`, value)

	case "mb":
		return fmt.Sprintf(`margin-bottom: %s;`, value)

	case "fv":
		return fmt.Sprintf(`font-variant: %s;`, value)

	case "d":
		return fmt.Sprintf(`display: %s;`, value)

	case "br":
		return fmt.Sprintf(`border-right: %s;`, value)

	case "colms":
		return fmt.Sprintf(`column-span: %s;`, value)

	case "bgbk":
		return fmt.Sprintf(`background-break: %s;`, value)

	case "whsc":
		return fmt.Sprintf(`white-space-collapse: %s;`, value)

	case "animic":
		return fmt.Sprintf(`animation-iteration-count: %s;`, value)

	case "animtf":
		return fmt.Sprintf(`animation-timing-function: %s;`, value)

	case "jc":
		return fmt.Sprintf(`justify-content: %s;`, value)

	case "bdlw":
		return fmt.Sprintf(`border-left-width: %s;`, value)

	case "td":
		return fmt.Sprintf(`text-decoration: %s;`, value)

	case "wow":
		return fmt.Sprintf(`word-wrap: %s;`, value)

	case "colm":
		return fmt.Sprintf(`columns: %s;`, value)

	case "tt":
		return fmt.Sprintf(`text-transform: %s;`, value)

	case "bdtri":
		return fmt.Sprintf(`border-top-right-image:url(%s);`, value)

	case "trfo":
		return fmt.Sprintf(`transform-origin: %s;`, value)

	case "fw":
		return fmt.Sprintf(`font-weight: %s;`, value)

	case "va":
		return fmt.Sprintf(`vertical-align: %s;`, value)

	case "trs":
		return fmt.Sprintf(`transition: %s;`, value)

	case "bdb+":
		return fmt.Sprintf(`border-bottom: %s;`, value)

	case "lh":
		return fmt.Sprintf(`line-height: %s;`, value)

	case "fxg":
		return fmt.Sprintf(`flex-grow: %s;`, value)

	case "bgr":
		return fmt.Sprintf(`background-repeat: %s;`, value)

	case "te":
		return fmt.Sprintf(`text-emphasis: %s;`, value)

	case "q":
		return fmt.Sprintf(`quotes: %s;`, value)

	case "pgbb":
		return fmt.Sprintf(`page-break-before: %s;`, value)

	case "bdbri":
		return fmt.Sprintf(`border-bottom-right-image:url(%s);`, value)

	case "ta":
		return fmt.Sprintf(`text-align: %s;`, value)

	case "bt":
		return fmt.Sprintf(`border-top: %s;`, value)

	case "bds":
		return fmt.Sprintf(`border-style: %s;`, value)

	case "us":
		return fmt.Sprintf(`user-select: %s;`, value)

	case "anim-":
		return fmt.Sprintf(`animation: %s;`, value)

	case "maw":
		return fmt.Sprintf(`max-width: %s;`, value)

	case "lisp":
		return fmt.Sprintf(`list-style-position: %s;`, value)

	case "fxf":
		return fmt.Sprintf(`flex-flow: %s;`, value)

	case "bdtli":
		return fmt.Sprintf(`border-top-left-image:url(%s);`, value)

	case "cps":
		return fmt.Sprintf(`caption-side: %s;`, value)

	case "bgsz":
		return fmt.Sprintf(`background-size: %s;`, value)

	case "r":
		return fmt.Sprintf(`right: %s;`, value)

	case "bdsp":
		return fmt.Sprintf(`border-spacing: %s;`, value)

	case "bb":
		return fmt.Sprintf(`border-bottom: %s;`, value)

	case "wfsm":
		return fmt.Sprintf(`-webkit-font-smoothing: %s;`, value)

	case "colmf":
		return fmt.Sprintf(`column-fill: %s;`, value)

	case "bdt+":
		return fmt.Sprintf(`border-top: %s;`, value)

	case "ta-lst":
		return fmt.Sprintf(`text-align-last: %s;`, value)

	case "th":
		return fmt.Sprintf(`text-height: %s;`, value)

	case "l":
		return fmt.Sprintf(`left: %s;`, value)

	case "mir":
		return fmt.Sprintf(`min-resolution: %s;`, value)

	case "olw":
		return fmt.Sprintf(`outline-width: %s;`, value)

	case "animdur":
		return fmt.Sprintf(`animation-duration: %s;`, value)

	case "ap":
		return fmt.Sprintf(`appearance: %s;`, value)

	case "f+":
		return fmt.Sprintf(`font: %s;`, value)

	case "mar":
		return fmt.Sprintf(`max-resolution: %s;`, value)

	case "bdtrrs":
		return fmt.Sprintf(`border-top-right-radius: %s;`, value)

	case "lts":
		return fmt.Sprintf(`letter-spacing: %s;`, value)

	case "b":
		return fmt.Sprintf(`bottom: %s;`, value)

	case "fxw":
		return fmt.Sprintf(`flex-wrap: %s;`, value)

	case "bdr":
		return fmt.Sprintf(`border-right: %s;`, value)

	case "ovs":
		return fmt.Sprintf(`overflow-style: %s;`, value)

	case "cp":
		return fmt.Sprintf(`clip: %s;`, value)

	case "ov":
		return fmt.Sprintf(`overflow: %s;`, value)

	case "bdci":
		return fmt.Sprintf(`border-corner-image:url(%s);`, value)

	case "ff":
		return fmt.Sprintf(`font-family: %s;`, value)

	case "bxsh":
		return fmt.Sprintf(`box-shadow: %s;`, value)

	case "bdlen":
		return fmt.Sprintf(`border-length: %s;`, value)

	case "mt":
		return fmt.Sprintf(`margin-top: %s;`, value)

	case "bgcp":
		return fmt.Sprintf(`background-clip: %s;`, value)

	case "mr":
		return fmt.Sprintf(`margin-right: %s;`, value)

	case "bfv":
		return fmt.Sprintf(`backface-visibility: %s;`, value)

	case "bdblrs":
		return fmt.Sprintf(`border-bottom-left-radius: %s;`, value)

	case "to+":
		return fmt.Sprintf(`text-outline: %s;`, value)

	case "ai":
		return fmt.Sprintf(`align-items: %s;`, value)

	case "bdli":
		return fmt.Sprintf(`border-left-image:url(%s);`, value)

	case "bdf":
		return fmt.Sprintf(`border-fit: %s;`, value)

	case "wid":
		return fmt.Sprintf(`widows: %s;`, value)

	case "lis":
		return fmt.Sprintf(`list-style: %s;`, value)

	case "fems":
		return fmt.Sprintf(`font-emphasize-style: %s;`, value)

	case "bdtc":
		return fmt.Sprintf(`border-top-color: %s;`, value)

	case "bdbc":
		return fmt.Sprintf(`border-bottom-color: %s;`, value)

	case "trstf":
		return fmt.Sprintf(`transition-timing-function: %s;`, value)

	case "olo":
		return fmt.Sprintf(`outline-offset: %s;`, value)

	case "bdlc":
		return fmt.Sprintf(`border-left-color: %s;`, value)

	case "bdl":
		return fmt.Sprintf(`border-left: %s;`, value)

	case "tj":
		return fmt.Sprintf(`text-justify: %s;`, value)

	case "trsp":
		return fmt.Sprintf(`transition-property: %s;`, value)

	case "fst":
		return fmt.Sprintf(`font-stretch: %s;`, value)

	case "t":
		return fmt.Sprintf(`top: %s;`, value)
	}
	return ""
}
