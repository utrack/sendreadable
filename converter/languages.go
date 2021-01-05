package converter

import (
	"strings"

	"github.com/ryboe/q"
)

var codeToBabelName = map[string]string{
	"en": "english",
	"ru": "russian",
	"es": "spanish",
	"it": "italian",
	"de": "ngerman",
	"fr": "french",
	"nl": "dutch",
}

func langsToArray(htmlLang string, ctypes []string) []string {
	var ret []string

	if v, ok := codeToBabelName[htmlLang]; ok {
		ret = append(ret, v)
	}
	for _, c := range ctypes {
		c = strings.Split(c, "-")[0]
		if v, ok := codeToBabelName[c]; ok {
			ret = append(ret, v)
		}
	}
	q.Q("langs", htmlLang, ctypes, ret)
	return ret
}
