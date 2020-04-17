package main

//go:generate ./sqlboiler --wipe psql

import (
	"github.com/ystv/web-api/utils"
)

func main() {
	e := utils.InitRoutes()
	e.Logger.Fatal(e.Start(":8080"))
}
