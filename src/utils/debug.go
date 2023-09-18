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
	t.Logf("\n\n>>>>> %T: %s <<<<<\n\n", data, string(dataJSON))
}

func PrintWD() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(">>> error reading the path")
	}
	fmt.Println(">>>", dir)
}

func PrintColor(color string, values ...interface{}) {
	switch color {
	case "black":
		fmt.Printf("\x1b[30m")
	case "red":
		fmt.Printf("\x1b[31m")
	case "green":
		fmt.Printf("\x1b[32m")
	case "blue":
		fmt.Printf("\x1b[34m")
	case "yellow":
		fmt.Printf("\x1b[33m")
	case "cyan":
		fmt.Printf("\x1b[36m")
	case "magenta":
		fmt.Printf("\x1b[35m")
	case "gray":
		fmt.Printf("\x1b[90m")
	case "white":
		fmt.Printf("\x1b[37m")
	default:
		fmt.Println("Invalid color:", color)
		return
	}

	fmt.Print(values...)

	fmt.Println("\x1b[0m")
}

func PrintBlueTesting(t *testing.T, values ...interface{}) {
	fmt.Printf("\n\n\x1b[34m >>> ") // blue
	fmt.Print(values...)
	fmt.Printf("\x1b[0m\n\n")
}
