package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestCheckFlagNotSet(t *testing.T) {
	err := checkFlag("")

	if err == nil || err.Error() != "both of the flags need to be set: src and dst" {
		t.Error("checkFlag should return an error for empty flag")
	}
}

func TestCheckFlagDirNoExist(t *testing.T) {
	appFs = afero.NewMemMapFs()
	err := checkFlag("/test")

	if err == nil || err.Error() != "Dir does not exist: /test" {
		t.Error("checkFlag should return an error for non existing dir")
	}
}

func TestCheckFlagNotADir(t *testing.T) {
	appFs = afero.NewMemMapFs()
	afero.WriteFile(appFs, "/file", []byte("file c"), 0644)
	err := checkFlag("/file")

	if err == nil || err.Error() != "Not a dir: /file" {
		t.Error("checkFlag should return an error for non existing dir")
	}
}

func TestCheckFlagExistingADir(t *testing.T) {
	appFs = afero.NewMemMapFs()
	appFs.Mkdir("/testDir", 0644)
	err := checkFlag("/testDir")

	if err != nil {
		t.Error("checkFlag should not return an error for existing dir")
	}
}
