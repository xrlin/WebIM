package main

import (
	"github.com/xrlin/WebIM/server/routes"
)

func main() {
	routes.RouterEngin().Run(":8080")
}