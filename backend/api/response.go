package api

import (
	"resolve/sfx"
)

type response struct {
	Errors  []string                    `json:"errors""`
	Records sfx.MultiObjXMLResponseBody `json:"records"`
}
