package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/underscore"
)

var (
	flag_underscore *bool = flag.Bool("underscore", true, "Load underscore into the runtime environment")
	flag_repl       *bool = flag.Bool("i", false, "Start Otto in interactive mode (REPL)")
)

func readSource(filename string) ([]byte, error) {
	if filename == "" || filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filename)
}

func main() {
	flag.Parse()

	if !*flag_underscore {
		underscore.Disable()
	}

	err := func() error {
		vm := otto.New()
		if *flag_repl {
			return NewREPL(vm).Run()
		}
		src, err := readSource(flag.Arg(0))
		if err == nil {
			_, err = vm.Run(src)
		}
		return err
	}()
	if err != nil {
		switch err := err.(type) {
		case *otto.Error:
			fmt.Print(err.String())
		default:
			fmt.Println(err)
		}
		os.Exit(64)
	}
}
