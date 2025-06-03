package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

func ProcessFileForURL(filepath string, outputPath string) {

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

	
    ReplaceServiceArgs(input.Services)

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

	
	// services := []Service{
	// 	{
	// 		Name: "el",
	// 		Args: []string{
	// 			"--authrpc.port", "{{Service \"el\" \"authrpc\" \"http\" \"\"}}",
	// 		},
	// 		Ports: []Port{
	// 			{Name: "authrpc", Port: 8551, Protocol: "tcp"},
	// 		},
	// 	},
	// 	{
	// 		Name: "mev-boost",
	// 		Args: []string{
	// 			"--mev-boost", "{{Service \"mev-boost\" \"http\" \"http\" \"\"}}",
	// 		},
	// 		Ports: []Port{
	// 			{Name: "http", Port: 5555, Protocol: "tcp"},
	// 		},
	// 	},
	// }

	// ReplaceServiceArgs(services)

	// // Print results
	// for _, svc := range services {
	// 	fmt.Printf("Service %s Args:\n", svc.Name)
	// 	for _, arg := range svc.Args {
	// 		fmt.Println("  ", arg)
	// 	}
	// }
}

// ReplaceServiceArgs replaces all {{Service ...}} patterns with http://localhost:<port>
func ReplaceServiceArgs(services []Service) {
	// Regex pattern to match: {{Service "service" "portname" "scheme" ""}}
	re := regexp.MustCompile(`{{\s*Service\s+"([^"]+)"\s+"([^"]+)"\s+"([^"]+)"\s+"[^"]*"\s*}}`)

	for i := range services {
		for j, arg := range services[i].Args {
			services[i].Args[j] = re.ReplaceAllStringFunc(arg, func(match string) string {
				parts := re.FindStringSubmatch(match)
				if len(parts) != 4 {
					return match // Leave unchanged if pattern doesn't match
				}
				serviceName := parts[1]
				portName := parts[2]
				scheme := parts[3]

				// Lookup the service
				for _, svc := range services {
					if svc.Name != serviceName {
						continue
					}
					// Lookup the port
					for _, p := range svc.Ports {
						if p.Name == portName {
							// Return formatted address
							return fmt.Sprintf("%s://localhost:%d", scheme, p.Port)
						}
					}
				}
				// If not found, leave as-is
				return match
			})
		}
	}
}
