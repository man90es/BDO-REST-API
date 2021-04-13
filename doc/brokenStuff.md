# Broken Stuff
BDO website, where all the data is taken from, has a number of bugs, which are replicated by this API by design and, sadly, nothing can be done about it on my side.

## List of known bugs
This is a list of bugs that I am aware of:
1. Opening some guild profiles returns profiles of different guilds. (Example: VoS guild on EU)
1. Opening some guild profiles returns a 404 page. (Example: Valiant guild on EU)
1. Sometimes members who left the guild are still displayed as members.
1. Data is not updated immediately after it is updated in game.

## Workarounds and tips
List numbers match list numbers in the previous section:
1. Use the guild search to find out who the real guild master is. Is his name the same as the one in the guild profile? If yes, you most probably got the right profile.
1. You can get some information like creation date, guild master's name and population by searching for guild instead of requesting its profile. Not much, but better than nothing.
1. I believe maintenances remove "ghost members" from guilds. If you don't feel like waiting, request profiles of those players and check if they are:
	1. Still in the guild
	1. Members of another guild or not in a guild
	1. Have their guild set as private, so you can't check.
1. The lag is around a few hours, and you can only wait. This API may introduce additional lag (≤2 hours) in some cases, so if you need the most fresh data possible, consider disabling cache (I assume you understand the consequences).

## Contribute to this list
If you found a bug on the original BDO website that affects this API and is not listed in this file, you can contribute by:
1. Making a pull request on [GitLab](https://gitlab.com/man90/black-desert-social-rest-api) or [GitHub](https://github.com/octoman90/BDO-REST-API)
2. Sending an email to an address found [here on GitLab](https://gitlab.com/man90) or [here on GitHub](https://github.com/octoman90)
3. Messaging deadMNGO#8312 on Discord.
