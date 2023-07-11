# BDO-REST-API
[![license](https://img.shields.io/github/license/octoman90/BDO-REST-API)](https://github.com/octoman90/BDO-REST-API/blob/master/LICENSE)

A scraper for Black Desert Online player in-game data with a REST API. It currently supports European, North American, and South American servers (Korean server support is in progress).

## Projects using this API
- [BDO Leaderboards](https://bdo.hemlo.cc/leaderboards) ([Source](https://github.com/octoman90/BDO-Leaderboards)): web-based leaderboards for Black Desert Online.
- [Ikusa](https://ikusa.site) ([Source](https://github.com/sch-28/ikusa_api)): a powerful tool that allows you to analyze your game logs and gain valuable insights into your combat performance.

## How to start using it
There are two ways to use this API in your apps:
* https://bdo.hemlo.cc/communityapi/v1 is the "official" instance hosted by me. Using it is easy but it may be slow or even become unresponsive. The API documentation can be viewed [here](https://gitlab.com/man90/black-desert-social-rest-api/-/tree/master/doc/api/openapi.json).
* If you want to have more control over the API, host the scraper yourself. The repository is preconfigured to be deployable on Heroku, can be easily deployed with Docker or built manually for a VPS/VDS. To do this, follow these steps:
	1. Build the server from the source code following [this guide](doc/buildingFromSource.md).
	2. (Optional) Set the environment variables. The list of available variables is in the section below.
	3. Run the binary. Possible flags are described in a section below.
	4. Use the API as described in the [documentation](https://gitlab.com/man90/black-desert-social-rest-api/-/tree/master/doc/api/openapi.json).

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

## Known bugs
There is a number of bugs that the official BDO website has. This collector does not do anything about them for the sake of simplicity, so your apps may need to use the [workarounds](doc/brokenStuff.md).

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## By the way
This is a fan-created project that is not affiliated with or endorsed by Pearl Abyss.
