FROM golang:1.25.1-alpine AS build

WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/web_server ./cmd/web_server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/agent_server ./cmd/agent_server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/agent ./cmd/agent
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/report ./cmd/report

FROM alpine:3.20

WORKDIR /app
RUN adduser -D -u 10001 appuser
COPY --from=build /out/ /app/
USER appuser

ENTRYPOINT ["/app/web_server"]
