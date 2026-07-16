FROM golang:1.26-alpine AS build
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Static build: the runtime stage has no libc to link against, and the binary
# already carries the migrations and the frontend via //go:embed.
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.22
WORKDIR /app

RUN addgroup --system threedice && adduser --system --ingroup threedice threedice
COPY --from=build /out/server /app/server

USER threedice
EXPOSE 8081
ENTRYPOINT ["/app/server"]
