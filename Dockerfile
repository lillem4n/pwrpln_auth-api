FROM golang:1.17.5-alpine3.14 AS builder

# Install missing pkgs
RUN apk add --no-cache git

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set workdir in GOPATH
WORKDIR $GOPATH/src/gitlab.larvit.se/power-plan/auth

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY ./src ./src

# Build the application
RUN go build -o /build/main ./src

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM scratch

COPY --from=builder /dist/main /

# Command to run
ENTRYPOINT ["/main"]