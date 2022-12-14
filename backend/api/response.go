package api

import (
	"ariadne/sfx"
)

type Response struct {
	Errors  []string                    `json:"errors"`
	Records sfx.MultiObjXMLResponseBody `json:"records"`
}
