package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/pkg/errors"
	"github.com/utrack/sendreadable/images"
	"github.com/utrack/sendreadable/tpl"
)

func runToPdf(ctx context.Context, url string) (string, error) {

	dir, err := ioutil.TempDir(os.TempDir()+"/sendreadable", "*")
	if err != nil {
		return "", errors.Wrap(err, "failed to create tempdir")
	}

	fmt.Println("pulling ", url, " , dir ", dir)

	a, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse")
	}

	dwn := images.NewRunner(dir)

	art, err := NewArticle(ctx, a, url, dwn)
	if err != nil {
		return "", err
	}

	if art.Image != "" {
		art.Image, err = dwn.Download(ctx, art.Image)
		if err != nil {
			log.Println("when pulling cover image: ", err.Error())
			art.Image = ""
		}
	}

	dst, err := ioutil.TempFile(dir, "main.tex")
	if err != nil {
		return "", errors.Wrap(err, "failed to create tempfile")
	}
	defer dst.Close()

	tplReq := tpl.Request{
		Title:         art.Title,
		Author:        art.Author,
		URL:           art.URL,
		SourceName:    art.Source,
		AvgTimeString: art.AvgTimeString,
		ImagePath:     art.Image,
		Content:       art.Content,
	}
	err = tpl.Render(tplReq, dst)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute")
	}

	tmpName := dst.Name()
	err = dst.Close()
	if err != nil {
		return "", errors.Wrap(err, "failed to close tmpfile")
	}

	cmd := exec.Command("xelatex", tmpName)
	cmd.Dir = dir
	logs, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "when running pdflatex, stderr: ("+string(logs)+")")
	}
	return dir, nil
}
