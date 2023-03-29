package finitestatemachine

import (
	"fmt"
	"go/token"
)

type expect int

const (
	pkg expect = iota
	pkgPeriod
	fun
	lParen
	str
	rParen
)

type Match struct {
	Pos int
	Val string
}

type Machine struct {
	expect   expect
	pkgName  string
	funcName string

	Results []Match
}

func New(pkgName, funcName string) Machine {
	return Machine{
		pkgName:  pkgName,
		funcName: funcName,
	}
}

func (m *Machine) Consume(pos token.Pos, tk token.Token, lit string) error {
	switch m.expect {
	case pkg:
		if tk == token.IDENT && lit == m.pkgName {
			m.expect = pkgPeriod
			return nil
		}
	case pkgPeriod:
		if tk == token.PERIOD {
			m.expect = fun
			return nil
		}
	case fun:
		if tk == token.IDENT && lit == m.funcName {
			m.expect = lParen
			return nil
		}
	case lParen:
		if tk == token.LPAREN {
			m.expect = str
			return nil
		}
	case str:
		if tk == token.STRING {
			m.Results = append(m.Results, Match{int(pos), lit})
		} else {
			return fmt.Errorf("unexpected token. Expected string, got token %s with value %s", tk.String(), lit)
		}
	}

	m.expect = pkg
	return nil
}
