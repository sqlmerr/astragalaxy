package testutil

import (
	"path/filepath"
	"runtime"
)

func ProjectRoot() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(currentFile), "../..")
}
