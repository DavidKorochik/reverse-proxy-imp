package main

import (
	"davidk/reverse-proxy-imp/pkg/connection"
	"davidk/reverse-proxy-imp/pkg/proxy"
)

func main() {
	connections, err := connection.GetConnectionsBatch()
	if err != nil {
		panic(err)
	}

	proxy.NewProxy(connections, 8800).Run()
}
