package config

import (
	"fmt"

	"testing"
)

func TestLoad(t *testing.T) {
	config := Load("config.yaml")

	fmt.Print(config)
}
