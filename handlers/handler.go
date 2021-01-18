package handlers

import (
	"crypto/rsa"

	"github.com/utrack/sendreadable/converter"
	"github.com/utrack/sendreadable/pkg/rmclient"
)

type Handler struct {
	svc           *converter.Service
	chopDirPrefix string

	rm     *rmclient.Client
	jwtKey *rsa.PrivateKey

	secure bool
}

func New(svc *converter.Service,
	dirPrefix string,
	rm *rmclient.Client,
	jkey *rsa.PrivateKey,
	secure bool,
) Handler {
	return Handler{
		svc:           svc,
		chopDirPrefix: dirPrefix,
		rm:            rm,
		jwtKey:        jkey,
		secure:        secure,
	}
}
