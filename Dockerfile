FROM golang:1.27-alpine AS build
RUN apk add --no-cache git
WORKDIR /src/bdo-rest-api
COPY . .
RUN go mod download
RUN go build -o /bdo-rest-api -ldflags="-s -w" .

FROM alpine:latest AS bin
RUN addgroup --system --gid 1001 go
RUN adduser --system --uid 1001 go
COPY --from=build --chown=go:go /bdo-rest-api .
USER go
ENV PROXY=
EXPOSE 8001
ENTRYPOINT ["/bdo-rest-api"]
