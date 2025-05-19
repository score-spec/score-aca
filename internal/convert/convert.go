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
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/score-spec/score-go/framework"
	scoretypes "github.com/score-spec/score-go/types"

	"github.com/score-spec/score-aca/internal/state"
)

// ContainerAppProperties represents the properties of an Azure Container App
type ContainerAppProperties struct {
	Configuration ContainerAppConfiguration `json:"configuration"`
	Template      ContainerAppTemplate      `json:"template"`
	// Optional fields
	EnvironmentID       string `json:"environmentId,omitempty"`
	WorkloadProfileName string `json:"workloadProfileName,omitempty"`
}

// ContainerAppConfiguration represents the configuration of an Azure Container App
type ContainerAppConfiguration struct {
	ActiveRevisionsMode  string               `json:"activeRevisionsMode,omitempty"`
	Ingress              *ContainerAppIngress `json:"ingress,omitempty"`
	MaxInactiveRevisions int                  `json:"maxInactiveRevisions,omitempty"`
}

// ContainerAppIngress represents the ingress configuration of an Azure Container App
type ContainerAppIngress struct {
	External               bool                    `json:"external"`
	TargetPort             int                     `json:"targetPort"`
	Transport              string                  `json:"transport,omitempty"`
	AllowInsecure          bool                    `json:"allowInsecure,omitempty"`
	Traffic                []ContainerAppTraffic   `json:"traffic,omitempty"`
	IPSecurityRestrictions []IPSecurityRestriction `json:"ipSecurityRestrictions,omitempty"`
}

// ContainerAppTraffic represents the traffic configuration of an Azure Container App
type ContainerAppTraffic struct {
	Weight         int    `json:"weight"`
	LatestRevision bool   `json:"latestRevision"`
	RevisionName   string `json:"revisionName,omitempty"`
	Label          string `json:"label,omitempty"`
}

// IPSecurityRestriction represents an IP security restriction
type IPSecurityRestriction struct {
	Name           string `json:"name"`
	Action         string `json:"action"`
	IPAddressRange string `json:"ipAddressRange"`
	Description    string `json:"description,omitempty"`
}

// ContainerAppTemplate represents the template of an Azure Container App
type ContainerAppTemplate struct {
	Containers []ContainerAppContainer `json:"containers"`
}

// ContainerAppContainer represents a container in an Azure Container App
type ContainerAppContainer struct {
	Name      string                `json:"name"`
	Image     string                `json:"image"`
	Command   []string              `json:"command,omitempty"`
	Args      []string              `json:"args,omitempty"`
	Env       []ContainerAppEnv     `json:"env,omitempty"`
	Resources ContainerAppResources `json:"resources"`
	Probes    []ContainerAppProbe   `json:"probes,omitempty"`
}

// ContainerAppEnv represents an environment variable in an Azure Container App
type ContainerAppEnv struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

// ContainerAppResources represents the resources of a container in an Azure Container App
type ContainerAppResources struct {
	CPU    float64 `json:"cpu"`
	Memory string  `json:"memory"`
}

// ContainerAppProbe represents a probe in an Azure Container App
type ContainerAppProbe struct {
	Type    string               `json:"type"`
	HTTPGet *ContainerAppHTTPGet `json:"httpGet,omitempty"`
}

// ContainerAppHTTPGet represents an HTTP GET probe in an Azure Container App
type ContainerAppHTTPGet struct {
	Path        string                   `json:"path"`
	Port        int                      `json:"port"`
	Scheme      string                   `json:"scheme,omitempty"`
	HTTPHeaders []ContainerAppHTTPHeader `json:"httpHeaders,omitempty"`
}

