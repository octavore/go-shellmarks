package main

import (
	"fmt"
	"os"
)

func main() {
	runCommands()
}

func exitIfError(err error) {
	if err != nil {
		errOut("error: %s", err.Error())
		os.Exit(1)
	}
}

func errOut(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
