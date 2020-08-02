package main

import (
	"fmt"
	"os"
)

const appName = "go-shellmarks"

const bashFuncTmpl = `function %s {
	target="$(SHELLMARKS_ALIAS=%s %s $*)"
	[[ $? != 0 ]] && return
	if [ -z "$target" ]; then
		return
	elif [ -d "$target" ]; then
		echo cd $target
		cd "$target"
	else
		echo target "$target" does not exist
	fi
}
`

func main() {
	Commands()
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
