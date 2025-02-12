package config

import (
	"testing"
	"fmt"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("./")
	fmt.Println(config)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
}