package helper

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func FileUpload(file *multipart.FileHeader) (newFileName string, err error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	extension := filepath.Ext(file.Filename)
	storagePath := fmt.Sprintf("%s/../storage/", GetRelativeDirPath())
	newFileName = fmt.Sprintf("%s%s", GenerateUuid(), extension)
	filePath := fmt.Sprintf("%s%s", storagePath, newFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return newFileName, nil
}

func FileDelete(fileName string) error {
	if fileName == "" {
		return errors.New("invalid file name")
	}

	storagePath := fmt.Sprintf("%s/../storage/", GetRelativeDirPath())
	filePath := fmt.Sprintf("%s%s", storagePath, fileName)
	err := os.Remove(filePath)

	return err
}
