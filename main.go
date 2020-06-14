package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const appName = "go-shellmarks"

const bashFuncTmpl = `function g {
	target="$(%s $*)"
	[[ $? != 0 ]] && return
	if [ -z "$target" ]; then
		return
	elif [ -d "$target" ]; then
		cd "$target"
	else
		echo target "$target" does not exist
	fi
}
`

func main() {
	b, err := newBookmarksManager()
	exitIfError(err)

	bookmarks, err := b.ensureAndLoad()
	exitIfError(err)

	addBookmark := flag.String("a", "", "add a new bookmark to the current folder")
	delBookmark := flag.String("d", "", "delete a bookmark")
	_ = flag.Bool("l", false, "list bookmarks")
	_ = flag.Bool("shell", false, "add this to your .bashrc: source $(g --shell)")
	_ = flag.Bool("h", false, "show help")
	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "a":
			path, err := os.Getwd()
			exitIfError(err)
			errOut("Add bookmark: %s -> %s", *addBookmark, path)
			bookmarks.Add(*addBookmark, path)
			err = b.writeBookmarks(bookmarks)
			exitIfError(err)
		case "d":
			errOut("Remove bookmark: %s", *delBookmark)
			bookmarks.Remove(*delBookmark)
			err := b.writeBookmarks(bookmarks)
			exitIfError(err)
		case "shell":
			binPath, err := filepath.Abs(os.Args[0])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(bashFuncTmpl, binPath)
		case "l":
			bookmarks.Print()
		case "h":
			bookmarkPath, _, _ := b.bookmarksPath()
			errOut("Bookmark path: %s\n", bookmarkPath)
			flag.Usage()
		}
		os.Exit(0)
	})

	switch {
	case len(os.Args) > 1:
		key := os.Args[1]
		path, ok := bookmarks.Get(key)
		if !ok {
			errOut("Unknown bookmark: %s", key)
			os.Exit(1)
		}
		fmt.Println(path)
		os.Exit(0)
	case len(os.Args) == 1:
		bookmarks.Print()
	}
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
