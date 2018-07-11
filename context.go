package gosnap

import (
	"fmt"
	"os"
	"testing"
)

// Context is used to isolate snapshot testing for given options
type Context struct {
	t             *testing.T
	Dir           string
	DirMode       os.FileMode
	FileMode      os.FileMode
	FileExtension string
	AutoUpdate    bool
}

// NewContext creates a new snapshot testing context
func NewContext(t *testing.T, d string) *Context {
	return &Context{
		t:             t,
		Dir:           d,
		DirMode:       0755,
		FileMode:      0644,
		FileExtension: ".snap",
		AutoUpdate:    false,
	}
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

// NewSnapshot creates a new snapshot attached to context
func (c *Context) NewSnapshot(name string) *Snapshot {
	return &Snapshot{
		Name: name,
		ctx:  c,
	}
}
