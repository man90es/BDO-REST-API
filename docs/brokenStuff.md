# Broken Stuff
BDO website, where all the data is taken from, has a number of bugs, which are replicated by this API by design and, sadly, not much can be done about it.

## List of known bugs
This is a list of bugs that either used to occur or still occur that I am aware of:
1. You request guild's profile by its name and the API returns another guild's profile.
2. You request a guild profile by its name and the API returns a "Not Found" status although the guild does exist.
3. You search for a player or request his or her profile by profileTarget and the API returns no results or "Not Found" status although that player exists. [\[See issue\]](https://github.com/octoman90/BDO-REST-API/issues/5)
4. Members who left the guild remain on the guild members' list for some time.
5. Data is not updated immediately after it is updated in game.

## Workarounds and tips
List numbers match list numbers in the previous section:
1. Always check if the guild profile in response has the same name you specified. See tip #2.
2. You can get some information like creation date, guild master's name and population by searching for guild instead of requesting its profile. It's not much, but better than nothing.
3. There are no known workarounds.
4. I believe maintenances remove "ghost members" from guilds. If you don't feel like waiting, request profiles of those players. Guild membership status in player profiles is more reliable, unless it's set to private.
5. The lag is around a few hours, and you can only wait. This API may introduce additional lag (≤2 hours) in some cases, so if you need the most fresh data possible, consider disabling cache as it's described in [building from source manual](./buildingFromSource.md).

## Contribute to this list
If you found a bug on the original BDO website that affects this API and is not listed in this file, you can contribute by either:
- making a pull request or creating an issue [on GitHub](https://github.com/octoman90/BDO-REST-API)
- using one of the contact methods listed [on my website](https://www.hemlo.cc/)
