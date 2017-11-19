# pomodoro-thing
### Pomodoro Slack Bot Experiment

This current version can track pomodoros a basic way, works as a slash
command in Slack. It prevents to run multiple pomodoro sessions for
a single user and sessions could have task description.


### Usage
    go run app.go

Then expose the service using ngrok.io - that's it for now.

### Todo
* need a way to have a specific user token for every user to maintain Slack status
* dockerize
