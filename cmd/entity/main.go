package main

import (
	"HolyCrusade/pkg/core"
	"fmt"
)

func main() {
	var a core.Application
	a.Init("config/service_one.yml").WithKafka()

	fmt.Println(a.Config.Version)
}
