# Broken Stuff
BDO website, where all the data is taken from, has a number of bugs, which are replicated by this API by design and, sadly, nothing can be done about it on my side.

## List of known bugs
This is a list of bugs that I am aware of:
1. When profile of one of a number of guilds is requested, another guild's profile is returned. (e.g. requesting the profile of VoS guild on EU returns the profile of Ocrana guild on EU).
1. Opening some guild profiles returns a 404 page despite the fact that that guild exists. (e.g. Valiant guild on EU)
2. Sometimes members who left the guild are still displayed as members for some time.
3. Data is not updated immediately after it is updated in game.

## Workarounds and tips
List numbers match list numbers in the previous section:
1. Check the name of the guild in the response. Is it the same as the one you requested?
1. You can get some information like creation date, guild master's name and population by searching for guild instead of requesting its profile. Not much, but better than nothing.
2. I believe maintenances remove "ghost members" from guilds. If you don't feel like waiting, request profiles of those players. Guild membership status in player profiles is more reliable, unless it's set to private.
3. The lag is around a few hours, and you can only wait. This API may introduce additional lag (≤2 hours) in some cases, so if you need the most fresh data possible, consider disabling cache as it's described in the main README (I assume that you understand the consequences).

## Contribute to this list
If you found a bug on the original BDO website that affects this API and is not listed in this file, you can contribute by either:
- making a pull request or creating an issue [on GitHub](https://github.com/octoman90/BDO-REST-API)
- or sending an email to an address found [in my GitHub profile](https://github.com/octoman90)
- or messaging deadMNGO#8312 on Discord.
