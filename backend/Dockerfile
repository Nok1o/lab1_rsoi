FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/person-service ./cmd/person-service

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
COPY --from=build /out/person-service /usr/local/bin/person-service
USER nobody:nobody
EXPOSE 8080
ENTRYPOINT ["person-service"]
