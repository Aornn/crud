package filestorer_test

import (
	"crud/internal/module/filestorer"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteUser(t *testing.T) {
	fs := filestorer.NewFileStorer()
	filename := "file"
	t.Run("WriteUser valid", func(t *testing.T) {
		err := fs.WriteUser(filename, "data")
		assert.Nil(t, err)
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
		if err := os.Remove(filename); err != nil {
			panic(err)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	fs := filestorer.NewFileStorer()
	filename := "file"
	t.Run("Delete valid", func(t *testing.T) {
		_, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		err = fs.DeleteUser(filename)
		assert.Nil(t, err)
		if err != nil && err.Error() == "file not found" {
			if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
				panic(err)
			}
			if err := os.Remove(filename); err != nil {
				panic(err)
			}
		}
	})

	t.Run("Delete returns error", func(t *testing.T) {
		err := fs.DeleteUser(filename)
		assert.EqualError(t, err, "file not found")
	})
}

func TestUpdateUser(t *testing.T) {
	fs := filestorer.NewFileStorer()
	filename := "file"
	t.Run("UpdateUser file not found", func(t *testing.T) {
		err := fs.UpdateUser(filename, "data")
		assert.EqualError(t, err, "file not found")
	})

	t.Run("Valid", func(t *testing.T) {
		_, err := os.Create(filename)
		err = fs.UpdateUser(filename, "data")
		assert.Nil(t, err)
		if err != nil && err.Error() != "file not found" {
			if _, err = os.Stat(filename); errors.Is(err, os.ErrNotExist) {
				panic(err)
			}
			if err := os.Remove(filename); err != nil {
				panic(err)
			}
		}
		if err == nil {
			if err := os.Remove(filename); err != nil {
				panic(err)
			}
		}
	})
}
