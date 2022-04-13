package celeritas

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/brucebotes/celeritas/filesystems"
	"github.com/gabriel-vasile/mimetype"
)

func (c *Celeritas) UploadFile(r *http.Request, destination, field string, fs filesystems.FS) error {
	fileName, err := c.getFileToUpload(r, field)
	if err != nil {
		c.ErrorLog.Println(err)
		return err
	}

	if fs != nil {
		// we are uploading to a remote filesystem
		err = fs.Put(fileName, destination)
		if err != nil {
			c.ErrorLog.Println(err)
			return err
		}
	} else {
		os.Rename(fileName, fmt.Sprintf("%s/%s", destination, path.Base(fileName)))
		if err != nil {
			c.ErrorLog.Println(err)
			return err
		}
	}

	// delete the file that was uploaded into the tmp folder
	defer func() {
		_ = os.Remove(fileName)
	}()

	return nil
}

func (c *Celeritas) getFileToUpload(r *http.Request, fieldName string) (string, error) {
	_ = r.ParseMultipartForm(c.config.uploads.maxUploadSize)

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return "", err
	}

	// return file reader to start from beginning of file
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	if !inSlice(c.config.uploads.allowedMimeTypes, mimeType.String()) {
		return "", errors.New("invalid file type uploaded")
	}

	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}

func inSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
