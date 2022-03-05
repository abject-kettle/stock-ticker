# stock-ticker
Stock ticker web service

## Build binary
```shell
go build -o bin/stock-ticker ./cmd
```

## Run binary
```shell
SYMBOL=MSFT NDAYS=5 APIKEY=C227WD9W3LUVKVV9 ./bin/stock-ticker 
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
podman run -e SYMBOL=MSFT -e NDAYS=5 -e APIKEY=C227WD9W3LUVKVV9 -p 8080:8080/tcp stock-ticker
```