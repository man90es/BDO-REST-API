# Black Desert social REST API

An unofficial REST API for Black Desert Online that scrapes guild and player data into convenient JSONs.

## Projects using this API
- [BDO Leaderboards](https://man90.gitlab.io/bdo-leader-boards)

## Getting the binary
### Prebuilt
You can download prebuilt binaries [here](https://gitlab.com/man90/black-desert-social-rest-api/-/pipelines).

### Building from source
Prerequisites: <abbr title="Not tested on other platforms.">GNU/Linux</abbr>, Go >=1.15

Command:
```bash
go build
```

By default, scraped results are cached in memory and stored for up to 2 hours, it helps to ease the pressure on BDO servers and speeds up the response time in some situations. If you don't want to cache scraped results (for example if you just want to create a dump or need extra fresh data), use this command instead:
```bash
go build -tags "cacheless"
```

## Environment variables
Catch requests on a specific port (8001 by default):
```bash
export PORT=3000
```

Use a proxy to make requests to BDO servers (none by default):
```bash
export PROXY=http://123.123.123.123:8080
# or
export PROXY="http://123.123.123.123:8080 http://124.124.124.124:8081"
```

## Routes

### Guild data
`http://localhost:8001/v0/guildProfile`

| GET Parameter | Required | Wtf?                     |
|-----------|----------|--------------------------|
| guildName | Yes.     | The name of the guild.   |
| region    | Yes.     | Supported values: EU, NA |

You can find an example of a reply [here](https://gitlab.com/man90/black-desert-social-rest-api/-/blob/master/exampleDumps/guildProfile.json).

### Player data
`http://localhost:8001/v0/profile`
| GET Parameter     | Required | Wtf?                                                                  |
|---------------|----------|-----------------------------------------------------------------------|
| profileTarget | Yes.     | You can get this string from guild members data or from player search |

You can find an example of a reply [here](https://gitlab.com/man90/black-desert-social-rest-api/-/blob/master/exampleDumps/profile.json).

### Guild search
`http://localhost:8001/v0/guildProfileSearch`
| GET Parameter | Required | Wtf?                           |
|---------------|----------|--------------------------------|
| region        | Yes.     | Supported values: EU, NA       |
| query         | No.      |                                |
| page          | No.      | Each page has up to 10 guilds. |

You can find an example of a reply [here](https://gitlab.com/man90/black-desert-social-rest-api/-/blob/master/exampleDumps/guildProfileSearch.json).

### Player search
`http://localhost:8001/v0/profileSearch`
| GET Parameter  | Required                           | Wtf?                                        |
|----------------|------------------------------------|---------------------------------------------|
| region         | Yes.                               | Supported values: EU, NA                    |
| query          | Yes.                               |                                             |
| searchType     | Yes.                               | Supported values: familyName, characterName |
| page           | No.                                | Each page has up to 20 players.             |

You can find an example of a reply [here](https://gitlab.com/man90/black-desert-social-rest-api/-/blob/master/exampleDumps/profileSearch.json).

## Known bugs
A vast majority of bugs comes from the original BDO website, where data is taken from. You can find a list of known bugs and workarounds [here](doc/brokenStuff.md).

## To-do
* Error handling
* Provide info about which character is player's main character
* Provide meta info through the API routes?
