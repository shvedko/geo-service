FROM sourcemation/golang-1.24 AS builder
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
WORKDIR /source
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main -ldflags="-s -w" ./cmd/geo/main.go

FROM scratch
COPY --from=builder /source/main /main
EXPOSE 8080
ENTRYPOINT ["/main"]
