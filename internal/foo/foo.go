package foo

import "github.com/EduardGomezEscandell/go-gettext-adapter/internal/i18n"

func Bar() string {
	return i18n.G(`Hello, world\n`)
}
