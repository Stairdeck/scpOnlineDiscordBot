# scpOnlineDiscordBot
scpDiscordBot is a application that allow you to create Discord bots, which displaying current
online of your servers in status bar. 

The latest release can be found here: [Release link](https://github.com/Stairdeck/scpOnlineDiscordBot/releases/latest)

## Installation
Just place all files in any directory, edit config.json and run.

Release has two versions built: for linux and windows.

But you can download sources and compilation for system you need.

## Configuration
Example config.json
```json
{
  "updateTime": 15,
  "logger": true,
  "servers": [
    {
      "name": "Server1",
      "ip": "255.255.255.255",
      "port": 7777,
      "serverId": 25655,
      "discordUserBotToken": "TokenForBot1"
    },
    {
      "name": "Server2",
      "ip": "255.255.255.255",
      "port": 7778,
      "serverId": 25613,
      "discordUserBotToken": "TokenForBot2"
    }
  ]
}
```

You need to fill the serverId or ip and port for any server!

| Option | Type | Description |
| ------ | ------ | ------ |
| updateTime | int | Time in seconds after information is updated. Recommended > 15 |
| updateTime | bool | Additional log information in console |
| servers | array | Servers array with detail information of server |
| name | string | Name of your server, it is for logging |
| ip | string | IP of your server |
| port | int | Port of your server |
| serverId | int | ID of your server in [SCP:SL ServerList](https://api.scpslgame.com/lobbylist.php?format=json) |
| discordUserBotToken | string | Bot's token |

## TODO

- emojis in Discord's bots statuses

## Credits
Thanks a lot to @andersfylling for creating [disgord](https://github.com/andersfylling/disgord)

##
Please report any problems or suggestions in the issues tab at the top.

Thanks & Enjoy.