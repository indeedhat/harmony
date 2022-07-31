package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/indeedhat/harmony/internal/app"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	"github.com/indeedhat/harmony/internal/logger"
)

func main() {
	verbose := flag.Bool("v", false, "print logs to screen rather than log file")
	flag.Parse()
	flag.Usage = usage

	if !*verbose {
		logger.UseConfigFile()
	}

	conf := config.Load()
	if conf == nil {
		return
	}

	ctx := common.NewContext(conf)

	app, err := app.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(app.Run())
}

func usage() {
	fmt.Print(`Harmony HID
Share your mouse and keyboard over the network

Usage: 
    ./harmony-hid

Options:
`)

	flag.PrintDefaults()
}
