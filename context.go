package gosnap

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
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
	snapshots           map[string]*Snapshot
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
		snapshots:     map[string]*Snapshot{},
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

// shouldUpdateSnapshot checks if given snapshot should be automatically updated
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

// Has checks if a snapshot with the given name already exists
func (c *Context) Has(name string) bool {
	_, ok := c.snapshots[name]

	return ok
}

// Get returns a snapshot by its name
func (c *Context) Get(name string) *Snapshot {
	return c.snapshots[name]
}

// NewSnapshot creates a new snapshot attached to context.
// If a snapshot with the same name already exists, test will fail.
func (c *Context) NewSnapshot(name string) *Snapshot {
	if c.Has(name) {
		c.t.Errorf(color.RedString("snapshot %s already exists", name))
		c.t.FailNow()
		return nil
	}

	s := Snapshot{
		Name: name,
		ctx:  c,
	}

	c.snapshots[name] = &s

	return &s
}

func init() {
	flag.StringVar(&updateSnapshots, "update", "", "update generated snapshots")
	flag.StringVar(&updateSnapshots, "u", "", "update generated snapshots")

	flag.Parse()
}
