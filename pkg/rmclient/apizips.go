package rmclient

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juruen/rmapi/archive"
	"github.com/pkg/errors"
)

func createDirectoryZip(id string) (io.Reader, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	w := zip.NewWriter(buf)

	f, err := w.Create(fmt.Sprintf("%s.content", id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create content file")
	}

	f.Write([]byte("{}"))
	w.Flush()
	w.Close()

	return bytes.NewReader(buf.Bytes()), nil
}

func createFileZip(id string, doc io.Reader) (io.Reader, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1048576))
	w := zip.NewWriter(buf)

	f, err := w.Create(fmt.Sprintf("%s.pdf", id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pdf entry")
	}
	_, err = io.Copy(f, doc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pipe document")
	}

	f, err = w.Create(fmt.Sprintf("%s.pagedata", id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pagedata entry")
	}
	f.Write(make([]byte, 0))

	f, err = w.Create(fmt.Sprintf("%s.content", id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create content entry")
	}

	c, err := createZipContent("pdf")
	if err != nil {
		return nil, err
	}

	f.Write([]byte(c))
	w.Flush()
	w.Close()

	return bytes.NewReader(buf.Bytes()), nil
}

func createZipContent(ext string) (string, error) {
	c := archive.Content{
		DummyDocument: false,
		ExtraMetadata: archive.ExtraMetadata{
			LastPen:             "Finelinerv2",
			LastTool:            "Finelinerv2",
			LastFinelinerv2Size: "1",
		},
		FileType:       ext,
		PageCount:      0,
		LastOpenedPage: 0,
		LineHeight:     -1,
		Margins:        180,
		TextScale:      1,
		Transform: archive.Transform{
			M11: 1,
			M12: 0,
			M13: 0,
			M21: 0,
			M22: 1,
			M23: 0,
			M31: 0,
			M32: 0,
			M33: 1,
		},
	}

	cstring, err := json.Marshal(c)

	if err != nil {
		return "", errors.Wrap(err, "failed to serialize content file")
	}

	return string(cstring), nil
}
