package esapi

import (
	"encoding/base64"
	"fmt"
)

type Authorization interface {
	AuthorizationHeader() string
}

type BasicAuth struct {
	Username string
	Password string
}

func (a BasicAuth) AuthorizationHeader() string {
	creds := fmt.Sprintf("%s:%s", a.Username, a.Password)
	b64Creds := base64.StdEncoding.EncodeToString([]byte(creds))
	return fmt.Sprintf("Basic %s", b64Creds)
}
