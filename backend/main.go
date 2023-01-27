package main

import (
	"ariadne/cmd"
	"ariadne/util"
)

func main() {
	util.InitSentry()
	defer util.RecoverWithSentry()
	cmd.Execute()
}
