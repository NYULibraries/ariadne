package api

import (
	"ariadne/sfx"
)

type response struct {
	Errors  []string                    `json:"errors""`
	Records sfx.MultiObjXMLResponseBody `json:"records"`
}
