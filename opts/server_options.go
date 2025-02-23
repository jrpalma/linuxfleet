package opts

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerOptions struct {
	DatabaseCluster []string `yaml:"database_cluster"`
}

// Marshal the Options struct to YAML format
func (o *ServerOptions) Marshal() ([]byte, error) {
	return yaml.Marshal(o)
}

// Unmarshal a byte slice into the Options struct from YAML format
func (o *ServerOptions) Unmarshal(data []byte) error {
	return yaml.Unmarshal(data, o)
}

// WriteOptions writes the Options object to a YAML file
func (o *ServerOptions) WriteOptions(filename string) error {
	yamlData, err := o.Marshal()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, yamlData, 0644)
}

// ReadOptions reads the Options object from a YAML file
func (o *ServerOptions) ReadOptions(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = o.Unmarshal(data)
	if err != nil {
		return err
	}
	return nil
}
