package dirparse

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// WalkDir is a walker through the source dir(also visits nested dirs) which
// copies jpeg images with exif data to the destination folder with name format
// dstFlder/YYYY/YYYY-MM-DD--HH-MM-SS-old-fileame.jpg
type WalkDir struct {
	srcDir, dstDir string
	appFs          afero.Fs
}

// NewWalkDir creates new instance of WalkDir
func NewWalkDir(appFs afero.Fs, srcDir, dstDir string) *WalkDir {
	return &WalkDir{
		srcDir: srcDir,
		dstDir: dstDir,
		appFs:  appFs,
	}
}

// Walk is actual method which does the passthrough and copy of the images
func (wd *WalkDir) Walk() error {
	if err := afero.Walk(wd.appFs, wd.srcDir, wd.walkfn); err != nil {
		return errors.Wrap(err, "failed to walk dir")
	}

	return nil
}

func (wd *WalkDir) walkfn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if !isRelevantPath(path) {
		return nil
	}

	if dstFilepath, err := wd.processFile(path); err != nil {
		log.Errorf("failed processing file: %s %q", path, err)
	} else {
		log.WithFields(log.Fields{"source": path, "destination": dstFilepath}).Info("image copied")
	}

	return nil
}

func (wd *WalkDir) processFile(filepath string) (dstFilepath string, err error) {
	file, err := wd.appFs.Open(filepath)
	if err != nil {
		return "", errors.Wrapf(err, "failed opening file %s", filepath)
	}
	defer file.Close()

	exifData, err := exif.Decode(file)
	if err != nil {
		return "", errors.Wrapf(err, "failed parsing EXIF data %s", filepath)
	}
	dateTime, err := exifData.DateTime()
	if err != nil {
		return "", errors.Wrapf(err, "failed getting image's datatime %s", filepath)
	}

	newFilepath := newFilePathFromDateTime(dateTime, filepath, wd.dstDir)
	// need to reset file pointer since it is already not at the beginnig after reading exif data
	file.Seek(0, 0)
	if err := copyFile(wd.appFs, file, newFilepath); err != nil {
		return "", errors.Wrapf(err, "failed to copy image %s", filepath)
	}

	return newFilepath, nil
}

func newFilePathFromDateTime(dt time.Time, oldFilepath string, dstDir string) string {
	return fmt.Sprintf(
		"%s/%s-%s",
		dstDir,
		dt.Format("2006/2006-01-02--15-04-05"),
		filepath.Base(oldFilepath))
}

func isRelevantPath(path string) bool {
	fileExt := strings.ToLower(filepath.Ext(path))

	return fileExt == ".jpg" || fileExt == ".jpeg"
}

func copyFile(appFs afero.Fs, srcFile io.Reader, dstFilepath string) error {
	if _, err := appFs.Stat(dstFilepath); err == nil {
		return errors.New("skipping. Destination file already exists: " + dstFilepath)
	}

	if err := appFs.MkdirAll(filepath.Dir(dstFilepath), 0744); err != nil {
		return errors.New("failed to create dir structure: " + dstFilepath)
	}

	content, err := afero.ReadAll(srcFile)
	if err != nil {
		return errors.Wrapf(err, "error reading file %s", dstFilepath)
	}

	if err := afero.WriteFile(appFs, dstFilepath, content, 0744); err != nil {
		return errors.Wrapf(err, "failed to writДe to file %s", dstFilepath)
	}

	return nil
}
