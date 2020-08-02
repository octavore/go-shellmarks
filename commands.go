package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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

func runCommands() {
	b, err := newBookmarksManager()
	exitIfError(err)

	bookmarks, err := b.ensureAndLoad()
	exitIfError(err)

	rootCmd := &cobra.Command{
		Use:   "go-shellmarks",
		Short: "go-shellmarks is bookmarks for your shell",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			_, isAliased := os.LookupEnv("SHELLMARKS_ALIAS")
			if !isAliased {
				helpFunc(cmd, args)
				return
			}
			var p string
			var err error
			if len(args) == 0 {
				p, err = bookmarks.Select()
			} else {
				key := args[0]
				_p, ok := bookmarks.Get(key)
				if !ok {
					err = fmt.Errorf("Unknown bookmark: %s", key)
				}
				p = _p
			}
			exitIfError(err)
			fmt.Println(p)
		},
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "@get <bookmark>",
		Short: "show stored value for <bookmark>",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			p, ok := bookmarks.Get(key)
			if !ok {
				exitIfError(fmt.Errorf("Unknown bookmark: %s", key))
			}
			errOut(p)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "@add <bookmark>",
		Short: "bookmark the current directory as <bookmark>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bookmark := args[0]
			path, err := os.Getwd()
			if err != nil {
				return err
			}
			errOut("Add bookmark: %s -> %s", bookmark, path)
			bookmarks.Add(bookmark, path)
			return b.writeBookmarks(bookmarks)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "@rm <bookmark>",
		Short: "remove a stored bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bookmark := args[0]
			errOut("Remove bookmark: %s", bookmark)
			bookmarks.Remove(bookmark)
			return b.writeBookmarks(bookmarks)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "@ls",
		Short: "list all bookmarks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bookmarks.Print()
		},
	})

	var g *string
	shellCmd := &cobra.Command{
		Use:   "@shell",
		Short: "go-shellmarks is bookmarks for your shell",
		RunE: func(cmd *cobra.Command, args []string) error {
			binPath, err := filepath.Abs(os.Args[0])
			if err != nil {
				return err
			}
			fmt.Printf(bashFuncTmpl, *g, *g, binPath)
			return nil
		},
	}
	g = shellCmd.Flags().StringP("alias", "a", "g", "alias for shell command")
	rootCmd.AddCommand(shellCmd)

	rootCmd.SetHelpFunc(helpFunc)
	_ = rootCmd.Execute()
}

func helpFunc(cmd *cobra.Command, args []string) {
	usage := "\n" + cmd.UsageString()
	alias, ok := os.LookupEnv("SHELLMARKS_ALIAS")
	green := color.New(color.FgGreen)
	green.EnableColor()
	if ok {
		usage = strings.ReplaceAll(usage,
			"  go-shellmarks [flags]",
			"  go-shellmarks [bookmark]\n  go-shellmarks [flags]")
		usage = strings.ReplaceAll(usage, "go-shellmarks", fmt.Sprintf("%s", alias))
		usage = fmt.Sprintf("\nAlias:\n  go-shellmarks is currently aliased to: %s\n%s", green.Sprint(alias), usage)
	} else {
		p, _ := filepath.Abs(os.Args[0])
		aliasHelp := green.Sprintf("source <(%s @shell --alias z)", p)
		usage = fmt.Sprintf("\nSetup:\n  Add this to your .bashrc to use go-shellmarks as z:\n  %s\n%s", aliasHelp, usage)
	}
	usage = fmt.Sprintf("go-shellmarks: bookmarks for your shell\n%s", usage)
	errOut(usage)
}
