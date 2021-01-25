package rmclient

import (
	"github.com/gofrs/uuid"
	"github.com/juruen/rmapi/model"
	"github.com/juruen/rmapi/transport"
	"github.com/pkg/errors"
)

func getDeviceToken(code string) (string, error) {
	cli := transport.CreateHttpClientCtx(model.AuthTokens{})

	uuid, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(err, "cannot generate new UUID v4")
	}

	req := model.DeviceTokenRequest{code, deviceDesc, uuid.String()}

	resp := transport.BodyString{}
	err = cli.Post(transport.EmptyBearer, newTokenDevice, req, &resp)
	return resp.Content, err
}

func getUserToken(devTok string) (string, error) {
	cli := transport.CreateHttpClientCtx(model.AuthTokens{DeviceToken: devTok})
	resp := transport.BodyString{}
	err := cli.Post(transport.DeviceBearer, newUserDevice, nil, &resp)
	return resp.Content, err
}
