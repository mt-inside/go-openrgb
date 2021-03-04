package main

import (
	"fmt"
	"os"

	"github.com/mt-inside/go-openrgb/pkg/model"
	"github.com/mt-inside/logging"
)

func main() {
	log := logging.GetLogger(false)

	m, err := model.NewModel(log, "localhost:6742", "go-openrgb info example")
	if err != nil {
		log.Error(err, "Couldn't synchronise devices and colors from server")
		os.Exit(1)
	}

	fmt.Println(m)
}
