# Black Desert social REST API

An unofficial REST API for Black Desert Online that scrapes guild and player data into convenient JSONs.

## Getting the binary
### Building from source
Prerequisites: GNU/Linux, Go >=1.15
```bash
go build
```

### Prebuilt
You can download prebuild binaries [here](https://gitlab.com/man90/black-desert-social-rest-api/-/pipelines).

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
| query          | No.                                |                                             |
| searchType     | Only if «query» param is provided. | Supported values: familyName, characterName |
| page           | No.                                | Each page has up to 20 players.             |

You can find an example of a reply [here](https://gitlab.com/man90/black-desert-social-rest-api/-/blob/master/exampleDumps/profileSearch.json).

## To-do
* Cache
* Error handling
