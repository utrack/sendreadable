package images

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"

	"github.com/pkg/errors"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func convertImageToPng(r io.Reader, w io.Writer) error {
	r = io.LimitReader(r, 2097152)
	imgData, _, err := image.Decode(r)
	if err != nil {
		return errors.Wrap(err, "reading image")
	}

	err = png.Encode(w, imgData)
	return err
}
