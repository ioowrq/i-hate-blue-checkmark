# I Hate Blue Checkmark
## Overview
Update your Twitter profile picture regularly to remove the blue checkmark. This tool uses your current profile picture, so no new image will be uploaded.

## Prerequisites
- You must have an active Twitter Developer account and an app registered on developer.twitter.com to obtain the `API_KEY` and `API_SECRET`.

## Usage
### Authenticating Your Twitter Account
Download the binary from the release page, build your own binary file, or run directly with `go run main.go`.
The program will prompt you to enter your API Key and API Secret.
You will be redirected to the Twitter OAuth page. Authorize the app, and you will be redirected to a page indicating successful authentication.
A `twitter_token.json` and a `credentials.json` file will be created.

### Updating Your Profile Image
#### Running Locally

Whenever you run the program again, it will check for the existence of `credentials.json` and `twitter_token`.json files.

### Running on GitHub Actions

Set the value of `credentials.json` as a repository secret named `CREDENTIALS`.
Set the value of `twitter_token.json` as a repository secret named `TOKEN`.