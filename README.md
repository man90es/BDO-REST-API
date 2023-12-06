# BDO-REST-API

A scraper for Black Desert Online player in-game data with a REST API. It currently supports European, North American, and South American servers (Korean server support is in progress).

## Projects using this API
- [BDO Leaderboards](https://bdo.hemlo.cc/leaderboards) ([Source](https://github.com/man90es/BDO-Leaderboards)): web-based leaderboards for Black Desert Online.
- [Ikusa](https://ikusa.site) ([Source](https://github.com/sch-28/ikusa_api)): a powerful tool that allows you to analyze your game logs and gain valuable insights into your combat performance.

## How to start using it
There are two ways to use this scraper for your needs:
* By querying https://bdo.hemlo.cc/communityapi/v1 — this is the "official" instance hosted by me.
* If you want to have more control over the API, host the scraper yourself using one of the following methods:
  - As a Docker container: the image is available on [DockerHub](https://hub.docker.com/r/man90/bdo-rest-api).
  - Natively: build the binary from source as described in [this guide](docs/buildingFromSource.md).

API documentation can be viewed [here](https://man90es.github.io/BDO-REST-API/).

## Flags
If you host the API yourself, either via Docker or natively, you can control some of its features by executing it with flags.

Available flags:
- `-cachecap`
	- Allows to specify response cache capacity
	- Type: unsigned integer
	- Default value: `10000`
- `-cachettl`
	- Allows to specify cache TTL in minutes
	- Type: unsigned integer
	- Default value: `180`
- `-port`
	- Allows to specify API server's port
	- Type: unsigned integer
	- Default value: `8001`
	- Also available as `PORT` environment variable (doesn't work in Docker)
- `-proxy`
	- Allows to specify a list of proxies to make requests to BDO servers through
	- Type: string, space-separated list of IP addresses or URLs
	- Default value: none, requests are made directly
	- Also available as `PROXY` environment variable
- `-verbose`
	- Allows to put the app into verbose mode and print out additional logs to stdout
	- Default value: none, no additional output is produced

You can use them like this:
```bash
./bdo-rest-api -cachettl 30
# or
docker container run -p 8001:8001 bdo-rest-api -cachettl 30
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Known bugs
There is a number of bugs that the official BDO website has. This scraper does not do anything about them for the sake of simplicity, so your apps may need to use [workarounds](docs/brokenStuff.md).

## By the way
This is a fan-created project that is not affiliated with or endorsed by Pearl Abyss.
