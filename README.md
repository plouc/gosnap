# gosnap

[![Build Status](https://travis-ci.org/plouc/gosnap.png?branch=master)](https://travis-ci.org/plouc/gosnap)
[![GoDoc](https://godoc.org/github.com/plouc/gosnap?status.svg)](https://godoc.org/github.com/plouc/gosnap)
[![GitHub license](https://img.shields.io/github/license/plouc/gosnap.svg)](https://github.com/plouc/gosnap/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/plouc/gosnap)](https://goreportcard.com/report/github.com/plouc/gosnap)
[![GitHub issues](https://img.shields.io/github/issues/plouc/gosnap.svg)](https://github.com/plouc/gosnap/issues)

**gosnap** is a go package to perform snapshot testing.

Snapshot testing is a convenient way to test large outputs with ease,
this package was initially created to test the [`go-gitlab-client` CLI](https://github.com/plouc/go-gitlab-client).

## Install

```
go get github.com/plouc/gosnap
```

## Usage

**gosnap** requires a context to run, from which you can create snapshots.

````go
package whatever

import (
    "testing"
    "github.com/plouc/gosnap"
)

func TestSomething(t *testing.T) {
    // creates a new context
    // snapshots will be stored inside the `snapshots` directory
    ctx := gosnap.NewContext(t, "snapshots")

    // creates a new snapshot
    // it will be stored in `snapshots/my_first_snapshot.snap`
    s := ctx.NewSnapshot("my_first_snapshot")
    
    actual := DoSomethingWhichReturnsString()
    
    // this will load the snapshot content and check it matches `actual`
    s.AssertString(actual)
}
````

This will check that `actual` matches current snapshot (`./snapshots/my_first_snapshot.snap`) content.

The first time you run your tests, the snapshot will be created automatically,
then if the current result does not match snapshot's content, you'll have to
update it, you can add a command-line flag to the `go test command` to do so:

```
# will update all stale snapshots
go test -v -update all ./...

# will just update snapshot whose name is `my_snapshot`
go test -v -update my_snapshot ./...
```

For complete usage of **gosnap**, see the full [package docs](https://godoc.org/github.com/plouc/gosnap).
