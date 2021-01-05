package converter

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type pullRsp struct {
	b    []byte
	lang []string
}

func pull(ctx context.Context, timeout time.Duration, url string) (*pullRsp, error) {

	// Make sure URL is valid
	_, err := nurl.ParseRequestURI(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Fetch page from URL
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating GET request")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; SendReadable/1.0; +https://sendreadable.utrack.dev)")
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	// Make sure content type is HTML
	cp := resp.Header.Get("Content-Type")
	if !strings.Contains(cp, "text/html") {
		resp.Body.Close()
		return nil, fmt.Errorf("URL is not a HTML document")
	}

	lang := resp.Header.Get("Content-Language")

	b, err := ioutil.ReadAll(resp.Body)

	return &pullRsp{b: b, lang: strings.Split(lang, ",")}, err
}
