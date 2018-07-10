package gosnap

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/src-d/go-git.v4/utils/diff"
)

// StringsDiff returns a pretty diff between given strings
func StringsDiff(a, b string) string {
	ds := []string{}
	diffs := diff.Do(a, b)
	for _, d := range diffs {
		if d.Type != diffmatchpatch.DiffEqual {
			if d.Type == diffmatchpatch.DiffDelete {
				ds = append(ds, fmt.Sprintf("- %s", d.Text))
			} else {
				ds = append(ds, fmt.Sprintf("+ %s", d.Text))
			}
		}
	}

	return strings.Join(ds, "")
}
