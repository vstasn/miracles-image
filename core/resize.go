package core

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"

	"github.com/disintegration/imaging"
)

type Task struct {
	UserId       int    `json:"user_id"`
	Url          string `json:"url"`
	Filename     string `json:"filename"`
	Filename_50  string `json:"filename_50"`
	Filename_200 string `json:"filename_200"`
	Folder       string `json:"folder"`
}

type UploadImage struct {
	Body        io.Reader
	Key         string
	ContentType string
}

type ResizedImage struct {
	reader      io.Reader
	contentType string
}

func makePath(folder string, key string) string {
	return fmt.Sprintf("%s/%s", folder, key)
}

func DownloadImage(URL string) (image.Image, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("bad status: " + response.Status)
	}

	w := &bytes.Buffer{}
	_, err = io.Copy(w, response.Body)
	if err != nil {
		return nil, err
	}

	srcImage, err := imaging.Decode(w)
	if err != nil {
		return nil, err
	}

	return srcImage, nil
}

func resizeImage(img image.Image, width int, format imaging.Format) (output ResizedImage, err error) {
	src := imaging.Resize(img, width, 0, imaging.NearestNeighbor)

	var buf bytes.Buffer
	if err = imaging.Encode(&buf, src, format); err != nil {
		return
	}
	data := buf.Bytes()
	output.contentType = http.DetectContentType(data)
	output.reader = bytes.NewReader(data)
	return
}

func (t *Task) ResizeImages() (result []UploadImage, err error) {
	srcImg, err := DownloadImage(t.Url)
	if err != nil {
		return
	}

	if srcImg.Bounds().Max.X < 500 {
		err = errors.New("image is too small")
		return
	}

	format, err := imaging.FormatFromFilename(t.Filename)
	if err != nil {
		return nil, err
	}

	if img50, err := resizeImage(srcImg, 50, format); t.Filename_50 != "" && err == nil {
		result = append(result, UploadImage{
			Key:         makePath(t.Folder, t.Filename_50),
			Body:        img50.reader,
			ContentType: img50.contentType,
		})
	}

	if img200, err := resizeImage(srcImg, 200, format); t.Filename_200 != "" && err == nil {
		result = append(result, UploadImage{
			Key:         makePath(t.Folder, t.Filename_200),
			Body:        img200.reader,
			ContentType: img200.contentType,
		})
	}

	if img500, err := resizeImage(srcImg, 500, format); err == nil {
		result = append(result, UploadImage{
			Key:         makePath(t.Folder, t.Filename),
			Body:        img500.reader,
			ContentType: img500.contentType,
		})
	}

	return
}
