# Broken Stuff
BDO website, where all the data is taken from, has a number of bugs, which are replicated by this API by design and, sadly, not much can be done about it.

## List of known bugs
This is a list of bugs that either used to occur or still occur that I am aware of:
1. When you request one guild, the server may return another guild.
2. When you request a guild, the server may return a "not found" status despite the fact that that guild does exist.
3. Members who left the guild are still displayed as members for some time.
4. Data is not updated immediately after it is updated in game.

## Workarounds and tips
List numbers match list numbers in the previous section:
1. If the response has a different guild name than the one you requested, it's a bug. See tip #2.
2. You can get some information like creation date, guild master's name and population by searching for guild instead of requesting its profile. Not much, but it's better than nothing.
3. I believe maintenances remove "ghost members" from guilds. If you don't feel like waiting, request profiles of those players. Guild membership status in player profiles is more reliable, unless it's set to private.
4. The lag is around a few hours, and you can only wait. This API may introduce additional lag (≤2 hours) in some cases, so if you need the most fresh data possible, consider disabling cache as it's described in [building from source manual](./buildingFromSource.md).

## Contribute to this list
If you found a bug on the original BDO website that affects this API and is not listed in this file, you can contribute by either:
- making a pull request or creating an issue [on GitHub](https://github.com/octoman90/BDO-REST-API)
- or sending an email to an address found [in my GitHub profile](https://github.com/octoman90)
- or messaging mvngo#8312 on Discord.
