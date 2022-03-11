FROM golang:1.17-alpine AS build
RUN apk add --no-cache git
WORKDIR /src/bdo-rest-api
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./bin/bdo-rest-api .

FROM alpine:3.14 AS bin
COPY --from=build /src/bdo-rest-api/bin/bdo-rest-api /app/bdo-rest-api
EXPOSE 8001
CMD ["/app/bdo-rest-api"]
