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
	h   []string
}

func NewREPL(vm *otto.Otto) repl {
	return repl{vm, term.NewTTY(os.Stdin), make([]string, 1)}
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
			io.WriteString(r.tty, "\r\nGoodbye!\r\n")
			return nil
		case term.CarriageReturn, term.NewLine:
			if expr = strings.Trim(expr, " \t"); len(expr) > 0 {
				if value, err := r.vm.Run(expr); nil != err {
					if err := perr(err); nil != err && err.Message == err_UnexpectedEndOfInput {
						// Unexpected end of input, continue and see if the next line completes
						// the expression.
						continue
					}
					io.WriteString(r.tty, err.Error()+"\r\n")
				} else {
					io.WriteString(r.tty, fmt.Sprintf(
						">>>>> %s\r\n", value.String()))
					r.h = append(r.h, expr)
				}
				expr = ""
			}
			r.prompt()
		default:
			if strings.IndexByte(in, ':') == 0 {
				switch cmd := strings.Split(in, " "); cmd[0] {
				case ":help", ":h", ":?":
					io.WriteString(r.tty, "Available commands:\r\n\r\n"+
						"\t:help, :h, or :? - Show this help message\r\n"+
						"\t:load <file>     - Load JavaScript file\r\n"+
						"\t:dump [<file>]   - Dump REPL history to <file> or stdout\r\n"+
						"\r\n")
				case ":load":
					if len(cmd) < 2 {
						io.WriteString(r.tty, "Please tell me what file to load.\r\n")
					}
					if src, err := readSource(cmd[1]); err == nil {
						if _, err = r.vm.Run(src); nil == err {
							continue
						}
					} else {
						io.WriteString(r.tty, err.Error()+"\r\n")
					}
				case ":dump":
					if err := dump(cmd, r.h); err != nil {
						io.WriteString(r.tty, err.Error()+"\r\n")
					}
				default:
					io.WriteString(r.tty, fmt.Sprintf(
						"Unknown command '%s'\r\n", cmd[0]))
				}
			} else {
				expr += in
			}
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

func dump(cmd []string, history []string) error {
	var w io.Writer = os.Stdout
	var nl string = "\r\n"
	if len(cmd) > 1 {
		fp, err := os.OpenFile(cmd[1], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer fp.Close()
		nl = "\n"
		w = fp
	}
	for _, l := range history {
		if _, err := io.WriteString(w, l+nl); err != nil {
			return err
		}
	}
	return nil
}
