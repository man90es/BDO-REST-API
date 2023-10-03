# BDO-REST-API
[![license](https://img.shields.io/github/license/man90es/BDO-REST-API)](https://github.com/man90es/BDO-REST-API/blob/master/LICENSE)

A scraper for Black Desert Online player in-game data with a REST API. It currently supports European, North American, and South American servers (Korean server support is in progress).

## Projects using this API
- [BDO Leaderboards](https://bdo.hemlo.cc/leaderboards) ([Source](https://github.com/man90es/BDO-Leaderboards)): web-based leaderboards for Black Desert Online.
- [Ikusa](https://ikusa.site) ([Source](https://github.com/sch-28/ikusa_api)): a powerful tool that allows you to analyze your game logs and gain valuable insights into your combat performance.

## How to start using it
There are two ways to use this scraper for your needs:
* By querying https://bdo.hemlo.cc/communityapi/v1 — this is the "official" instance hosted by me.
* If you want to have more control over the API, host the scraper yourself. It's [available on DockerHub](https://hub.docker.com/r/man90/bdo-rest-api) (exposes port 8001), preconfigured for Heroku, and can be built from source as described in [this guide](docs/buildingFromSource.md) (this gives you a bit more control over how the scraper behaves). You can also build Docker image on your machine like it's explained [here](docs/buildingDocker.md).

API documentation can be viewed [here](https://man90es.github.io/BDO-REST-API/).

## Known bugs
There is a number of bugs that the official BDO website has. This scraper does not do anything about them for the sake of simplicity, so your apps may need to use [workarounds](docs/brokenStuff.md).

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## By the way
This is a fan-created project that is not affiliated with or endorsed by Pearl Abyss.
