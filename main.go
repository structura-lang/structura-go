// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"structura/interpreter"
)

func main() {
	programFile, err := os.Open("program.json")
	if err != nil {
		log.Fatalf("failed to open program file: %v", err)
	}
	defer programFile.Close()

	dataFile, err := os.Open("data.json")
	if err != nil {
		log.Fatalf("failed to open data file: %v", err)
	}
	defer dataFile.Close()

	programBytes, err := io.ReadAll(programFile)
	if err != nil {
		log.Fatalf("failed to read program file: %v", err)
	}

	dataBytes, err := io.ReadAll(dataFile)
	if err != nil {
		log.Fatalf("failed to read data file: %v", err)
	}

	var p map[string]any
	if err := json.Unmarshal(programBytes, &p); err != nil {
		log.Fatalf("failed to unmarshal into Function: %v", err)
	}

	var d map[string]any
	if err := json.Unmarshal(dataBytes, &d); err != nil {
		log.Fatalf("failed to unmarshal into map: %v", err)
	}

	result, err := interpreter.Execute(p, d, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("FUNC[main]: %w", err))
	}

	fmt.Println(result)
}
