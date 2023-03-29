package foo

import "example.com/sample/internal/i18n"

func Bar() string {
	return i18n.G("Hello, world\\n")
}
