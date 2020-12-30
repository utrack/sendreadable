package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func pullImage(url string, timeout time.Duration, dir string) (string, error) {
	fmt.Println("pulling image ", url)
	// Fetch page from URL
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "downloading image")
	}
	defer resp.Body.Close()

	f, err := ioutil.TempFile(dir, "image-*.png")
	if err != nil {
		return "", err
	}
	defer f.Close()

	return f.Name(), convertImageToPng(resp.Body, f)
}

func convertImageToPng(r io.Reader, w io.Writer) error {
	r = io.LimitReader(r, 2097152)
	imgData, _, err := image.Decode(r)
	if err != nil {
		return errors.Wrap(err, "reading image")
	}

	err = png.Encode(w, imgData)
	return err
}
