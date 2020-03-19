package vgdir

import (
	"os"
	"path/filepath"
)

// ConvDirAbs 相对路径转换绝对路径
func ConvDirAbs(dir string) string {
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return dir
	}
	return filepath.Join(workDir, dir)
}
