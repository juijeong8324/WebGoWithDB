# ---- Stage 1: Builder ----
# Use Go + Alpine as the build environment
FROM golang:1.20.1-alpine3.17 AS builder
WORKDIR /work
COPY go.mod ./
RUN go mod download
COPY . /work/
RUN go build -o server .

# ---- Stage 2: Runtime ----
# Use a minimal Alpine image (no Go compiler needed at runtime)
FROM alpine:3.14
COPY --from=builder /work/server /work/server
# Run the server when the container starts
ENTRYPOINT ["/work/server"]
