package main

import (
	"log"
	"os"
	"pc28/xmd"
)

func main() {
	log.Printf("当前版本 %q >>> \n", "2023.10.18-rc3")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s \n", err.Error())
	}

	cache, err := xmd.NewCache(dir)
	if err != nil {
		log.Fatalf("xmd.NewCache() fail : %s\n", err.Error())
	}

	xmd.Run(cache)
}
