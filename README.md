# This software is DEPRECATED
Get better version here [DiscordMBM](https://github.com/Stairdeck/discordMBM)

# scpOnlineDiscordBot
scpOnlineDiscordBot is a application that allow you to create Discord bots, which displaying current
online of your servers in status bar. 

The latest release can be found here: [Release link](https://github.com/Stairdeck/scpOnlineDiscordBot/releases/latest)

![example_bots](https://stairdeck.com/bots_en.png)

## Installation
Just place all files in any directory, edit config.json and run.

Release has two versions built: for linux and windows.

But you can download sources and make a compilation for the system you need.

## Configuration
Example config.json
```json
{
  "logger": true,
  "accountId": 1234,
  "apiKey": "XXXXXXXX",
  "servers": [
    {
      "name": "Server1",
      "serverId": 1234,
      "discordUserBotToken": "XXXXXXXXXXXXX"
    },
    {
      "name": "Server2",
      "serverId": 1234,
      "discordUserBotToken": "XXXXXXXXXXXX"
    }
  ]
}
```

You need to fill the serverId for any server!

| Option | Type | Description |
| ------ | ------ | ------ |
| logger | bool | Additional log information in console |
| accountId | int | Your account id. You can find it here [ServerList](https://servers.scpslgame.com/) (Click on your server to expand)
| apiKey | string | Your api key. Type `!api` in your scp:sl server console |
| servers | array | Servers array with detail information of server |
| name | string | Name of your server, it is for logging |
| serverId | int | ID of your server in [SCP:SL ServerList](https://api.scpslgame.com/lobbylist.php?format=json) |
| discordUserBotToken | string | Bot's token |

## TODO

- emojis in Discord's bots statuses

## Credits
Thanks a lot to @andersfylling for creating [disgord](https://github.com/andersfylling/disgord)

##
Please report any problems or suggestions in the issues tab at the top.

Thanks & Enjoy.
