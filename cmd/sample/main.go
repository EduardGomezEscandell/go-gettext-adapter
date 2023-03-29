package main

import (
	"fmt"

	"example.com/sample/internal/foo"
	"example.com/sample/internal/i18n"
)

func main() {
	fmt.Println(i18n.G("Welcome!"))
	fmt.Println(foo.Bar())
}
