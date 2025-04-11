// Copyright 2024 Humanitec
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package convert

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	scoretypes "github.com/score-spec/score-go/types"
	"github.com/stretchr/testify/assert"
)

// func TestConvertToBicep(t *testing.T) {
// 	// Helper function to assert string contains
// 	assertContains := func(t *testing.T, result, expected, message string) {
// 		t.Helper()
// 		if !strings.Contains(result, expected) {
// 			t.Errorf("%s, got: %s", message, result)
// 		}
// 	}

// 	// Create test cases based on the inferred structure of scoretypes.Workload
// 	tests := []struct {
// 		name         string
// 		workload     scoretypes.Workload
// 		workloadName string
// 		wantErr      bool
// 		checkFunc    func(t *testing.T, result string)
// 	}{
// 		{
// 			name: "basic workload without service",
// 			workload: scoretypes.Workload{
// 				Containers: map[string]scoretypes.Container{
// 					"main": {
// 						Image: "nginx:latest",
// 					},
// 				},
// 			},
// 			workloadName: "test-workload",
// 			wantErr:      false,
// 			checkFunc: func(t *testing.T, result string) {
// 				// Check for basic structure
// 				assertContains(t, result, "targetScope = 'subscription'", "Expected Bicep header with targetScope")
// 				assertContains(t, result, "resource resourceGroup", "Expected resource group definition")
// 				assertContains(t, result, "module containerAppEnvironment", "Expected container app environment module")
// 				assertContains(t, result, "resource containerApp", "Expected container app resource")
// 				assertContains(t, result, "image: 'nginx:latest'", "Expected container image to be set")

// 				// Check for default scale settings
// 				assertContains(t, result, "minReplicas: 1", "Expected default minReplicas")
// 				assertContains(t, result, "maxReplicas: 10", "Expected default maxReplicas")
// 			},
// 		},
// 		{
// 			name: "workload with service",
// 			workload: scoretypes.Workload{
// 				Containers: map[string]scoretypes.Container{
// 					"main": {
// 						Image: "nginx:latest",
// 					},
// 				},
// 				Service: &scoretypes.WorkloadService{
// 					Ports: map[string]scoretypes.ServicePort{
// 						"http": {
// 							Port: 80,
// 						},
// 					},
// 				},
// 			},
// 			workloadName: "test-workload-with-service",
// 			wantErr:      false,
// 			checkFunc: func(t *testing.T, result string) {
// 				assertContains(t, result, "ingress: {", "Expected ingress configuration")
// 				assertContains(t, result, "external: true", "Expected external ingress")
// 				assertContains(t, result, "targetPort: 80", "Expected target port in ingress")

// 			},
// 		},
// 		{
// 			name: "workload with container variables",
// 			workload: scoretypes.Workload{
// 				Containers: map[string]scoretypes.Container{
// 					"main": {
// 						Image: "nginx:latest",
// 						Variables: map[string]string{
// 							"ENV_VAR":     "value",
// 							"DB_PASSWORD": "secret-password",
// 						},
// 					},
// 				},
// 			},
// 			workloadName: "test-workload-with-vars",
// 			wantErr:      false,
// 			checkFunc: func(t *testing.T, result string) {
// 				assertContains(t, result, "env: [", "Expected environment variables")
// 				assertContains(t, result, "name: 'ENV_VAR'", "Expected ENV_VAR environment variable")
// 				assertContains(t, result, "value: 'value'", "Expected ENV_VAR value")

// 				// Check for secrets section - DB_PASSWORD should be detected as a secret
// 				assertContains(t, result, "secrets: [", "Expected secrets section")
// 				assertContains(t, result, "name: 'db_password-test-workload-with-vars'", "Expected DB_PASSWORD to be treated as a secret")
// 				assertContains(t, result, "value: 'secret-password'", "Expected secret value to be set")

