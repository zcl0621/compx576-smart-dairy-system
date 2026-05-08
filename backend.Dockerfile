FROM golang:1.25.1-bookworm AS build

WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/web_server ./cmd/web_server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/agent_server ./cmd/agent_server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/agent ./cmd/agent
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/report ./cmd/report

FROM debian:bookworm-slim

WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && useradd -r -u 10001 -g users appuser \
    && rm -rf /var/lib/apt/lists/*
COPY --from=build /out/ /app/
USER appuser

ENTRYPOINT ["/app/web_server"]
