# stock-ticker
Stock ticker web service

## Build binary
```shell
go build -o bin/stock-ticker ./cmd
```

## Run binary
```shell
SYMBOL=MSFT NDAYS=5 APIKEY=C227WD9W3LUVKVV9 PORT=8080 ./bin/stock-ticker 
```

## Query stock ticker
```shell
curl http://localhost:8080 | jq .
```

## Build image
```shell
podman build --tag stock-ticker .
```

## Run image
```shell
podman run -e SYMBOL=MSFT -e NDAYS=5 -e APIKEY=C227WD9W3LUVKVV9 -e PORT=8080 -p 8080:8080/tcp stock-ticker
```

## Deploy to kubernetes

### Prerequisites
* A running cluster.
* An ingress controller running in the cluster.
* A default ingress class.
* A created namespace for the deployment. The default namespace is "stock-ticker".

### Deploying
See the section below on customization to fit your desired deployment.

```shell
kubectl apply -k ./manifests
```

### Customization

The `./manifests/kustomization.yaml` file, and its dependent files, can be adjusted to fit your desired deployment

#### Ticker Symbol

Set ticker symbol to lookup by setting the value of `SYMBOL` in `./manifests/server-config.properties`.

#### Number of Days

Set the number of days of prices to lookup by setting the value of `NDAYS` in `./manifests/server-config.properties`.

#### API Key

Set the API key for the price service by setting the value of `APIKEY` in `./manifests/apikey.properties`.

#### Namespace

The namespace in which to locate the kubernetes resources can be set by changing the `namespace` field in `./manifests/kustomization.yaml`.

#### Image

The image to use for the stock-ticker server can be set by changing the `images` field in `./manifests/kustomization.yaml`.

#### Replicas

The number of replicas to use for the stock-ticker server can be set by changing the `replicas` field in `./manifests/kustomization.yaml`.

#### URL

The URL to use for accessing the stock-ticker server can be set by changing the `patches` field in `./manifests/kustomization.yaml`.
