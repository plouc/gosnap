package gosnap

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Snapshot represents a snapshot file
type Snapshot struct {
	Name           string
	content        string
	ctx            *Context
	hasBeenLoaded  bool
	hasBeenUpdated bool
}

// FileName returns snapshot file name
func (s *Snapshot) FileName() string {
	return fmt.Sprintf("%s%s", s.Name, s.ctx.FileExtension)
}

// FilePath returns full snapshot path on disk
func (s *Snapshot) FilePath() string {
	return filepath.Join(s.ctx.Dir, s.FileName())
}

// File returns snapshot file for reading
func (s *Snapshot) File() (*os.File, error) {
	file, err := os.Open(s.FileName())
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Load loads snapshot file content
func (s *Snapshot) Load() error {
	err := s.ctx.ensureDir()
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(s.FilePath())
	if err == nil {
		s.content = string(content)
		s.hasBeenLoaded = true
	}

	return err
}

// Content returns snapshot content and loads it if required
func (s *Snapshot) Content() (string, error) {
	if s.hasBeenLoaded || s.hasBeenUpdated {
		return s.content, nil
	}

	err := s.Load()
	if err != nil {
		return "", err
	}

	return s.content, nil
}

// Update writes given content to disk and refresh content
func (s *Snapshot) Update(c string) error {
	err := s.ctx.ensureDir()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.FilePath(), []byte(c), s.ctx.FileMode)
	if err != nil {
		return err
	}

	s.content = c
	s.hasBeenLoaded = true

	return nil
}

// AssertString test given string against stored content
func (s *Snapshot) AssertString(expected string) {
	c, err := s.Content()
	if err != nil {
		s.ctx.t.Error(err)
		s.ctx.t.FailNow()
	}

	if c != expected {
		s.ctx.t.Errorf("Snapshot does not match: %s\n%s", s.FilePath(), StringsDiff(expected, c))
	}
}
