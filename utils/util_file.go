package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

/** Move upload file to a folder. */
func MoveFile(uploadFile multipart.File, newPath string) error {
	/* Condition validation */
	if uploadFile == nil {
		return fmt.Errorf("Upload file is empty.")
	}
	defer uploadFile.Close()

	// Create new file
	output, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer output.Close()

	// Move file to new location
	_, err = io.Copy(output, uploadFile)
	return err
}

/** Check if file exist at path or not. */
func FileExisted(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
