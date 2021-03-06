FROM golang:1.12-alpine as builder
RUN apk update && \
    apk add --update make git
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
WORKDIR /src
RUN go build -o sdmap-test-agent ./cmd/mqtt-pub/*.go

FROM golang:1.12-alpine
RUN apk update
WORKDIR /app
COPY --from=builder /src/pkg/decoder/samples/logs/bsm-10-23.log .
COPY --from=builder /src/pkg/decoder/samples/logs/spat.log .
COPY --from=builder /src/sdmap-test-agent .
ENTRYPOINT [ "/app/sdmap-test-agent" ]