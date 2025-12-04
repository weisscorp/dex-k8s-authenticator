package main

import (
	"os"
	"reflect"
	"testing"
)

func TestSubstituteEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "no variables",
			input:    "simple text",
			envVars:  map[string]string{},
			expected: "simple text",
		},
		{
			name:     "single variable",
			input:    "value is ${VAR1}",
			envVars:  map[string]string{"VAR1": "test"},
			expected: "value is test",
		},
		{
			name:     "multiple variables",
			input:    "${VAR1} and ${VAR2}",
			envVars:  map[string]string{"VAR1": "first", "VAR2": "second"},
			expected: "first and second",
		},
		{
			name:     "variable with dashes and underscores",
			input:    "${VAR-NAME_123}",
			envVars:  map[string]string{"VAR-NAME_123": "value"},
			expected: "value",
		},
		{
			name:     "variable not set",
			input:    "value is ${MISSING}",
			envVars:  map[string]string{},
			expected: "value is ",
		},
		{
			name:     "nested variables",
			input:    "${VAR1}${VAR2}",
			envVars:  map[string]string{"VAR1": "a", "VAR2": "b"},
			expected: "ab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			result := substituteEnvVars(tt.input)
			if result != tt.expected {
				t.Errorf("substituteEnvVars() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSubstituteEnvVarsRecursive(t *testing.T) {
	// Set test environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	original := Config{
		Listen:          "http://${TEST_VAR}:5555",
		Web_Path_Prefix: "/",
		Clusters: []Cluster{
			{
				Name:        "${TEST_VAR}-cluster",
				Issuer:      "https://${TEST_VAR}.example.com",
				Client_ID:   "client-${TEST_VAR}",
				Client_Secret: "secret-${TEST_VAR}",
			},
		},
	}

	copy := reflect.New(reflect.TypeOf(original)).Elem()
	substituteEnvVarsRecursive(copy, reflect.ValueOf(original))

	result := copy.Interface().(Config)

	if result.Listen != "http://test_value:5555" {
		t.Errorf("Listen = %v, want http://test_value:5555", result.Listen)
	}
	if result.Clusters[0].Name != "test_value-cluster" {
		t.Errorf("Cluster Name = %v, want test_value-cluster", result.Clusters[0].Name)
	}
	if result.Clusters[0].Issuer != "https://test_value.example.com" {
		t.Errorf("Issuer = %v, want https://test_value.example.com", result.Clusters[0].Issuer)
	}
}
