package main

import (
	"flag"
	"log"

	"github.com/indeedhat/harmony/internal/app"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	"github.com/indeedhat/harmony/internal/logger"
)

func main() {
	verbose := flag.Bool("v", false, "print logs to screen rather than log file")
	flag.Parse()

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
