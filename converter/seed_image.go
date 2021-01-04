package converter

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
)

var seedImg []byte

func init() {
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(seedImgBase64)))
	var err error
	seedImg, err = ioutil.ReadAll(dec)
	if err != nil {
		panic(err)
	}
}

var seedImgBase64 = `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`
