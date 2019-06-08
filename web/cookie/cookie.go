package cookie

import (
	"fmt"
	"github.com/gascore/gas/web"
	"syscall/js"
	"strings"
)

// Set set cookie by key and value
func Set(key, value string) {
	js.Global().Set("cookie", fmt.Sprintf("%s=%s", key, value))
}

// Get get cookie by cookie name
func Get(key string) (string, error) {
	cookies := js.Global().Get("cookie").String()
	for _, cookie := range strings.Split(cookies, "; ") {
		cookieSplit := strings.Split(cookie, "=")
		if len(cookieSplit) != 2 {
			return "", web.ErrInvalidCookie
		}

		if cookieSplit[0] == key {
			return cookieSplit[1], nil
		}
	}

	return "", web.ErrCookieNotFound
}
