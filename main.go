package main

import (
	"log"
	"os"
	"web_server/webserver"
)

func main() {
	rootDirEnv := os.Getenv("ROOT_DIR")
	if rootDirEnv == "" {
		rootDirEnv = "./www"
	}
	wb := webserver.NewWebServer(rootDirEnv)
	defer wb.Close()

	RegisterRoutes(wb)
	err := wb.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
