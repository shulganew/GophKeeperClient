# GophKeeperClient
Password keeper - client (Yandex praktikum final project)
Server:
https://github.com/shulganew/GophKeeper

Client:
https://github.com/shulganew/GophKeeperClient

## Build Client
// set flag
git tag -a v1.0.0 -m "release app. First stable."

// makefile crossplatform build
```
make build_win build_linux build_mac
```
valid combination GOOS  GOARCH 
```
go tool dist list
```

## Generate oapi
Use make or bash command or //TODO build generate
```
make oapi


## Current questions
Привет!
Пытаюсь делать проект по частям, сделал базовую регистрацию и логирование для клиента и севера. Посмотри, пожалуйста.
После проверки запишусь к тебе на 1:1, хочу обсудить вопросы:

дальнейший план работы по проекту
нужен ли в этом проекте gRPC
поговорить о шифровании, хочу рассказать, как я понял и понять, правильно ли))
о синхронизации между клиентами
Клиент живет в отдельном репозитории:
User register GophKeeperClient#1

# TUI
## tcell and tview exapmles (not used, complicated)
https://github.com/rivo/tview/tree/master

## bubbletea TUI framework
https://github.com/charmbracelet/bubbletea

##