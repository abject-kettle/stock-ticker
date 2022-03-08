# stock-ticker

This is a stock ticker web service for reporting the historical closing price of a particular stock.

## Client Usage

Example usage:
```shell
# curl http://stock-ticker.test
{"average":"295.9420","historical":[{"date":"2022-03-04","close":"289.8600"},{"date":"2022-03-03","close":"295.9200"},{"date":"2022-03-02","close":"300.1900"},{"date":"2022-03-01","close":"294.9500"},{"date":"2022-02-28","close":"298.7900"}]}
```

The server responds with a JSON document that includes the historical closing prices for the preceding set number of
days.

### API

The JSON document contains two fields: `average` and `historical`.
```
"average": The average closing price over the set number of days.
"historical": A list of historical closing prices covering the set number of days.
    "date": The day of the price.
    "close": The closing price on the day.
```

## Deploy locally

### Build binary
```shell
go build -o bin/stock-ticker ./cmd
```

### Run binary
```shell
SYMBOL=MSFT NDAYS=5 APIKEY=C227WD9W3LUVKVV9 PORT=8080 ./bin/stock-ticker 
```

### Query stock ticker
```shell
curl http://localhost:8080
```

## Deploy locally in container

### Build image
```shell
podman build --tag stock-ticker .
```

### Run image
```shell
podman run -e SYMBOL=MSFT -e NDAYS=5 -e APIKEY=C227WD9W3LUVKVV9 -e PORT=8080 -p 8080:8080/tcp stock-ticker
```

### Query stock ticker
```shell
curl http://localhost:8080
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

### Query stock ticker
```shell
curl http://stock-ticker.test
```

Note that your URL may differ depending upon the configuration of your kubernetes cluster and the configuration that you
specify in `./manifests/kustomization.yaml`. See the Customization section below for more details.

### Customization

The `./manifests/kustomization.yaml` file, and its dependent files, can be adjusted to fit your desired deployment

#### Ticker Symbol

Set ticker symbol to lookup by setting the value of `SYMBOL` in `./manifests/server-config.properties`.
The default ticker symbol is MSFT.

#### Number of Days

Set the number of days of prices to lookup by setting the value of `NDAYS` in `./manifests/server-config.properties`.
The default number of days is 5.

#### API Key

Set the API key for the price service by setting the value of `APIKEY` in `./manifests/apikey.properties`.

#### Namespace

The namespace in which to locate the kubernetes resources can be set by changing the `namespace` field in `./manifests/kustomization.yaml`.
The default namespace is "stock-ticker".

#### Image

The image to use for the stock-ticker server can be set by changing the `images` field in `./manifests/kustomization.yaml`.
The default image is quay.io/matthew_staebler/stock-ticker:latest.

#### Replicas

The number of replicas to use for the stock-ticker server can be set by changing the `replicas` field in `./manifests/kustomization.yaml`.
The default number of replicas to use is 3.

#### URL

The URL to use for accessing the stock-ticker server can be set by changing the `patches` field in `./manifests/kustomization.yaml`.
The default URL is http://stock-ticker.test.
