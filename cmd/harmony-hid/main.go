package main

import (
	"log"

	"github.com/indeedhat/harmony/internal/app"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
)

func main() {

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
