package config

import (
	"os"
	"log"
)

var BaasFabricEngine string

func init() {
	BaasFabricEngine = os.Getenv("BaasFabricEngine")
	log.Println("BaasFabricEngine:", BaasFabricEngine)
	if BaasFabricEngine == "" {
		log.Fatal("no env BaasFabricEngine")
	}
}
