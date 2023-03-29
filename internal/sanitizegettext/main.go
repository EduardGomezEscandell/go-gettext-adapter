//go:build tools
// +build tools

package main

import (
	"fmt"
	"log"
	"os"

	"example.com/sample/internal/sanitizegettext"
)

const usage = `usage:

	sanitizegettext DST SRC package func
		Copies the directory from SRC to DST replacing all the invalid gettext strings.
		Only go files are copied. You can run xgettext on DST.
`

func main() {
	if len(os.Args) != 5 {
		log.Fatal(usage)
	}

	dst := os.Args[1]
	src := os.Args[2]
	pkg := os.Args[3]
	fun := os.Args[4]

	fmt.Println(dst, src, pkg, fun)

	if err := sanitizegettext.Sanitize(dst, src, pkg, fun); err != nil {
		log.Fatal(err)
	}
}
