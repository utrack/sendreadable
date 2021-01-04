package images

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	tickets chan struct{}

	cli *http.Client

	dir string

	seedImg []byte

	wg sync.WaitGroup
}

const parallelImagesPerJob = 5
const defaultTimeout = time.Second * 5

func NewRunner(dir string, seedImg []byte) *Runner {
	ch := make(chan struct{}, parallelImagesPerJob+1)
	for i := 0; i <= parallelImagesPerJob; i++ {
		ch <- struct{}{}
	}

	return &Runner{tickets: ch,
		cli:     &http.Client{Timeout: defaultTimeout},
		dir:     dir,
		seedImg: seedImg,
	}
}

var randSrc = rand.NewSource(time.Now().UnixNano())

func (r *Runner) Wait() {
	r.wg.Wait()
}

func (r *Runner) Download(ctx context.Context, url string) (string, error) {
	f, err := ioutil.TempFile(r.dir, "image-*.png")
	if err != nil {
		return "", err
	}

	src := bytes.NewReader(r.seedImg)
	_, err = io.Copy(f, src)
	if err != nil {
		f.Close()
		return "", errors.Wrap(err, "writing seed image")
	}

	err = f.Sync()
	if err != nil {
		f.Close()
		return "", errors.Wrap(err, "syncing seed image")
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		f.Close()
		return "", errors.Wrap(err, "seeking from seed image to start")
	}

	<-r.tickets
	r.wg.Add(1)

	go func() {
		defer func() {
			r.tickets <- struct{}{}
		}()
		defer f.Close()
		defer r.wg.Done()

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logrus.Error(errors.Wrapf(err, "cannot create GET request to '%v'", url))
			return
		}
		req = req.WithContext(ctx)
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; SendReadable/1.0; +https://sendreadable.utrack.dev)")

		resp, err := r.cli.Do(req)
		if err != nil {
			logrus.Error(errors.Wrap(err, "downloading image"))
			return
		}
		defer resp.Body.Close()

		buf := bytes.NewBuffer(nil)
		err = convertImageToPng(resp.Body, buf)
		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot convert image to PNG"))
		}
		src := bytes.NewReader(buf.Bytes())

		_, err = io.Copy(f, src)
		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot dump image to file"))
		}
		f.Sync()
	}()
	return f.Name(), nil
}
