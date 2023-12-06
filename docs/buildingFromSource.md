# Building from source
Building the scraper from source may be preferable in some cases. This way, the app has a smaller footprint and gives you more control than a [Docker image](https://hub.docker.com/r/man90/bdo-rest-api).

## Prerequisites:
- GNU/Linux (other platforms should work as well but I haven't tested them)
- Go compiler >=v1.15

## Compilation
By default, scraped results are cached in memory and stored for up to 2 hours. This helps to ease the pressure on BDO servers and decreases the response time tremendously (for cached responses). Use this command to compile the app:
```bash
go build
```

If you don't want to cache scraped results (e.g., if you are 100% sure that there will be no similar requests sent to the API), you can also use this command instead:
```bash
go build -tags "cacheless"
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