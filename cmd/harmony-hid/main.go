package main

import (
	"log"

	"github.com/indeedhat/harmony/internal/app"
	"github.com/indeedhat/harmony/internal/common"
)

func main() {
	ctx := common.NewContext()

	app, err := app.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(app.Run())
}
