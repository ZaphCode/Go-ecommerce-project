package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func PrettyPrint(data any) {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("--- error formating %T ---", data)
	}
	fmt.Printf(">>>>> %T: %s <<<<<\n", data, string(dataJSON))
}

func PrintWD() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(">>> error reading the path")
	}
	log.Println(">>>", dir)
}
