package filestorer

import (
	"errors"
	"fmt"
	"os"
)

//go:generate mockgen -destination mock_test.go -source $GOFILE -package filestorer_test

type FileStorer struct{}

func NewFileStorer() *FileStorer {
	return &FileStorer{}
}

func (f *FileStorer) Checkfile(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (f *FileStorer) WriteUser(filename, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileStorer) DeleteUser(id string) error {
	if _, err := os.Stat(id); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file not found")
	}
	return os.Remove(id)
}

func (f *FileStorer) UpdateUser(filename, data string) error {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file not found")
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	if string(content) != data {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(data)
		if err != nil {
			return err
		}
	}
	return nil
}
