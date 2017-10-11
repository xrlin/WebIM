package main

import (
	"github.com/xrlin/WebIM/server/routes"
)

func main() {
	routes.RouterEngine().Run("127.0.0.1:8080")
}
