package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func PrettyPrint(data any) {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("--- error formating %T ---", data)
	}
	fmt.Printf(">>>>> %T: %s <<<<<\n", data, string(dataJSON))
}

func PrettyPrintTesting(t *testing.T, data any) {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Logf("--- error formating %T ---", data)
	}
	t.Logf(">>>>> %T: %s <<<<<\n", data, string(dataJSON))
}

func PrintWD() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(">>> error reading the path")
	}
	fmt.Println(">>>", dir)
}
