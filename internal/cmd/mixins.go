package cmd

import (
	"path/filepath"
)

type GoFileMixin struct {
	GoFilePath string
}
func (m *GoFileMixin) GetGofilePath() (path string, err error) {
	abs, err := filepath.Abs(m.GoFilePath)
	if err != nil {
		return "", err
	}
	return abs, nil
}
