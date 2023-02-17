## Build stage
FROM golang:alpine AS builder
ENV GO111MODULE=on

# Copy files to image
COPY . /falcondb/src
WORKDIR /falcondb/src

# Install Git / Dependencies
RUN apk add git ca-certificates
RUN go mod download

# Build image
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/falcondb

## Image creation stage
FROM scratch
# Copy user from build stage
COPY --from=builder /etc/passwd /etc/passwd

# Copy falcondb
COPY --from=builder /go/bin/falcondb /go/bin/falcondb
COPY --from=builder /falcondb/src/db /go/bin/db
WORKDIR /go/bin

# Set entrypoint
ENTRYPOINT ["./falcondb"]