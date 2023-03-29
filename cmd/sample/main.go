package main

import (
	"fmt"

	"github.com/EduardGomezEscandell/go-gettext-adapter/internal/foo"
	"github.com/EduardGomezEscandell/go-gettext-adapter/internal/i18n"
)

func main() {
	fmt.Println(i18n.G("Welcome!"))
	fmt.Println(foo.Bar())
}
