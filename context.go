package gosnap

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

var updateSnapshots string

// Context is used to isolate snapshot testing for given options
type Context struct {
	t                   *testing.T
	Dir                 string
	DirMode             os.FileMode
	FileMode            os.FileMode
	FileExtension       string
	AutoUpdate          bool
	AutoUpdateSnapshots []string
}

// NewContext creates a new snapshot testing context
//
// By default the `AutoUpdate` will be enabled according
// to `-update` or `-u` command-line flag, if it equals `all`
// `AutoUpdate` will be enabled.
//
// If the update flag is set but doesn't equal `all`,
// then its value will be split using a comma,
// this allows to update specific snapshots only,
// if the snapshot's name matches one of the extracted value
// it will be updated, otherwise code has to be fixed :)
//
// Note that even if it initially uses flag value,
// you can manually set `AutoUpdate` & `AutoUpdateSnapshots`
// manually.
func NewContext(t *testing.T, d string) *Context {
	ctx := Context{
		t:             t,
		Dir:           d,
		DirMode:       0755,
		FileMode:      0644,
		FileExtension: ".snap",
		AutoUpdate:    updateSnapshots == "all",
	}

	if updateSnapshots != "" && updateSnapshots != "all" {
		ctx.AutoUpdateSnapshots = strings.Split(updateSnapshots, ",")
	}

	return &ctx
}

func (c *Context) ensureDir() error {
	s, err := os.Stat(c.Dir)
	switch {
	case err != nil && os.IsNotExist(err):
		err = os.MkdirAll(c.Dir, c.DirMode)
		if err != nil {
			return err
		}

		return err

	case err == nil && !s.IsDir():
		return fmt.Errorf("%s is not a directory", c.Dir)

	default:
		return err
	}
}

func (c *Context) shouldUpdateSnapshot(s *Snapshot) bool {
	if c.AutoUpdate {
		return true
	}

	for _, name := range c.AutoUpdateSnapshots {
		if name == s.Name {
			return true
		}
	}

	return false
}

// NewSnapshot creates a new snapshot attached to context
func (c *Context) NewSnapshot(name string) *Snapshot {
	return &Snapshot{
		Name: name,
		ctx:  c,
	}
}

func init() {
	flag.StringVar(&updateSnapshots, "update", "whatever", "do whatever")
	flag.StringVar(&updateSnapshots, "u", "whatever", "do whatever")

	flag.Parse()
}
