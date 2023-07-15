# TeleHook

A simple programm ment to create discord-like webhooks for telegram

## Usage

Send post requests to [your ip]/[secret path]/[bot token]

The secret path is hardcoded and `Q75k9anIncOQO9peWkF0HMTkIyjQVsSd` by default.
You can get the bottoken from @Botfather on telegram after creating a bot

In `example.json` you can find an example json in the format of the json you're supposed to give the webhook when making the post request

## Errors

`TH_WEB_NO_TOKEN` You forgot to append your bots token to the url

`TH_WEB_JSON_FAIL` The Json you gave the webhook is malformed 

`TH_WEB_IMG_FAIL` The Image can only be base64. It's either not base64 or malformed

`Not Found` The token you gave the webhook is invalid
