package sanitizegettext

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"example.com/sample/internal/finitestatemachine"
)

func Sanitize(dstRoot string, srcRoot string, i18nPkg string, gettextFunc string) error {
	dstTemp, err := SanitizeTemp(srcRoot, i18nPkg, gettextFunc)
	if err != nil {
		return err
	}

	return os.Rename(dstTemp, dstRoot)
}

func SanitizeTemp(srcRoot string, i18nPkg string, gettextFunc string) (string, error) {
	dstTmp, err := os.MkdirTemp(os.TempDir(), "sanitation*")
	if err != nil {
		return "", fmt.Errorf("could not create temp destination dir: %v", err)
	}

	err = mapGoFiles(dstTmp, srcRoot, func(path string, r io.Reader, w io.Writer) error {
		contents, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read file %q: %v", path, err)
		}

		fs := token.NewFileSet()
		f := fs.AddFile(path, fs.Base(), len(contents))

		var s scanner.Scanner
		s.Init(f, contents, nil, scanner.ScanComments)

		machine := finitestatemachine.New(i18nPkg, gettextFunc)
		for {
			pos, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}

			if err := machine.Consume(pos, tok, lit); err != nil {
				line, char := findPos(contents, int(pos))
				log.Fatalf("%s:%d:%d: %v", path, line, char, err)
			}
		}

		var i int
		for _, res := range machine.Results {
			// Skipping properly formatted files
			if res.Val[0] == '"' {
				continue
			}

			replaceWith := res.Val[1 : len(res.Val)-1]

			// Logging
			ln, ch := findPos(contents, res.Pos)
			fmt.Printf("%s:%d:%d: %s => %q\n", path, ln, ch, res.Val, replaceWith)

			// Mapping
			if _, err := w.Write(contents[i : res.Pos-1]); err != nil {
				return fmt.Errorf("could not write on tmp copy of %q: %v", path, err)
			}

			if _, err := fmt.Fprintf(w, `%q`, replaceWith); err != nil {
				return fmt.Errorf("could not write on tmp copy of %q: %v", path, err)
			}

			i = res.Pos + len(res.Val) - 1
		}

		if _, err := w.Write(contents[i:]); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		os.RemoveAll(dstTmp)
	}

	return dstTmp, err
}

func mapGoFiles(dstRoot string, srcRoot string, f func(name string, r io.Reader, w io.Writer) error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	errch := make(chan error)

	if srcRoot, err := filepath.Abs(srcRoot); err != nil {
		return fmt.Errorf("could not get absoulte path for %q: %v", srcRoot, err)
	}

	err := filepath.WalkDir(srcRoot, func(src string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcRoot, src)
		if err != nil {
			return err
		}
		dest := filepath.Join(dstRoot, rel)

		if d.IsDir() {
			if err := os.MkdirAll(dest, 0750); err != nil {
				return fmt.Errorf("could not mkdir %q: %v", dest, err)
			}
			return nil
		}

		// Ignore non-go files
		if filepath.Ext(src) != ".go" {
			return nil
		}

		// Ignore test files
		if strings.HasSuffix(src, "_test.go") {
			return nil
		}

		// Call custom function asyncronously
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Quick escape
			select {
			case <-ctx.Done():
				errch <- ctx.Err()
			default:
			}

			// Run custom func
			err := func() error {
				// Open read-only file
				r, err := os.Open(src)
				if err != nil {
					return fmt.Errorf("could not read %q: %v", src, err)
				}
				defer r.Close()

				// Open destination file
				w, err := os.Create(dest)
				if err != nil {
					return fmt.Errorf("could not write %q: %v", dest, err)
				}
				defer w.Close()

				// Call custom function
				return f(src, r, w)
			}()

			// Try to return error
			select {
			case <-ctx.Done():
				errch <- ctx.Err()
			case errch <- err:
			}
		}()

		return nil
	})

	go func() {
		wg.Wait()
		close(errch)
	}()

	for e := range errch {
		err = errors.Join(err, e)
		if err != nil {
			cancel()
		}
	}

	return err
}

func findPos(contents []byte, pos int) (line, ch int) {
	var acc int
	for i, line := range bytes.Split(contents, []byte("\n")) {
		if acc+len(line) > pos {
			return i, pos - acc
		}
		acc += len(line)
	}
	return 0, 0
}
