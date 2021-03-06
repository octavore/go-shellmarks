package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/fatih/color"
)

type Bookmark struct {
	Key  string `json:"key"`
	Path string `json:"path"`
}

type Config struct {
	Bookmarks []*Bookmark `json:"bookmarks"`
}

func (c *Config) Get(key string) (string, bool) {
	for _, b := range c.Bookmarks {
		if b.Key == key {
			return b.Path, true
		}
	}
	return "", false
}

func (c *Config) Add(key, path string) {
	for _, b := range c.Bookmarks {
		if b.Key == key {
			b.Path = path
			return
		}
	}
	c.Bookmarks = append(c.Bookmarks, &Bookmark{
		Key:  key,
		Path: path,
	})
}

func (c *Config) Remove(key string) {
	updated := []*Bookmark{}
	for _, b := range c.Bookmarks {
		if b.Key != key {
			updated = append(updated, b)
		}
	}
	c.Bookmarks = updated
}

func (c *Config) Select() (string, error) {
	bookmarks, err := c.sorted()
	if err != nil {
		return "", err
	}
	bold := color.New(color.FgYellow, color.Bold)
	bold.EnableColor()

	// create question
	q := &survey.Question{
		Prompt: &survey.Select{Message: "Choose a shortcut"},
		Transform: func(s interface{}) interface{} {
			ans := s.(core.OptionAnswer)
			return core.OptionAnswer{
				Index: ans.Index,
				Value: bookmarks[ans.Index].Path,
			}
		},
	}
	for _, b := range bookmarks {
		o := q.Prompt.(*survey.Select).Options
		o = append(o, fmt.Sprintf("%s %s", bold.Sprintf("%-12s", b.Key), b.Path))
		q.Prompt.(*survey.Select).Options = o
	}

	// run select
	opt := ""
	err = survey.Ask([]*survey.Question{q}, &opt, survey.WithStdio(os.Stdin, os.Stderr, os.Stderr))
	if err != nil {
		return "", err
	}
	if opt == "" {
		errOut("nothing selected")
	}
	return opt, nil
}

func (c *Config) Print() error {
	bold := color.New(color.FgYellow, color.Bold)
	bold.EnableColor()
	bookmarks, err := c.sorted()
	if err != nil {
		return err
	}
	for _, b := range bookmarks {
		errOut("%s %s", bold.Sprintf("%-12s", b.Key), b.Path)
	}
	return nil
}

func (c *Config) sorted() ([]*Bookmark, error) {
	bookmarks := c.Bookmarks
	if len(bookmarks) == 0 {
		return nil, fmt.Errorf("(no bookmarks; add the current directory with `g -a`)")
	}
	sort.Slice(bookmarks, func(i, j int) bool {
		return bookmarks[i].Key < bookmarks[j].Key
	})
	return bookmarks, nil
}
