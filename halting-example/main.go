package main

import (
	Otto ".."
	"fmt"
	"os"
	"time"
)

func main() {
	otto := Otto.New()
	_, err := otto.RunUntil(`
while (true) {
    // Loop forever
}
    `, time.Second*10)
	fmt.Fprintln(os.Stderr, err)
}
