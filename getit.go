package main

import (
	"getit/utils"
)

type ContextObjects utils.ContextObjects

type OpenUrl struct {
	url string
}

func (url OpenUrl) toContextObj() (contextObj ContextObjects) {
	return
}

func postToSfx(ctxObj ContextObjects) (response string) {
	return
}

func (ctxObj ContextObjects) toJson() (json string) {
	return
}

func main() {

}
