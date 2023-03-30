# Go gettext adaptor

The interesting part of this repository is `internal/sanitizegettext`. This module takes a directory
and recursively copies all the `.go` (but not the `_test.go`) files into a destination directory.

In the destination directory, the calls to a specific function `i18n.G` are analyzed (the package and function names are customizable).

## Bad quotation marks
All the instances where
it is succeeded by a string that is not surrounded by `"double quotes"`, the strings are replaced (and sanitized)
with their double-quoted version.

File `internal/foo/foo.go` has such an instance:
```go
return i18n.G(`Hello, world\n`)
```
These cases would silently fail with gettext as it does not understand ```tilde quotes```, so it is ignored. With the help
of this module, the string does show up in the `.pot` files, see `po/sample.pot`.

## No string literals
If any call to `i18n.G` is not followed by a string (e.g `i18n.G(variable)`) the program will report it and exit. This would also silently fail with gettext, and the entry would be ignored.

## Usage

The end result is that `go generate internal/i18n/i18n.go` performs this fix, and then calls gettext.

You can also use it to autofix other projects:
```bash
# Find all issues
go run ./internal/sanitizegettext/main.go /tmp/out/ /path/to/repo i18n G

# Replace with fixes
cp -Tr /tmp/out/ /path/to/repo
```