// 				// Check that the env var references the secret
// 				assertContains(t, result, "name: 'DB_PASSWORD'", "Expected DB_PASSWORD environment variable")
// 				assertContains(t, result, "secretRef: 'db_password-test-workload-with-vars'", "Expected DB_PASSWORD to reference secret")
// 			},
// 		},
// 		{
// 			name: "workload with container resources",
// 			workload: scoretypes.Workload{
// 				Containers: map[string]scoretypes.Container{
// 					"main": {
// 						Image: "nginx:latest",
// 						Resources: &scoretypes.ContainerResources{
// 							Limits: &scoretypes.ResourcesLimits{
// 								Cpu:    stringPtr("500m"),
// 								Memory: stringPtr("512Mi"),
// 							},
// 						},
// 					},
// 				},
// 			},
// 			workloadName: "test-workload-with-resources",
// 			wantErr:      false,
// 			checkFunc: func(t *testing.T, result string) {
// 				assertContains(t, result, "resources: {", "Expected resources section")
// 				assertContains(t, result, "cpu: json('0.5')", "Expected CPU resource with correct value")
// 				assertContains(t, result, "memory: '512Mi'", "Expected memory resource with correct value")
// 			},
// 		},
// 		{
// 			name: "workload with probes",
// 			workload: scoretypes.Workload{
// 				Containers: map[string]scoretypes.Container{
// 					"main": {
// 						Image: "nginx:latest",
// 						LivenessProbe: &scoretypes.ContainerProbe{
// 							HttpGet: &scoretypes.HttpProbe{
// 								Path: "/health",
// 								Port: 8080,
// 							},
// 						},
// 						ReadinessProbe: &scoretypes.ContainerProbe{
// 							HttpGet: &scoretypes.HttpProbe{
// 								Path: "/ready",
// 								Port: 8080,
// 							},
// 						},
// 					},
// 				},
// 			},
// 			workloadName: "test-workload-with-probes",
// 			wantErr:      false,
// 			checkFunc: func(t *testing.T, result string) {
// 				assertContains(t, result, "probes: [", "Expected probes section")
// 				assertContains(t, result, "type: 'Liveness'", "Expected liveness probe")
// 				assertContains(t, result, "type: 'Readiness'", "Expected readiness probe")
// 				assertContains(t, result, "path: '/health'", "Expected health path")
// 				assertContains(t, result, "port: 8080", "Expected port in health probe")
// 				assertContains(t, result, "path: '/ready'", "Expected ready path")
// 				assertContains(t, result, "port: 8080", "Expected port in ready probe")
// 			},
// 		},
// 		{
// 			name: "workload with multiple containers",
// 			workload: scoretypes.Workload{
// 				Containers: map[string]scoretypes.Container{
// 					"main": {
// 						Image: "nginx:latest",
// 					},
// 					"sidecar": {
// 						Image: "redis:alpine",
// 						Command: []string{
// 							"redis-server",
// 						},
// 						Args: []string{
// 							"--appendonly", "yes",
// 						},
// 					},
// 				},
// 			},
// 			workloadName: "test-workload-multi-container",
// 			wantErr:      false,
// 			checkFunc: func(t *testing.T, result string) {
// 				// Check for both containers
// 				assertContains(t, result, "name: 'main'", "Expected main container")
// 				assertContains(t, result, "image: 'nginx:latest'", "Expected main container image")
// 				assertContains(t, result, "name: 'sidecar'", "Expected sidecar container")
// 				assertContains(t, result, "image: 'redis:alpine'", "Expected sidecar container image")

