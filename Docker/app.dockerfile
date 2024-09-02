ARG GO_BUILD_TAGS=all

# ============================================================================

FROM node:latest AS webpack

WORKDIR /app

COPY package.json .
RUN npm install

COPY . .
RUN npm run build

# ============================================================================

FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download 

COPY --from=webpack /app/cmd/static ./cmd/static

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -tags prod -o app ./cmd

# ============================================================================

FROM scratch

COPY --from=build /app/app .

CMD ["./app"]
