FROM node:latest AS webpack

WORKDIR /app

# Copy package.json and install packages.
COPY package.json .
RUN npm install

# Copy the rest and build into a static directory.
COPY . .
RUN npm run build

# ============================================================================

FROM golang:alpine AS build

WORKDIR /app

# Copy and download go module dependencies.
COPY go.* .
RUN go mod download 

# Copy the the static directory from the webpack stage.
COPY --from=webpack /app/cmd/static ./cmd/static

# Copy the rest and build the binary.
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -tags prod -o app cmd/main.go cmd/static_prod.go

# ============================================================================

FROM scratch

# Copy the binary from the build stage.
COPY --from=build /app/app .

CMD ["./app"]
