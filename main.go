package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/xdrive/photosort/dirparse"
)

var appFs = afero.NewOsFs()

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	srcDir := flag.String("src", "", "path to the source dir with images")
	dstDir := flag.String("dst", "", "path to the destination dir")
	flag.Parse()

	if err := checkFlag(*srcDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := checkFlag(*dstDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wd := dirparse.NewWalkDir(appFs, *srcDir, *dstDir)
	if err := wd.Walk(); err != nil {
		log.Errorf("got an error when walking dir %q", err)
		fmt.Println(err)
	} else {
		fmt.Println("done!")
	}
}

func checkFlag(dir string) error {
	if dir == "" {
		return errors.New("both of the flags need to be set: src and dst")
	}

	stat, err := appFs.Stat(dir)

	if os.IsNotExist(err) {
		return errors.New("Dir does not exist: " + dir)
	}

	if !stat.IsDir() {
		return errors.New("Not a dir: " + dir)
	}

	return nil
}
