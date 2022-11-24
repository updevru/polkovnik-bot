# PolkovnikBot

[![Go](https://github.com/updevru/polkovnik-bot/workflows/Go/badge.svg)](https://github.com/updevru/polkovnik-bot/actions)

Бот-помощник по организации работы команды и увеличения ее эффективности.
Позволяет выполнять определенные действия по расписанию.

![PolkovnikBot Screenshot](/docs/images/screen.PNG)

## Возможности

- Напоминание о списании времени по задачам
- Уведомление об отпусках
- Отправка сообщений в командный чат по расписанию
- Интеграция с такс трекерами - Jira
- Интеграция с чатами - Telegram, Webex
- Управление через web интерфейс
- API

## Установка

**Сборка из исходниго кода:**

```bash
git clone https://github.com/updevru/polkovnik-bot.git
cd polkovnik-bot/
go build

cd ui
npm install
npm run build

./polkovnik -c ./config.json
```

**Запуск в контейнере Docker:**

```bash
docker run updev/polkovnik-bot -v ./config.json:/app/var/config.json -v ./data.db:/app/var/data.db -p 8080:8080
```

## Запуск

```bash
./PolkovnikBot.exe -o
```

Параметры запуска:
```
-c string Config file (default "./var/config.json")
-db string Database file (default "./var/data.db")
-o Send logs to stdout
-p HTTP port for UI (default 8080)
```

## Документация API

После запуска зайти по адресу /doc/api/

## Разработка

Для запуска проекта в режиме разработке удобно использовать Docker.

В корне репоизтория выполнить.
```bash
docker compose up
```

Открыть в браузере http://localhost:3000