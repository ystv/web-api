package main

//go:generate ./sqlboiler --wipe psql --add-global-variants

import (
	"github.com/ystv/web-api/routes"
)

func main() {
	e := routes.Init()
	e.Logger.Fatal(e.Start(":8080"))
}