// ContainerAppHTTPHeader represents an HTTP header in an Azure Container App
type ContainerAppHTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Workload converts a Score workload to a Bicep manifest
func Workload(currentState *state.State, workloadName string) (string, error) {
	resOutputs, err := currentState.GetResourceOutputForWorkload(workloadName)
	if err != nil {
		return "", fmt.Errorf("failed to generate outputs: %w", err)
	}
	sf := framework.BuildSubstitutionFunction(currentState.Workloads[workloadName].Spec.Metadata, resOutputs)

	spec := currentState.Workloads[workloadName].Spec
	containers := maps.Clone(spec.Containers)
	for containerName, container := range containers {
		if container.Variables, err = convertContainerVariables(container.Variables, sf); err != nil {
			return "", fmt.Errorf("workload: %s: container: %s: variables: %w", workloadName, containerName, err)
		}

		if container.Files, err = convertContainerFiles(container.Files, currentState.Workloads[workloadName].File, sf); err != nil {
			return "", fmt.Errorf("workload: %s: container: %s: files: %w", workloadName, containerName, err)
		}
		containers[containerName] = container
	}
	spec.Containers = containers
	resources := maps.Clone(spec.Resources)
	for resName, res := range resources {
		resUid := framework.NewResourceUid(workloadName, resName, res.Type, res.Class, res.Id)
		resState, ok := currentState.Resources[resUid]
		if !ok {
			return "", fmt.Errorf("workload '%s': resource '%s' (%s) is not primed", workloadName, resName, resUid)
		}
		res.Params = resState.Params
		resources[resName] = res
	}
	spec.Resources = resources

	// Convert the Score workload to a Bicep manifest
	bicepManifest, err := convertToBicep(spec, workloadName)
	if err != nil {
		return "", fmt.Errorf("workload: %s: failed to convert to Bicep: %w", workloadName, err)
	}

	return bicepManifest, nil
}

// convertToBicep converts a Score workload to a Bicep manifest
func convertToBicep(spec scoretypes.Workload, workloadName string) (string, error) {
	// Create the Bicep manifest
	bicepContent := generateBicepHeader()

	// Add parameters
	params, err := generateBicepParameters(workloadName)
	if err != nil {
		return "", fmt.Errorf("failed to generate Bicep parameters: %w", err)
	}
	bicepContent += params

	// Add container app environment
	bicepContent += generateContainerAppEnvironment()

	// Add container app
	containerApp, err := generateContainerApp(spec, workloadName)
	if err != nil {
		return "", fmt.Errorf("failed to generate container app: %w", err)
	}
	bicepContent += containerApp

	// Add outputs
	bicepContent += bicepOutputs

	return bicepContent, nil
}

// generateBicepHeader generates the header of the Bicep manifest
func generateBicepHeader() string {
	return bicepHeader
}

// generateBicepParameters generates the parameters section of the Bicep manifest
func generateBicepParameters(workloadName string) (string, error) {

	t, err := template.New("bicepParameters").Parse(bicepParameters)
	if err != nil {
		return "", err
	}
	// Create a buffer to hold the generated parameters
	var buf bytes.Buffer

	// Execute the template with the workload name
	if err := t.Execute(&buf, map[string]string{"WorkloadName": workloadName}); err != nil {
		return "", err
	}
	// Return the generated parameters as a string
	return buf.String(), nil
}

// generateContainerAppEnvironment generates the container app environment section of the Bicep manifest
func generateContainerAppEnvironment() string {
	return bicepContainerAppEnvironment
}

// generateContainerApp generates the container app section of the Bicep manifest
func generateContainerApp(spec scoretypes.Workload, workloadName string) (string, error) {
	// Create the container app properties
	properties, err := createContainerAppProperties(spec)
	if err != nil {
		return "", fmt.Errorf("failed to create container app properties: %w", err)
	}

	// Convert properties to YAML for debugging
	// propertiesYAML, _ := yaml.Marshal(properties)
	// slog.Debug("Generated container app properties", slog.Any("properties", string(propertiesYAML)))

	t, err := template.New("bicepContainerApp").Parse(bicepContainerApp)
	if err != nil {
		return "", err
	}
	// Create a buffer to hold the generated container app
	var buf bytes.Buffer

	data := struct {
		WorkloadName string
		Properties   *ContainerAppProperties
		Spec         scoretypes.Workload
	}{
		WorkloadName: workloadName,
		Properties:   properties,
		Spec:         spec,
	}

	// Execute the template with the properties
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	// Return the generated container app as a string
	return buf.String(), nil

}

