package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"regexp"
)

type Service struct {
	Name           string            `json:"name"`
	Args           []string          `json:"args"`
	FilesMapped    map[string]string `json:"files_mapped"`
	VolumesMapped  map[string]string `json:"volumes_mapped"`	
	OtherContent   map[string]any    `json:"-"` // catch-all for other keys
}

type Input struct {
	Services []Service `json:"services"`
}

func ProcessFile(filepath string, outputPath string) {
	// Load the original JSON
	inputBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	// Decode JSON into a map to preserve other keys
	var raw map[string]any
	if err := json.Unmarshal(inputBytes, &raw); err != nil {
		panic(err)
	}

	servicesData := raw["services"].([]any)
	mountPath := "/artifacts"
	portRegexp := regexp.MustCompile(`{{\s*(Port|PortUDP)\s+"[^"]+"\s+(\d+)\s*}}`)
	for i, s := range servicesData {
		serviceMap := s.(map[string]any)
		args := toStringSlice(serviceMap["args"])
		filesMapped := toStringMap(serviceMap["files_mapped"])
		volumesMapped := toStringMap(serviceMap["volumes_mapped"])

		// Replace args that match any key in files_mapped
		for j, arg := range args {
			if newVal, found := filesMapped[arg]; found {
				args[j] = mountPath + "/" + newVal
			}
		}

		// Replace {{Port "..." 12345}} and {{PortUDP "..." 12345}} with "12345"
		for k, argK := range args {
			args[k] = portRegexp.ReplaceAllString(argK, "$2")
		}

		for key := range volumesMapped {
			for j, arg := range args {
				if strings.HasPrefix(arg, key) {
					args[j] = mountPath + arg
				}
			}
        }

		// Add volumes configure for each container
		serviceMap["volumes"] = []map[string]string{
			{
				"name":      "artifacts",
				"mountPath": mountPath,
			},
		}

		// Put modified args back
		serviceMap["args"] = args
		servicesData[i] = serviceMap
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		panic(fmt.Errorf("failed to create output directory: %w", err))
	}

	// Write modified data back
	raw["services"] = servicesData
	outputBytes, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(outputPath+"/processed.json", outputBytes, 0644); err != nil {
		panic(err)
	}

	fmt.Println("✅ JSON processing complete. See output.json.")
}

func toStringSlice(i any) []string {
	raw := i.([]any)
	res := make([]string, len(raw))
	for i, v := range raw {
		res[i] = v.(string)
	}
	return res
}

func toStringMap(i any) map[string]string {
	raw, ok := i.(map[string]any)
	if !ok || raw == nil {
		return map[string]string{}
	}

	res := make(map[string]string)
	for k, v := range raw {
		strVal, ok := v.(string)
		if !ok {
			continue // or log/handle as needed
		}
		res[k] = strVal
	}
	return res
}
