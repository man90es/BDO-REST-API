# Broken Stuff
BDO website, where all the data is taken from, has a number of bugs, which affects by this API by design and, sadly, not much can be done about it.

## List of known bugs
This is a list of bugs that either used to occur that I'm aware of. They may be still occuring or be fixed.

### 🐞 Data is not updated immediately after it is updated in game
The website's lag is around a few hours, and you can only wait. This API introduces additional lag that depends on the cache TTL parameter.

### 🐞 Members who left the guild remain on the guild members' list for some time
I believe maintenances remove "ghost members" from guilds. If you don't feel like waiting, request profiles of those players. Guild membership status in player profiles is more reliable, unless it's set to private.

### 🐞 Profile of a different guild is returned instead of the one requested
Always check if the guild profile in response has the same name you specified. Also see next bug.

### 🐞 Guild profile returned as "Not Found" although the guild exists in the game
You can get some information like creation date, guild master's name and population by searching for guild instead of requesting its profile. It's not much, but better than nothing.

### 🐞 Player profile returned as "Not Found" although that player exists in the game
There are no known workarounds. [See issue #5](https://github.com/man90es/BDO-REST-API/issues/5).

### 🐞 Players whose family name is longer than 16 characters aren't searchable
Family names longer than 16 characters aren't officially allowed in BDO, but there may be a bug that allows players to take them. [See issue #13](https://github.com/man90es/BDO-REST-API/issues/13). This API won't support them unless there will be many reports of players with long family names.

### 🐞 Some guild members are missing when requesting a guild profile
BDO website hides some guild members on its website based on unknown conditions. Supposedly, all hidden members are alt accounts with zero progress on the account.

### 🐞 The number of guild members when searching for a guild doesn't match the number of guild members when requesting the guild profile
See the previous bug on the list. The number when searching is the correct one.

## Contribute to this list
If you found a bug on the original BDO website that affects this API and is not listed in this file, you can contribute by either:
- making a pull request or creating an issue [on GitHub](https://github.com/man90es/BDO-REST-API)
- using one of the contact methods listed [on my website](https://www.hemlo.cc/)
