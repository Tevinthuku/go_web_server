package main

import (
	"log"
	"net"
	"os"
	"web_server/webserver"
)

func main() {
	listner, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listner.Close()
	rootDirEnv := os.Getenv("ROOT_DIR")
	if rootDirEnv == "" {
		rootDirEnv = "./www"
	}
	wb := webserver.NewWebServer(rootDirEnv)
	defer wb.Close()

	RegisterRoutes(wb)

	if err := wb.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
