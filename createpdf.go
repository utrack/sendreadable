package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/pkg/errors"
)

var tpl = template.Must(template.New("").Delims("$$", "$$").Parse(tplLatex))

func runToPdf(url string) error {

	fmt.Println("pulling ", url)

	a, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return errors.Wrap(err, "failed to parse")
	}

	art, err := NewArticle(a, url)
	if err != nil {
		return err
	}

	//fmt.Println(art.Content)
	//return nil

	dir, err := ioutil.TempDir(os.TempDir(), "sendreadable-*")
	if err != nil {
		return errors.Wrap(err, "failed to create tempdir")
	}
	defer os.RemoveAll(dir)

	if art.Image != "" {
		art.Image, err = pullImage(art.Image, time.Second*10, dir)
		if err != nil {
			log.Println("when pulling cover image: ", err.Error())
			art.Image = ""
		}
	}

	dst, err := ioutil.TempFile(dir, "main.tex")
	if err != nil {
		return errors.Wrap(err, "failed to create tempfile")
	}
	defer dst.Close()

	err = tpl.Execute(dst, art)
	if err != nil {
		return errors.Wrap(err, "failed to execute")
	}

	tmpName := dst.Name()
	err = dst.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close tmpfile")
	}

	cmd := exec.Command("xelatex", tmpName)
	cmd.Dir = dir
	fmt.Printf("running '%v'\n", cmd.String())
	logs, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "when running pdflatex("+string(logs)+")")
	}
	fmt.Println(string(logs))
	fmt.Println(dir + "/main.pdf")
	<-time.After(time.Minute * 5)
	return nil
}
