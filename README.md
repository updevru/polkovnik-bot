# PolkovnikBot

Бот-помощник по огранизации работы команды и увелияения ее эффективности.
Позволяет выполнять определенные действия по расписанию.

![PolkovnikBot Screenshot](/docs/images/screen.PNG)

## Возможности

- Напоминание о списании времени по задачам
- Отправка сообщений в командный чат по расписанию
- Интеграция с такс трекерами - Jira
- Интеграция с чатами - Telegram
- Управление через web интерфейс

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
docker run updev/polkovnik-bot -v ./config.json:/app/var/config.json -p 8080:8080
```

## Запуск

```bash
./PolkovnikBot.exe -o
```

Параметры запуска:
```
-c string Config file (default "./var/config.json")
-o Send logs to stdout
-p HTTP port for UI (default 8080)
-u Folder with UI (default "./ui/build")
```