// createContainerAppProperties creates the properties of an Azure Container App from a Score workload
func createContainerAppProperties(spec scoretypes.Workload) (*ContainerAppProperties, error) {
	properties := &ContainerAppProperties{
		Configuration: ContainerAppConfiguration{
			ActiveRevisionsMode: "Single",
		},
		Template: ContainerAppTemplate{
			Containers: []ContainerAppContainer{},
		},
	}

	// Set ingress if service is defined
	if spec.Service != nil && len(spec.Service.Ports) > 0 {
		// Use the first port for ingress
		var port int
		for _, p := range spec.Service.Ports {
			if p.TargetPort != nil {
				port = *p.TargetPort
			}
			if port == 0 {
				port = p.Port
			}
			break
		}
		properties.Configuration.Ingress = &ContainerAppIngress{
			External:   true,
			TargetPort: port,
			Transport:  "auto",
		}
	}

	// Add containers
	for name, container := range spec.Containers {
		// Create container
		containerApp := ContainerAppContainer{
			Name:  name,
			Image: container.Image,
			Resources: ContainerAppResources{
				CPU:    0.25,    // Default CPU
				Memory: "0.5Gi", // Default memory
			},
		}

		// Add command if any
		if len(container.Command) > 0 {
			containerApp.Command = container.Command
		}

		// Add args if any
		if len(container.Args) > 0 {
			containerApp.Args = container.Args
		}

		// Add environment variables
		for key, value := range container.Variables {
			env := ContainerAppEnv{
				Name:  key,
				Value: value,
			}

			containerApp.Env = append(containerApp.Env, env)
		}

		// Add resources if defined
		if container.Resources != nil {
			if container.Resources.Limits != nil {
				slog.Warn(fmt.Sprintf("%s: Resource limits are not supported in Azure Container Apps. Set the wanted values in the requests section.", name))
			}
			// TODO: Evaluate Resource Requests for correct values from Azure
			if container.Resources.Requests != nil {
				if container.Resources.Requests.Cpu != nil {
					cpuValue, err := parseCPU(*container.Resources.Requests.Cpu)
					if err == nil {
						containerApp.Resources.CPU = cpuValue
					}
				}
				if container.Resources.Requests.Memory != nil {
					containerApp.Resources.Memory = *container.Resources.Requests.Memory
				}
			}
		}

		// Add probes if defined
		if container.LivenessProbe != nil {
			probe := ContainerAppProbe{}
			if container.LivenessProbe.HttpGet != nil {
				probe.HTTPGet = &ContainerAppHTTPGet{
					Path: container.LivenessProbe.HttpGet.Path,
					Port: container.LivenessProbe.HttpGet.Port,
				}
				if container.LivenessProbe.HttpGet.Scheme != nil {
					probe.HTTPGet.Scheme = string(*container.LivenessProbe.HttpGet.Scheme)
				}

			} else if spec.Containers[name].LivenessProbe != nil {
				if spec.Containers[name].LivenessProbe.HttpGet != nil {
					probe.HTTPGet = &ContainerAppHTTPGet{
						Path: spec.Containers[name].LivenessProbe.HttpGet.Path,
						Port: spec.Containers[name].LivenessProbe.HttpGet.Port,
					}
					if spec.Containers[name].LivenessProbe.HttpGet.Scheme != nil {
						probe.HTTPGet.Scheme = string(*spec.Containers[name].LivenessProbe.HttpGet.Scheme)
					}
				}
			}
			if probe != (ContainerAppProbe{}) {
				probe.Type = "Liveness"
				containerApp.Probes = append(containerApp.Probes, probe)
			}
		}

		if container.ReadinessProbe != nil {
			probe := ContainerAppProbe{}
			if container.ReadinessProbe != nil {
				if container.ReadinessProbe.HttpGet != nil {
					probe.HTTPGet = &ContainerAppHTTPGet{
						Path: container.ReadinessProbe.HttpGet.Path,
						Port: container.ReadinessProbe.HttpGet.Port,
					}
					if container.ReadinessProbe.HttpGet.Scheme != nil {
						probe.HTTPGet.Scheme = string(*container.ReadinessProbe.HttpGet.Scheme)
					}
				}
				containerApp.Probes = append(containerApp.Probes, probe)
			} else if spec.Containers[name].ReadinessProbe != nil {
				if spec.Containers[name].ReadinessProbe.HttpGet != nil {
					probe.HTTPGet = &ContainerAppHTTPGet{
						Path: spec.Containers[name].ReadinessProbe.HttpGet.Path,
						Port: spec.Containers[name].ReadinessProbe.HttpGet.Port,
					}
					if spec.Containers[name].ReadinessProbe.HttpGet.Scheme != nil {
						probe.HTTPGet.Scheme = string(*spec.Containers[name].ReadinessProbe.HttpGet.Scheme)
					}
				}
			}
			if probe != (ContainerAppProbe{}) {
				probe.Type = "Readiness"
				containerApp.Probes = append(containerApp.Probes, probe)
			}
		}

		properties.Template.Containers = append(properties.Template.Containers, containerApp)
	}

	return properties, nil
}

