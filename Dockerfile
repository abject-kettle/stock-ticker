FROM golang:1.17-alpine AS builder
WORKDIR /go/src/stock-ticker
COPY . .
RUN go build -o bin/stock-ticker ./cmd

FROM alpine:latest
COPY --from=builder /go/src/stock-ticker/bin/stock-ticker /bin/stock-ticker
ENV PATH /bin
ENTRYPOINT ["/bin/stock-ticker"]