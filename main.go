package main

//go:generate ./sqlboiler --wipe psql --add-global-variants

import (
	utils "github.com/ystv/web-api/routes"
)

func main() {
	e := utils.InitRoutes()
	e.Logger.Fatal(e.Start(":8080"))
}