// 				// Check for command and args
// 				assertContains(t, result, "command: ['redis-server']", "Expected command in sidecar container")
// 				assertContains(t, result, "args: ['--appendonly', 'yes']", "Expected args in sidecar container")
// 			},
// 		},
// 		{
// 			name: "workload with service and target port",
// 			workload: scoretypes.Workload{
// 				Containers: map[string]scoretypes.Container{
// 					"main": {
// 						Image: "nginx:latest",
// 					},
// 				},
// 				Service: &scoretypes.WorkloadService{
// 					Ports: map[string]scoretypes.ServicePort{
// 						"http": {
// 							Port:       80,
// 							TargetPort: intPtr(8080),
// 						},
// 					},
// 				},
// 			},
// 			workloadName: "test-workload-target-port",
// 			wantErr:      false,
// 			checkFunc: func(t *testing.T, result string) {
// 				assertContains(t, result, "targetPort: 8080", "Expected target port to be 8080")
// 			},
// 		},
// 		{
// 			name:     "error case - empty workload",
// 			workload: scoretypes.Workload{
// 				// Empty workload with no containers
// 			},
// 			workloadName: "test-empty-workload",
// 			wantErr:      true, // Expect an error for empty workload
// 			checkFunc:    nil,  // No check function needed for error case
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := convertToBicep(tt.workload, tt.workloadName)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("convertToBicep() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if err == nil && tt.checkFunc != nil {
// 				tt.checkFunc(t, result)
// 			}
// 		})
// 	}
// }

// TestCreateContainerAppProperties tests the createContainerAppProperties function
func TestCreateContainerAppProperties(t *testing.T) {
	tests := []struct {
		name         string
		workload     scoretypes.Workload
		workloadName string
		wantErr      bool
		check        func(t *testing.T, props *ContainerAppProperties)
	}{
		{
			name: "basic container app properties",
			workload: scoretypes.Workload{
				Containers: map[string]scoretypes.Container{
					"main": {
						Image: "nginx:latest",
					},
				},
			},
			workloadName: "test-workload",
			wantErr:      false,
			check: func(t *testing.T, props *ContainerAppProperties) {
				if props == nil {
					t.Fatal("Expected non-nil ContainerAppProperties")
				}
				if len(props.Template.Containers) != 1 {
					t.Errorf("Expected 1 container, got %d", len(props.Template.Containers))
				}
				if props.Template.Containers[0].Image != "nginx:latest" {
					t.Errorf("Expected image 'nginx:latest', got '%s'", props.Template.Containers[0].Image)
				}
			},
		},
		{
			name: "container app with service",
			workload: scoretypes.Workload{
				Containers: map[string]scoretypes.Container{
					"main": {
						Image: "nginx:latest",
					},
				},
				Service: &scoretypes.WorkloadService{
					Ports: map[string]scoretypes.ServicePort{
						"http": {
							Port: 80,
						},
					},
				},
			},
			workloadName: "test-workload-with-service",
			wantErr:      false,
			check: func(t *testing.T, props *ContainerAppProperties) {
				if props.Configuration.Ingress == nil {
					t.Fatal("Expected non-nil Ingress")
				}
				if !props.Configuration.Ingress.External {
					t.Error("Expected External to be true")
				}
				if props.Configuration.Ingress.TargetPort != 80 {
					t.Errorf("Expected TargetPort 80, got %d", props.Configuration.Ingress.TargetPort)
				}
				if props.Configuration.Ingress.Transport != "auto" {
					t.Errorf("Expected Transport 'auto', got '%s'", props.Configuration.Ingress.Transport)
				}
			},
		},
		{
			name: "container app with resources",
			workload: scoretypes.Workload{
				Containers: map[string]scoretypes.Container{
					"main": {
						Image: "nginx:latest",
						Resources: &scoretypes.ContainerResources{
							Requests: &scoretypes.ResourcesLimits{
								Cpu:    stringPtr("0.5"),
								Memory: stringPtr("1Gi"),
							},
						},
					},
				},
			},
			workloadName: "test-workload-with-resources",
			wantErr:      false,
			check: func(t *testing.T, props *ContainerAppProperties) {
				if len(props.Template.Containers) == 0 {
					t.Fatal("Expected at least one container")
				}
				container := props.Template.Containers[0]
				// CPU value should be 0.5 for "500m"
				if container.Resources.CPU != 0.5 {
					t.Errorf("Expected CPU 0.5, got %v", container.Resources.CPU)
				}
				if container.Resources.Memory != "1Gi" {
					t.Errorf("Expected Memory '1Gi', got '%s'", container.Resources.Memory)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			props, err := createContainerAppProperties(tt.workload)
			if (err != nil) != tt.wantErr {
				t.Errorf("createContainerAppProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				tt.check(t, props)
			}
		})
	}
}

// TestParseCPU tests the parseCPU function
func TestParseCPU(t *testing.T) {
	tests := []struct {
		name    string
		cpu     string
		want    float64
		wantErr bool
	}{
		{
			name:    "millicores",
			cpu:     "500m",
			want:    0.5,
			wantErr: false,
		},
		{
			name:    "cores",
			cpu:     "2",
			want:    2.0,
			wantErr: false,
		},
		{
			name:    "decimal cores",
			cpu:     "0.25",
			want:    0.25,
			wantErr: false,
		},
		{
			name:    "invalid format",
			cpu:     "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCPU(tt.cpu)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCPU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseCPU() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseFloat tests the parseFloat function
func TestParseFloat(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    float64
		wantErr bool
	}{
		{
			name:    "integer",
			s:       "5",
			want:    5.0,
			wantErr: false,
		},
		{
			name:    "decimal",
			s:       "2.5",
			want:    2.5,
			wantErr: false,
		},
		{
			name:    "invalid",
			s:       "not-a-number",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloat(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGenerateBicepHeader tests the generateBicepHeader function
func TestGenerateBicepHeader(t *testing.T) {
	header := generateBicepHeader()

	assert.Contains(t, header, "// Generated by score-aca")
	assert.Contains(t, header, "// Azure Container Apps Bicep manifest")
}

// TestGenerateBicepParameters tests the generateBicepParameters function
func TestGenerateBicepParameters(t *testing.T) {
	params, err := generateBicepParameters("test-name")
	expected := `
// Parameters
param environmentName string = 'test-name-environment'
param containerAppName string = 'test-name-container-app'
param location string = resourceGroup().location

`
	assert.NoError(t, err)
	assert.Equal(t, expected, params)
}

// TestGenerateContainerAppEnvironment tests the generateContainerAppEnvironment function
func TestGenerateContainerAppEnvironment(t *testing.T) {
	env := generateContainerAppEnvironment()
	expected := `// Container App Environment
resource containerAppEnvironment 'Microsoft.App/managedEnvironments@2024-03-01' = {
  name: environmentName
  location: location
  properties: {
    appLogsConfiguration: {
      destination: 'azure-monitor'
    }
  }
}
`
	assert.Equal(t, expected, env)
}

// TestConvertContainerVariables tests the convertContainerVariables function
func TestConvertContainerVariables(t *testing.T) {
	// Create a simple substitution function for testing
	sf := func(s string) (string, error) {
		if s == "TEST_VAR" {
			return "test-value", nil
		}
		if s == "ERROR_VAR" {
			return "", fmt.Errorf("test error")
		}
		return s, nil
	}

	tests := []struct {
		name    string
		input   map[string]string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "simple variables",
			input: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			want: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			wantErr: false,
		},
		{
			name: "variables with substitution",
			input: map[string]string{
				"VAR1": "${TEST_VAR}",
				"VAR2": "prefix-${TEST_VAR}-suffix",
			},
			want: map[string]string{
				"VAR1": "test-value",
				"VAR2": "prefix-test-value-suffix",
			},
			wantErr: false,
		},
		{
			name: "variables with error",
			input: map[string]string{
				"VAR1": "${ERROR_VAR}",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertContainerVariables(tt.input, sf)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertContainerVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("convertContainerVariables() got %v, want %v", got, tt.want)
					return
				}
				for k, v := range tt.want {
					if got[k] != v {
						t.Errorf("convertContainerVariables() got[%s] = %v, want %v", k, got[k], v)
					}
				}
			}
		})
	}
}

// TestConvertContainerFiles tests the convertContainerFiles function
func TestConvertContainerFiles(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test-file.txt")
	err := os.WriteFile(testFilePath, []byte("test content with ${TEST_VAR}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a simple substitution function for testing
	sf := func(s string) (string, error) {
		if s == "TEST_VAR" {
			return "test-value", nil
		}
		if s == "ERROR_VAR" {
			return "", fmt.Errorf("test error")
		}
		return s, nil
	}

	scoreFile := filepath.Join(tmpDir, "score.yaml")

	tests := []struct {
		name      string
		input     []scoretypes.ContainerFilesElem
		scoreFile *string
		want      []scoretypes.ContainerFilesElem
		wantErr   bool
	}{
		{
			name: "file with content",
			input: []scoretypes.ContainerFilesElem{
				{
					Target:  "/app/config.txt",
					Content: stringPtr("content with ${TEST_VAR}"),
				},
			},
			scoreFile: &scoreFile,
			want: []scoretypes.ContainerFilesElem{
				{
					Target:   "/app/config.txt",
					Content:  stringPtr("content with test-value"),
					NoExpand: boolPtr(true),
				},
			},
			wantErr: false,
		},
		{
			name: "file with source",
			input: []scoretypes.ContainerFilesElem{
				{
					Target: "/app/config.txt",
					Source: stringPtr(testFilePath),
				},
			},
			scoreFile: &scoreFile,
			want: []scoretypes.ContainerFilesElem{
				{
					Target:   "/app/config.txt",
					Content:  stringPtr("test content with test-value"),
					NoExpand: boolPtr(true),
				},
			},
			wantErr: false,
		},
		{
			name: "file with error in substitution",
			input: []scoretypes.ContainerFilesElem{
				{
					Target:  "/app/config.txt",
					Content: stringPtr("content with ${ERROR_VAR}"),
				},
			},
			scoreFile: &scoreFile,
			want:      nil,
			wantErr:   true,
		},
		{
			name: "file with no content or source",
			input: []scoretypes.ContainerFilesElem{
				{
					Target: "/app/config.txt",
				},
			},
			scoreFile: &scoreFile,
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertContainerFiles(tt.input, tt.scoreFile, sf)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertContainerFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("convertContainerFiles() got %v, want %v", got, tt.want)
					return
				}
				for i, v := range tt.want {
					if *got[i].Content != *v.Content {
						t.Errorf("convertContainerFiles() got[%d].Content = %v, want %v", i, *got[i].Content, *v.Content)
					}
					if got[i].Target != v.Target {
						t.Errorf("convertContainerFiles() got[%d].Target = %v, want %v", i, got[i].Target, v.Target)
					}
					if *got[i].NoExpand != *v.NoExpand {
						t.Errorf("convertContainerFiles() got[%d].NoExpand = %v, want %v", i, *got[i].NoExpand, *v.NoExpand)
					}
				}
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// Helper function to create bool pointers
func boolPtr(b bool) *bool {
	return &b
}
