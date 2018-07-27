# goproj
create different types of projects in golang

Get the Package:
```go get -u github.com/aliforever/goproj```

Install:
```go install```

- Create a Telegram bot using the library:
```goproj --type=bot --username=BOT_USERNAME --token=BOT_TOKEN```

- Add New Menu to Bot's Engine File:
```goproj --username=BOT_USERNAME --make=menu:MENU_NAME:LINE_NUMBER
(leave the username part to add menu to existing project)
