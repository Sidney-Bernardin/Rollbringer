FROM node:20.5.1 AS static

COPY package.json .

RUN npm install && npm run build

# ============================================================================

FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.* .
RUN go mod download 

COPY . .
COPY --from=static static static

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o app . 

# ============================================================================

FROM scratch

COPY --from=build /app/app .
COPY --from=build /app/static static

CMD ["./app"]
