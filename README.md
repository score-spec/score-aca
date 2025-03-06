# score-aca

`score-aca` is a Score implementation of the [Score specification](https://github.com/score-spec/spec) for [Azure Container Apps](https://azure.microsoft.com/products/container-apps).

## Demo

Write the following to `score.yaml`:

```yaml
apiVersion: score.dev/v1b1
metadata:
    name: example
containers:
    main:
        image: stefanprodan/podinfo
```

### Generate Bicep file

```sh
go run ./cmd/score-aca init
go run ./cmd/score-aca generate score.yaml
```

The output `manifests.bicep` contains the following which indicates:

1. The environment that will hold all other resources.
2. The Azure Container App based on Container Image.

```bicep
resource environment 'Microsoft.App/managedEnvironments@2024-03-01' = {
  name: 'example-environment'
  location: resourceGroup().location
  properties: {
    appLogsConfiguration: {
      destination: 'azure-monitor'
    }
  }
}

resource app 'Microsoft.App/containerApps@2024-03-01' = {
  name: 'example-app'
  location: resourceGroup().location
  properties: {
    managedEnvironmentId: environment.id
    configuration: {
      ingress: {
        external: true
      }
    }
    template: {
      containers: [
        {
          image: 'stefanprodan/podinfo'
          name: 'podinfo'
        }
      ]
    }
  }
}
```

### Deploy Container App in Azure

```sh
export LOCATION="eastus"
export RESOURCE_GROUP="rg-score-aca"

az login
# or
az login --use-device-code

az group create --name ${RESOURCE_GROUP} --location ${LOCATION}

az deployment group create --resource-group "${RESOURCE_GROUP}" --template-file ./manifests.bicep
```

## A note on licensing

Most code files here retain the Apache licence header since they were copied or adapted from the reference `score-compose` which is Apache licensed. Any modifications to these files should retain the Apache licence and attribution.
