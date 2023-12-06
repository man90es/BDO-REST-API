# Building from source
Building the scraper from source may be preferable in some cases. This way, the app has a smaller footprint than a [Docker image](https://hub.docker.com/r/man90/bdo-rest-api).

## Prerequisites:
- GNU/Linux (other platforms should work as well but I haven't tested them)
- Go compiler >=v1.15

## Compilation
Use this command to compile the app:
```bash
go build
```

## Environment variables
Catch requests on a specific port (8001 by default):
```bash
export PORT=3000
```

Use a proxy to make requests to BDO servers (direct by default):
```bash
export PROXY=http://123.123.123.123:8080
# or
export PROXY="http://123.123.123.123:8080 http://124.124.124.124:8081"
```

## Flags
These flags override environment variables.
```
-cachecap int
	Cache capacity (default 10000)
-cachettl int
	Cache TTL in minutes (default 180)
-port int
	Port to catch requests on (default 8001)
-proxy string
	Open proxy address to make requests to BDO servers
-verbose
	Print out additional logs into stdout
```

Use them like this:
```bash
./bdo-rest-api -proxy="http://192.168.0.0.1:8080" -cachettl=30
```
