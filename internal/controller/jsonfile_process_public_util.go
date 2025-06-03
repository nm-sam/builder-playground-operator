package controller
import (
	"encoding/json"
)

type Input struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name           string                  `json:"name"`
	Args           []string                `json:"args"`
	FilesMapped    map[string]string       `json:"files_mapped"`
	VolumesMapped  map[string]string       `json:"volumes_mapped"`	
	Ports          []Port                  `json:"ports"`
	OtherContent   map[string]any          `json:"-"` // catch-all for other keys
}

type Port struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"Protocol"`
	HostPort int    `json:"HostPort"`
}

func (s *Service) UnmarshalJSON(data []byte) error {
	type Alias Service // Create an alias to avoid recursion
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Capture unknown fields
	var rawMap map[string]any
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Remove known fields
	delete(rawMap, "name")
	delete(rawMap, "args")
	delete(rawMap, "files_mapped")
	delete(rawMap, "volumes_mapped")
	delete(rawMap, "ports")

	s.OtherContent = rawMap
	return nil
}

func (s Service) MarshalJSON() ([]byte, error) {
	// Reconstruct the full map with known + unknown fields
	type Alias Service
	aux := make(map[string]any)

	// Marshal known fields
	known, err := json.Marshal(Alias(s))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(known, &aux); err != nil {
		return nil, err
	}

	// Merge with OtherContent
	for k, v := range s.OtherContent {
		aux[k] = v
	}

	return json.Marshal(aux)
}