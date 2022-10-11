package main

import (
	"HolyCrusade/internal/service_one"
	"HolyCrusade/pkg/core"
)

func main() {
	var a core.Application
	a.Init("config/service_one.yml").WithKafka()

	s := service_one.ServiceOne{App: a}

	s.SendToQueue()
}
