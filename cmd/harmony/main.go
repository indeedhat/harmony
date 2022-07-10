package main

import (
	"log"

	"github.com/indeedhat/kvm/internal/screen"
	"github.com/indeedhat/kvm/internal/ui/router"
)

func main() {
	displays, err := screen.GetDisplayBounds()
	if err != nil {
		log.Fatal(err)
	}

	r := router.New(displays)
	r.Run()
}
