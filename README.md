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

#### Running on GitHub Actions

Set the value of `credentials.json` as a repository secret named `CREDENTIALS`.
Set the value of `twitter_token.json` as a repository secret named `TOKEN`.

> I'm not really a tech savvy but I want to remove the blue checkmark as well and how do I do it?

1. Sign up for a basic twitter Twitter developer account [here](https://developer.twitter.com/en/portal/petition/essential/basic-info). You can ask ChatGPT for the use case by, let's say, exploring Twitter API
2. You will be redirected to the dashboard. On project app, click on the setting button, go to `Keys and tokens` tab, and click the `Regenerate` button. Save the value of the `API Key` and `API Key Secret`
3. Make sure that you already have a GitHub account and logged in. On the top of the page of this repository, click `Fork`
4. You will be redirected to the forked repository. Click on `Actions` tab, and click on `I understand my workflows, go ahead and enable them`
5. On the left pane, there is an item called `Run` with `disabled` text next to it. Click on that and there will be a warning message `This scheduled workflow is disabled because scheduled workflows are disabled by default in forks`. Click on `Enable workflow`
6. Download the binary file from the [release](https://github.com/pr0ph0z/i-hate-blue-checkmark/releases) page with the suitable OS and architecture
7. Double click or run the downloaded binary from the terminal, it will prompt your `API Key` and `API Key Secret` that you obtain from step 2
8. Double click or run the downloaded binary from the terminal again. If it went well, it will open up a browser asking your authorization to an app but if it doesn't, check the terminal and you will see an URL that you need to open
9. After you authorized the app with your Twitter account, you will be redirected to a page indicating successful authentication and a `twitter_token.json` and a `credentials.json` file will be created.
10. Open both file with a text editor, then On the top of the repository, click on `Settings`, then click the `Secrets and variables` drop down, and click on `actions`
11. Click on `New repository secret`, fill the `Name *` with `CREDENTIALS`, and fill the `Secret *` with the content of the `credentials.json` file
12. Repeat the same step with `twitter_token.json` with the name `TOKEN`
13. Go back to `Actions` tab, choose `Run` on the left pane, then click on `Run workflow` to test if the authentication token is can be used