// parseCPU parses a CPU value from a string to a float64
func parseCPU(cpu string) (float64, error) {
	// Check if the CPU value is in millicores (e.g., "500m")
	if strings.HasSuffix(cpu, "m") {
		millicore := strings.TrimSuffix(cpu, "m")
		value, err := parseFloat(millicore)
		if err != nil {
			return 0, err
		}
		return value / 1000, nil
	}

	// Otherwise, parse as a float
	return parseFloat(cpu)
}

// parseFloat parses a string to a float64
func parseFloat(s string) (float64, error) {
	var value float64
	_, err := fmt.Sscanf(s, "%f", &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func convertContainerVariables(input scoretypes.ContainerVariables, sf func(string) (string, error)) (map[string]string, error) {
	outMap := make(map[string]string, len(input))
	for key, value := range input {
		out, err := framework.SubstituteString(value, sf)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}
		outMap[key] = out
	}
	return outMap, nil
}

func convertContainerFiles(input []scoretypes.ContainerFilesElem, scoreFile *string, sf func(string) (string, error)) ([]scoretypes.ContainerFilesElem, error) {
	outSlice := make([]scoretypes.ContainerFilesElem, 0, len(input))
	for i, fileElem := range input {
		var content string
		if fileElem.Content != nil {
			content = *fileElem.Content
		} else if fileElem.Source != nil {
			sourcePath := *fileElem.Source
			if !filepath.IsAbs(sourcePath) && scoreFile != nil {
				sourcePath = filepath.Join(filepath.Dir(*scoreFile), sourcePath)
			}
			if rawContent, err := os.ReadFile(sourcePath); err != nil {
				return nil, fmt.Errorf("%d: source: failed to read file '%s': %w", i, sourcePath, err)
			} else {
				content = string(rawContent)
			}
		} else {
			return nil, fmt.Errorf("%d: missing 'content' or 'source'", i)
		}

		var err error
		if fileElem.NoExpand == nil || !*fileElem.NoExpand {
			content, err = framework.SubstituteString(string(content), sf)
			if err != nil {
				return nil, fmt.Errorf("%d: failed to substitute in content: %w", i, err)
			}
		}
		fileElem.Source = nil
		fileElem.Content = &content
		bTrue := true
		fileElem.NoExpand = &bTrue
		outSlice = append(outSlice, fileElem)
	}
	return outSlice, nil
}
