FROM golang:1.14.4-alpine3.12 AS builder

WORKDIR /

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0  go build -o executeable cmd/main.go

FROM scratch
WORKDIR /
COPY --from=builder /executeable .
COPY --from=builder /config.json .

ENTRYPOINT ["/executeable","-configfile=/config.json"]
