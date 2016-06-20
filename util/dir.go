package util

import "os"

/** Create directories. */
func CreateDir(path string, perm os.FileMode) {
	os.MkdirAll(path, perm)
}

/** Check if file exist at path or not. */
func DirExisted(dirPath string) bool {
	fileInfo, err := os.Stat(dirPath)
	return err == nil && fileInfo.IsDir()
}
