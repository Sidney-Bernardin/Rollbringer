FROM golang:alpine AS build

WORKDIR /app

# Copy and download go module dependencies.
COPY go.* .
RUN go mod download 

# Copy the rest and build the binary.
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o app cmd/*

# ============================================================================

FROM scratch

# Copy the binary from the build stage.
COPY --from=build /app/app .

CMD ["./app"]
