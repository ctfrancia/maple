package main

import ()

var (
	build  = "dev" // todo: change to commit hash
	routes = "all" // go build -ldflags "-X main.routes=crud"
)

func main() {
	println("build:", build)
	println("routes:", routes)
}
