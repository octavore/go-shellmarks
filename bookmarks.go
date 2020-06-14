package main

import "sort"

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

func (c *Config) Print() {
	if len(c.Bookmarks) == 0 {
		errOut("(no bookmarks; add the current directory with `g -a`)")
	}
	sort.Slice(c.Bookmarks, func(i, j int) bool {
		return c.Bookmarks[i].Key < c.Bookmarks[j].Key
	})
	for _, b := range c.Bookmarks {
		errOut("%-12s %s", b.Key, b.Path)
	}
}
