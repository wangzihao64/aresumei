package main

import (
	"aresumei/config"
	"aresumei/routes"
)

func main() {
	config.Init()
	r := routes.NewRouter()
	r.Run(config.HttpPort)
}
