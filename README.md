# Go gettext adaptor

The interesting part of this repository is `internal/sanitizegettext`. This module takes a directory
and recursively copies all the `.go` (but not the `_test.go`) files into a destination directory.

In the destination directory, the calls to a specific function `i18n.G` are analyzed. All the instances where
it is succeeded by a string that is not surrounded by `"double quotes"`, the strings are replaced (and sanitized)
with their double-quoted version.

File `internal/foo/foo.go` has such an instance:
```go
return i18n.G(`Hello, world\n`)
```

This would silently fail with gettext as it does not understand ```tilde quotes```, so it is ignored. With the help
of this module, the string does show up in the `.pot` files, see `po/sample.pot`.

The end result is that `go generate internal/i18n/i18n.go` performs this fix, and then calls gettext.