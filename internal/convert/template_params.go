package convert

const bicepParameters = `{{ define "bicepParameters" }}
// Parameters
param environmentName string = '{{ .WorkloadName }}-environment'
param containerAppName string = '{{ .WorkloadName }}-container-app'
param location string = resourceGroup().location

{{ end }}
`
