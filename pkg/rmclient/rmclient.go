package rmclient

import (
	"context"
	"io"
	"time"

	"github.com/juruen/rmapi/api"
	"github.com/juruen/rmapi/log"
	"github.com/juruen/rmapi/model"
	"github.com/juruen/rmapi/transport"
	"github.com/pkg/errors"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

type AuthResponse struct {
	Tokens Tokens
}

const (
	deviceDesc = "desktop-linux"
	//deviceUuid = "c9136f4f-4bb4-4860-9a80-61246ed245b3"

	authHost         = "https://my.remarkable.com"
	docHost          = "https://document-storage-production-dot-remarkable-production.appspot.com"
	newTokenDevice   = authHost + "/token/json/2/device/new"
	newUserDevice    = authHost + "/token/json/2/user/new"
	urlUploadRequest = docHost + "/document-storage/json/2/upload/request"
	urlUpdateStatus  = docHost + "/document-storage/json/2/upload/update-status"
)

func init() {
	transport.RmapiUserAGent = "SendReadable/1.0; +https://sendreadable.utrack.dev"
	log.InitLog()
}

func (c *Client) Auth(ctx context.Context, code string) (*AuthResponse, error) {

	devToken, err := getDeviceToken(code)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create new device token")
	}

	userToken, err := getUserToken(devToken)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create new user token")
	}

	return &AuthResponse{
		Tokens: Tokens{
			Device:          devToken,
			User:            userToken,
			UserRefreshedAt: time.Now(),
		},
	}, nil
}

type Tokens struct {
	// Device token used to create and refresh user tokens.
	Device string
	// User tokens used for all operations.
	User            string
	UserRefreshedAt time.Time
}

func (c *Client) Upload(ctx context.Context, name string, r io.Reader, tok *Tokens) error {
	if tok.UserRefreshedAt.Before(time.Now().Add(-time.Hour)) {
		if tok.Device == "" {
			return errors.New("rM auth method had been changed to be more stable, please relogin. Sorry for the inconvenience :(")
		}
		userTok, err := getUserToken(tok.Device)
		if err != nil {
			return errors.Wrap(err, "cannot refresh user token")
		}
		tok.User = userTok
	}
	cli := transport.CreateHttpClientCtx(model.AuthTokens{UserToken: tok.User})

	tree, err := api.DocumentsFileTree(&cli)
	if err != nil {
		return errors.Wrap(err, "cannot get root FS structure")
	}
	var dirID string
	dir, err := tree.NodeByPath("/SendReadable", tree.Root())
	if err != nil {
		dirID, err = createDir(cli, tree.Root().Id(), "SendReadable")
		if err != nil {
			return errors.Wrap(err, "creating directory")
		}
	} else {
		dirID = dir.Id()
	}

	return uploadDoc(cli, dirID, name, r)
}

func createDir(cli transport.HttpClientCtx, parentId string, name string) (string, error) {
	dirUploadRsp, err := uploadRequest(cli, "", model.DirectoryType)
	if err != nil {
		return "", errors.Wrap(err, "cannot create request")
	}

	dr, err := createDirectoryZip(dirUploadRsp.ID)
	if err != nil {
		return "", errors.Wrap(err, "cannot create asset zip")
	}

	err = cli.PutStream(transport.UserBearer, dirUploadRsp.BlobURLPut, dr)
	if err != nil {
		return "", errors.Wrap(err, "cannot upload model's data")
	}

	metaDoc := model.CreateUploadDocumentMeta(dirUploadRsp.ID, model.DirectoryType, parentId, name)

	err = cli.Put(transport.UserBearer, urlUpdateStatus, metaDoc, nil)

	if err != nil {
		return "", errors.Wrap(err, "cannot move directory entry")
	}

	doc := metaDoc.ToDocument()

	return doc.ID, err
}

func uploadDoc(cli transport.HttpClientCtx, dirID string, name string, r io.Reader) error {
	rsp, err := uploadRequest(cli, "", model.DocumentType)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}
	if !rsp.Success {
		return errors.New("upload request did not succeed")
	}

	zip, err := createFileZip(rsp.ID, r)
	if err != nil {
		return errors.Wrap(err, "cannot create content zipfile")
	}

	err = cli.PutStream(transport.UserBearer, rsp.BlobURLPut, zip)

	if err != nil {
		return errors.Wrap(err, "failed to upload zip document")
	}

	metaDoc := model.CreateUploadDocumentMeta(rsp.ID, model.DocumentType, dirID, name)

	err = cli.Put(transport.UserBearer, urlUpdateStatus, metaDoc, nil)

	return errors.Wrap(err, "failed to move entry")
}

func uploadRequest(cli transport.HttpClientCtx, id string, entryType string) (model.UploadDocumentResponse, error) {
	uploadReq := model.CreateUploadDocumentRequest(id, entryType)
	uploadRsp := make([]model.UploadDocumentResponse, 0)

	err := cli.Put(transport.UserBearer, urlUploadRequest, uploadReq, &uploadRsp)

	if err != nil {
		return model.UploadDocumentResponse{}, errors.Wrap(err, "cannot send upload request")
	}

	return uploadRsp[0], nil
}
