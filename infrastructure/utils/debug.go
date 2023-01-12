package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(data any) {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("--- error formating %T ---", data)
	}
	fmt.Printf(">>>>> %T: %s <<<<<\n", data, string(dataJSON))
}
