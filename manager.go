package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type bookmarksManager struct {
	searchPaths []string
}

// newBookmarksManager sets up all the search paths and bookmarksPath
func newBookmarksManager() (*bookmarksManager, error) {
	cm := &bookmarksManager{}
	c := os.Getenv("XDG_CONFIG_HOME")
	if c != "" {
		cm.searchPaths = append(cm.searchPaths, path.Join(c, appName))
	}
	cm.searchPaths = append(cm.searchPaths, getHomeConfigDir())
	return cm, nil
}

// bookmarksPath() returns the active config file path, config file dir
// and whether it exists or not.
// This checks all search paths for an existing config file
// XDG_CONFIG_HOME is the preferred path, but also fallback
// gracefully to $HOME/.config
func (cm *bookmarksManager) bookmarksPath() (string, string, bool) {
	// default config path is the first search path
	bookmarksDir := cm.searchPaths[0]

	for _, dir := range cm.searchPaths {
		bookmarksPath := path.Join(dir, "bookmarks.json")
		fi, err := os.Stat(bookmarksPath)
		if fi != nil && err == nil {
			return bookmarksPath, dir, true
		}
		if !os.IsNotExist(err) {
			errOut("Unknown error: %s", err.Error())
			os.Exit(1)
		}
	}
	return path.Join(bookmarksDir, "bookmarks.json"), bookmarksDir, false
}

func (cm *bookmarksManager) ensureAndLoad() (*Config, error) {
	err := cm.ensure()
	if err != nil {
		return nil, err
	}
	bookmarksPath, _, exists := cm.bookmarksPath()
	config := &Config{}
	if exists {
		f, err := ioutil.ReadFile(bookmarksPath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(f, config)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

// writeBookmarks writes to the existing config file, or the first search path
func (cm *bookmarksManager) writeBookmarks(config *Config) error {
	bookmarksPath, _, _ := cm.bookmarksPath()
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(bookmarksPath, b, 0644)
}

func (cm *bookmarksManager) ensure() error {
	bookmarksPath, bookmarksDir, exists := cm.bookmarksPath()
	if exists {
		return nil
	} else {
		err := os.MkdirAll(bookmarksDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create dir %s: %s", bookmarksDir, err)
		}
	}
	err := cm.writeBookmarks(&Config{})
	if err != nil {
		return fmt.Errorf("failed to to create config.json file: %s", err)
	}
	fmt.Printf("created config.json file: %s\n", bookmarksPath)
	return nil
}

func getHomeConfigDir() string {
	u, err := user.Current()
	if uid := os.Getenv("SUDO_UID"); uid != "" {
		u, err = user.LookupId(uid)
	}
	if err != nil {
		panic(err)
	}
	return path.Join(u.HomeDir, ".config", appName)
}
