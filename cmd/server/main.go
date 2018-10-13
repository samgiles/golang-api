package main

import (
	"github.com/samgiles/health"
)

func main() {
	controller := health.NewHealthCheckController()
	defer controller.Stop()
}
