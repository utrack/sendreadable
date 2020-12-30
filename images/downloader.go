package images

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Runner struct {
	tickets chan struct{}

	cli *http.Client

	dir string
}

const parallelImagesPerJob = 5
const defaultTimeout = time.Second * 30

func NewRunner(dir string) *Runner {
	ch := make(chan struct{}, parallelImagesPerJob+1)
	for i := 0; i <= parallelImagesPerJob; i++ {
		ch <- struct{}{}
	}

	return &Runner{tickets: ch, cli: &http.Client{Timeout: defaultTimeout}}
}

func (r *Runner) Download(ctx context.Context, url string) (string, error) {
	<-r.tickets
	defer func() {
		r.tickets <- struct{}{}
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrapf(err, "cannot create GET request to '%v'", url)
	}
	req = req.WithContext(ctx)

	resp, err := r.cli.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "downloading image")
	}
	defer resp.Body.Close()

	f, err := ioutil.TempFile(r.dir, "image-*.png")
	if err != nil {
		return "", err
	}
	defer f.Close()

	return f.Name(), convertImageToPng(resp.Body, f)
}
