
# Getting Started with score-aca (podinfo example)

This tutorial shows how to use `score-aca` to deploy the `podinfo` application to Azure Container Apps using the provided `score.yaml`.

## Prerequisites

1. `score-aca` binary installed and available in your PATH
2. [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) installed and configured
3. An active Azure subscription with permissions to create resources

## 1. Verify Installation

Check that `score-aca` is installed:

```sh
score-aca --version
```

## 2. Prepare Your Project

Create a new directory and add your `score.yaml`:

```sh
mkdir podinfo-aca
cd podinfo-aca
```

Initialize the project by running:

```sh
score-aca init
```

This creates a sample `score.yaml` with the following for the `podinfo` application:

```yaml
apiVersion: score.dev/v1b1
containers:
    main:
        image: stefanprodan/podinfo
metadata:
    name: example
service:
    ports:
        web:
            port: 9898
```

You can customize the `score.yaml` as needed. The provided example includes a health check configuration and environment variables:

```yaml
apiVersion: score.dev/v1b1

metadata:
  name: podinfo-app

containers:
  podinfo:
    image: stefanprodan/podinfo
    # Environment variables
    variables:
      PODINFO_UI_COLOR: "#a3c94cff"
      PODINFO_UI_MESSAGE: "Welcome to Podinfo"
      PODINFO_METRICS: "true"
      PODINFO_DEBUG: "false"

    # Health checks
    livenessProbe:
      httpGet:
        path: /healthz     # Health check endpoint
        port: 9999         # Port to check
        scheme: HTTP       # HTTP or HTTPS
        # Optional HTTP headers
        httpHeaders:
          - name: Custom-Header
            value: liveness-check

service:
  ports:
    http:
      port: 9898
      targetPort: 9898
      protocol: TCP
```

## 3. Generate the Bicep Manifest

Run the following command to generate the Azure Bicep manifest:

```sh
score-aca generate score.yaml
```

This creates a `manifest.bicep` file in your directory with the necessary Azure resources defined.

## 4. Deploy to Azure

First, log in to Azure:

```sh
az login
# or
az login --use-device-code
```

Set your desired location and resource group in environment variables for easier reuse:

```sh
export LOCATION="eastus"
export RESOURCE_GROUP="rg-score-aca"
```

Create a resource group (if needed):

```sh
az group create \
  --name ${RESOURCE_GROUP} \
  --location ${LOCATION}
```

Deploy the Bicep template:

```sh
az deployment group create \
  --resource-group "${RESOURCE_GROUP}" \
  --template-file ./manifest.bicep
```

## 5. Verify the Deployment

List your container apps:

```sh
az containerapp list \
  --resource-group ${RESOURCE_GROUP} \
  --output table
```

Get the application URL:

```sh
az containerapp show \
  --name podinfo-app-container-app \
  --resource-group ${RESOURCE_GROUP} \
  --query "properties.configuration.ingress.fqdn" \
  -o tsv
```

Visit the URL in your browser to see the running application.

## 6. Clean Up

To remove all resources:

```sh
az group delete \
  --name ${RESOURCE_GROUP} \
  --yes --no-wait
```

## Troubleshooting

- View logs:

```sh
az containerapp logs show \
  --name podinfo-app-container-app \
  --resource-group ${RESOURCE_GROUP}
```

- Validate the Bicep template:

```sh
az bicep build --file manifest.bicep
```

- Check deployment errors in the Azure Portal under the resource group's deployment history.

## Next Steps

- Customize the `score.yaml` for your needs (add environment variables, health checks, etc.)
- Learn more about [Score specification](https://score.dev/)
- Explore [Azure Container Apps documentation](https://docs.microsoft.com/en-us/azure/container-apps/)
