# PolkovnikBot

Бот-помощник по огранизации работы команды и увелияения ее эффективности.
Позволяет выполнять определенные действия по расписанию.

## Возможности

- Напоминание о списании времени по задачам
- Отправка сообщений в командный чат по расписанию
- Интеграция с такс трекерами - Jira
- Интеграция с чатами - Telegram
- Управление через web интерфейс

## Установка

Создать файл конфигурации на основе примера var/config.sample.json и сохранить его в var/config.sample.json

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

## Описание конфигурации

```json
{
  "Teams": [ //Список команд
    {
      "Title": "Team name", //Название команды
      "Users": [ //Состав команды
        {
          "Name": "Alex White", //Имя члена команды
          "Login": "alex", //Логин в системе задач (jira)
          "NickName": "alex" //Ник в чате
        },
        {
          "Name": "Den Black",
          "Login": "den",
          "NickName": "den_black"
        }
      ],
      "Tasks": [ //Список заданий для бота
        {
          "Schedule": { //Расписание задания проверки списания задач в конце рабочего дня
            "WeekDays": [ //Дни по которым отрабатывает задание
              "Monday",
              "Tuesday",
              "Wednesday",
              "Thursday",
              "Friday"
            ],
            "Hour": 18, //Час в которое отработает задание
            "Minute": 0 //Минута
          },
          "Type": "check_work_log", //Тип задания (проверка списанного времени)
          "Projects": [ //Проекты в которых учитывать списанное время
            "DEV"
          ]
        },
        {
          "Schedule": { //Расписание задания проверки списания времени в понедельник утром за пятницу
            "WeekDays": [
              "Monday"
            ],
            "Hour": 11,
            "Minute": 0
          },
          "Type": "check_work_log",
          "Projects": [
            "DEV"
          ],
          "DateModify": "-72h"
        },
        {
          "Schedule": { //Расписание задания
            "WeekDays": [ //Дни по которым отрабатывает задание
              "Monday",
              "Tuesday",
              "Wednesday",
              "Thursday",
              "Friday"
            ],
            "Hour": 10, //Час в которое отработает задание
            "Minute": 0 //Минута
          },
          "Type": "send_team_message", //Тип задания (отправка сообщения)
          "Message": "It's time to meet" //Текст сообщения, который будет отправлен
        }
      ],
      "Channel": { //Настроки канала для отправки уведомлений
        "Type": "telegram", //Тип канала
        "ChannelId": "-1001145000000", //ID канала
        "Settings": {
          "token": "331640000:AAEcl3yHv...." //Токен
        }
      },
      "Weekend": { //Настройка общего расписания всей команды
        "WeekDays": null,
        "Intervals": null
      },
      "IssueTracker": { //Настройка таск-трекера команды
        "Type": "jira", //Тип трекера
        "Settings": {
          "password": "d0kdsh89KR69K", //Пароль
          "url": "https://jira.domain.com", //Адрес трекера
          "username": "bot" //Пользователь
        }
      },
      "MinWorkLog": 20000 //Минимальное количество времени которое нужно списать в день
    }
  ]
}
```