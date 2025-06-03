package controller

import (
	"encoding/json"
	"fmt"
	"os"
)

func ProcessFileForPorts(filepath string, outputPath string) {
	// Read the input JSON file
	inputBytes, err := os.ReadFile(filepath + "/" + "processed.json")
	if err != nil {
		panic(fmt.Errorf("error reading file: %w", err))
	}

	// Parse JSON into Go struct
	var input Input
	if err := json.Unmarshal(inputBytes, &input); err != nil {
		panic(fmt.Errorf("error unmarshaling: %w", err))
	}

	// Normalize ports
	for si, service := range input.Services {
		nameCount := make(map[string]int)

		// Count occurrences of each port name
		for _, port := range service.Ports {
			nameCount[port.Name]++
		}

		// Modify duplicate names
		for pi, port := range service.Ports {
			if nameCount[port.Name] > 1 {
				normalizedName := fmt.Sprintf("%s-%s", port.Name, port.Protocol)
				input.Services[si].Ports[pi].Name = normalizedName
			}
		}
	}

	// Output normalized JSON if desired
	outputBytes, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		panic(fmt.Errorf("error marshaling: %w", err))
	}

	// Write to file
	if err := os.WriteFile(outputPath  + "/" + "processed.json", outputBytes, 0644); err != nil {
		panic(fmt.Errorf("error writing output: %w", err))
	}

	fmt.Println("✅ File processed and ports normalized.")
}
