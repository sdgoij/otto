package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/kylelemons/goat/term"
	"github.com/kylelemons/goat/termios"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
)

const err_UnexpectedEndOfInput = "Unexpected end of input"

type repl struct {
	vm  *otto.Otto
	tty *term.TTY
}

func NewREPL(vm *otto.Otto) repl {
	return repl{vm, term.NewTTY(os.Stdin)}
}

func (r repl) Run() error {
	tio, err := termios.NewTermSettings(0)
	if err != nil {
		return err
	}
	if err := tio.Raw(); err != nil {
		return err
	}
	defer tio.Reset()

	linebuf := make([]byte, 128)
	expr := ""

	r.welcome()
	r.prompt()

	for {
		n, err := r.tty.Read(linebuf)
		if err != nil {
			return err
		}
		switch in := string(linebuf[:n]); in {
		// Quit on ":quit", ^C, and ^D
		case ":quit", term.Interrupt, term.EndOfFile:
			io.WriteString(os.Stdout, "\r\nGoodbye!\r\n")
			return nil
		case term.CarriageReturn, term.NewLine:
			if r, err := r.vm.Run(strings.Trim(expr, " \t")); nil != err {
				if err := perr(err); nil != err && err.Message == err_UnexpectedEndOfInput {
					// Unexpected end of input, continue and see if the next line completes
					// the expression.
					continue
				}
				fmt.Printf("%s (%T)\r\n", err, err)
			} else {
				fmt.Println(">>>>>", r, "\r")
			}
			expr = ""
			r.prompt()
		default:
			expr += in
		}
	}
}

func (r repl) prompt() {
	io.WriteString(r.tty, "otto> ")
}

func (r repl) welcome() {
	io.WriteString(r.tty, fmt.Sprintf(
		"Welcome to otto (Go version: %s; %s %s %s) Copyright (c) 2012 Robert Krimen\r\n"+
			"Type \":quit\", ^C, or ^D to exit\r\n\r\n",
		runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH))
}

func perr(e error) (err *parser.Error) {
	switch t := e.(type) {
	case parser.ErrorList:
		err = t[0]
	case parser.Error:
		err = &t
	case *parser.Error:
		err = t
	}
	return
}
