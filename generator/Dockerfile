# Start with the OTEL collector stage
FROM otel/opentelemetry-collector-contrib:latest AS collector

# Final stage
FROM golang:1.24

# Copy the collector binary from the collector stage
COPY --from=collector /otelcol-contrib /otelcol-contrib

WORKDIR /go/src/app

# Copy and build the log generator
COPY go.mod go.sum ./
COPY *.go ./
COPY flog/ flog/
COPY log/ log/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /generator

COPY <<EOF /start.sh
#!/bin/sh

# Start the log generator
/generator \$@
EOF

RUN chmod +x /start.sh

ENTRYPOINT ["/start.sh"]
