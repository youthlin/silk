package main_test

import (
	"io/fs"
	"path/filepath"
	"testing"
)

func TestWalk(t *testing.T) {
	filepath.WalkDir("..", func(path string, d fs.DirEntry, err error) error {
		t.Logf("err=%+v, path=%v, d=%v", err, path, d)
		return nil
	})

}
