FROM golang:1.21rc2-alpine3.18 AS build
RUN apk add --no-cache git
WORKDIR /src/bdo-rest-api
COPY . .
RUN go mod download
ARG tags=none
RUN go build -tags $tags -o /bdo-rest-api -ldflags="-s -w" .

FROM alpine:3.14 AS bin
RUN addgroup --system --gid 1001 go
RUN adduser --system --uid 1001 go
COPY --from=build --chown=go:go /bdo-rest-api .
USER go
ENV PROXY=
EXPOSE 8001
CMD ["/bdo-rest-api"]
