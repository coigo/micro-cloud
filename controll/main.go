package main

import (
	"fmt"
	"github.com/coigo/micro-cloud/commandservice"
)

func main () {
	dockerId := commandservice.UpCommand()
	fmt.Printf("Container %v criado.\n", dockerId)
	commandservice.DownCommand(dockerId)
}