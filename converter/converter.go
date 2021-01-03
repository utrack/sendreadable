package converter

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/pkg/errors"
	"github.com/utrack/sendreadable/images"
	"github.com/utrack/sendreadable/tpl"
)

type Service struct {
	fontsPath string
	dirPrefix string
}

func New(fp string, dirPrefix string) *Service {
	return &Service{fontsPath: fp, dirPrefix: dirPrefix}
}

type Result struct {
	Filename    string
	ArticleName string
}

func (s *Service) Convert(ctx context.Context, url string) (*Result, error) {
	dir, err := ioutil.TempDir(s.dirPrefix, "*")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tempdir")
	}
	err = os.Chmod(dir, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "can't chmod directory")
	}

	go func() {
		<-time.After(time.Minute * 20)
		os.RemoveAll(dir)
	}()

	fmt.Println("pulling ", url, " , dir ", dir)

	a, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}

	dwn := images.NewRunner(dir)

	art, err := NewArticle(ctx, a, url, dwn)
	if err != nil {
		return nil, err
	}

	dst, err := os.OpenFile(filepath.Join(dir, "main.tex"), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tempfile")
	}
	defer dst.Close()

	tplReq := tpl.Request{
		Title:         art.Title,
		Author:        art.Author,
		URL:           art.URL,
		SourceName:    art.Source,
		AvgTimeString: art.AvgTimeString,
		Content:       art.Content,
		FontPath:      s.fontsPath,
	}
	err = tpl.Render(tplReq, dst)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute")
	}

	tmpName := dst.Name()
	err = dst.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close tmpfile")
	}

	cmd := exec.Command("xelatex", tmpName)
	cmd.Dir = dir
	logs, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "when running pdflatex, stderr: ("+string(logs)+")")
	}
	return &Result{
		Filename:    dir + "/main.pdf",
		ArticleName: art.Title}, nil
}
