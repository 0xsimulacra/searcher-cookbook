# log information and update your program via discord

Here I show you an example how you can use a discord bot to update configurations without restarting your program and also log some important information to your private discord channel

### Important

You should never log or update sensible informations over discord, even if it's your private channel. So don't think about updating your private key or something like that over this method.  
You could also add another security check and also only consume messages from a defined user (you). 

## Pre-condition
- A discord bot on your server, check out [this article](https://discord.com/developers/docs/getting-started#creating-an-app) until installing your app.
- The bot token
- The channel id

## Usage
Start
```bash
go run discord.go --token xxx --channel yyy 
```

## Discord commands
Print your current configuration
```
!show
bot-message: {"debug":false,"version":"1.0.0","amount":10000000}
```

Update a single key from you config - this will now send a message every 10 seconds to your channel
```
!config.debug true
bot-message: updated Debug to true
bot-message (every 10 second): debugging active !

!config.amount 420
bot-message: updated Amount to 420

!show
bot-message: {"debug":true,"version":"1.0.0","amount":420}
```

Update the whole config - you can change any value in config.json for this
```
!reload
bot-message: successfully reloaded config file!

!show
bot-message: {"debug":false,"version":"1.0.0","amount":10000000}
```