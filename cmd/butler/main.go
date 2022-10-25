package main

import (
	"HolyCrusade/pkg/core"
)

func main() {
	var a core.Application
	a.Init("config/butler.yml").
		// WithKafka().
		WithTelegramBot()
}
