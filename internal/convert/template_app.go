package convert

// TODO: Pull Container App version from Azure
const bicepContainerApp = `{{ define "bicepContainerApp" }}
// Container App
resource containerApp 'Microsoft.App/containerApps@2024-03-01' = {
  name: containerAppName
  location: location
  properties: {
    environmentId: containerAppEnvironment.id
    configuration: {
      {{- if and (ne .Spec.Service nil) (gt (len .Spec.Service.Ports) 0) }}
      ingress: {
        external: true
        {{- if (ne .Properties.Configuration.Ingress.TargetPort 0) }}
        targetPort: {{ .Properties.Configuration.Ingress.TargetPort }}
        {{- end }}
      }{{- end }}
    }
    template: {
      containers: [
        {{- range $containerName, $container := .Spec.Containers }}
        {
          name: '{{ $containerName }}'
          image: '{{ $container.Image }}'
          {{- if (gt (len $container.Command) 0) }}
          command: [
            {{- range $i, $cmd := $container.Command }}
            '{{ $cmd }}'{{ if (gt (len $container.Command) $i) }},{{ end }}
            {{- end }}
          ]
          {{- end }}

          {{- if (gt (len $container.Args) 0) }}
          args: [
            {{- range $i, $arg := $container.Args }}
            '{{ $arg }}'{{ if (gt (len $container.Args) $i) }},{{ end }}
            {{- end }}
          ]{{- end }}

          {{- if (ne $container.Resources nil) }}
          resources: {
            {{- if (ne $container.Resources.Requests nil) }}
            cpu: json('{{ $container.Resources.Requests.Cpu }}')
            memory: '{{ $container.Resources.Requests.Memory }}'
            {{- end }}
          }{{- end }}

          {{- if or (ne $container.LivenessProbe nil) (ne $container.ReadinessProbe nil) }}
          probes: [
            {{- if (ne $container.LivenessProbe nil) }}
            {
              type: 'liveness'
              initialDelaySeconds: 15
              periodSeconds: 30
              failureThreshold: 3
              timeoutSeconds: 1
              {{- if (ne $container.LivenessProbe.HttpGet nil) }}
              httpGet: {
                {{- if (ne $container.LivenessProbe.HttpGet.Port 0) }}
                port: {{ $container.LivenessProbe.HttpGet.Port }}
                {{- end}}
                {{- if (ne $container.LivenessProbe.HttpGet.Path "") }}
                path: '{{ $container.LivenessProbe.HttpGet.Path }}'
                {{- end }}
                {{- if (ne $container.LivenessProbe.HttpGet.Host nil) }}
                host: '{{ $container.LivenessProbe.HttpGet.Host }}'
                {{- end }}
                {{- if (ne $container.LivenessProbe.HttpGet.Scheme nil) }}
                scheme: '{{ $container.LivenessProbe.HttpGet.Scheme }}'
                {{- end }}
              }{{- end }}
            }{{- end }}
            {{- if (ne $container.ReadinessProbe nil) }}
            {
              type: 'readiness'
              initialDelaySeconds: 15
              periodSeconds: 30
              failureThreshold: 3
              timeoutSeconds: 1
              {{- if (ne $container.ReadinessProbe.HttpGet nil) }}
              httpGet: {
                {{- if (ne $container.ReadinessProbe.HttpGet.Port 0) }}
                port: {{ $container.ReadinessProbe.HttpGet.Port }}
                {{- end}}
                {{- if (ne $container.ReadinessProbe.HttpGet.Path "") }}
                path: '{{ $container.ReadinessProbe.HttpGet.Path }}'
                {{- end }}
                {{- if (ne $container.ReadinessProbe.HttpGet.Host nil) }}
                host: '{{ $container.ReadinessProbe.HttpGet.Host }}'
                {{- end }}
                {{- if (ne $container.ReadinessProbe.HttpGet.Scheme nil) }}
                scheme: '{{ $container.ReadinessProbe.HttpGet.Scheme }}'
                {{- end }}
              }{{- end }}
            }{{- end }}
          ]{{- end }}

          {{- if (gt (len $container.Variables) 0) }}
          env: [
            {{- range $variableName, $variableValue := $container.Variables }}
            {
              name: '{{ $variableName }}'
              value: '{{ $variableValue }}'
            }{{- end }}
          ]{{- end }}
        }{{- end }}
      ]
    }
  }
}
{{ end }}
`
