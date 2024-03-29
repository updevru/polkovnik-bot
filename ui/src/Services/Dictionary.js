class Dictionary {
    getTaskTypes() {
        return [
            {"value": "check_work_log", "label": "Проверка списания времени"},
            {"value": "check_work_log_by_period", "label": "Проверка списания времени за период"},
            {"value": "send_team_message", "label": "Отправка сообщения команде"},
            {"value": "check_user_weekend", "label": "Уведомление об отпусках"}
        ]
    }

    getTaskType(value) {
        let result
        this.getTaskTypes().forEach(function (item) {
            if (item.value === value) {
                result = item
                return
            }
        })

        return result;
    }

    getWeekdays() {
        return [
            { label: 'Monday', value: 'Monday' },
            { label: 'Tuesday', value: 'Tuesday' },
            { label: 'Wednesday', value: 'Wednesday' },
            { label: 'Thursday', value: 'Thursday' },
            { label: 'Friday', value: 'Friday' },
            { label: 'Saturday', value: 'Saturday' },
            { label: 'Sunday', value: 'Sunday' },
        ]
    }

    getTrackerTypes() {
        return [
            {"label": "Jira", "value": "jira"},
        ]
    }

    getChannelTypes() {
        return [
            {"label": "Telegram", "value": "telegram"},
            {"label": "Webex", "value": "webex"},
        ]
    }

    getReceiverTypes() {
        return [
            {"label": "Отправка сообщения команде", "value": "send_team_message"},
        ]
    }

    getReceiverType(value) {
        let result
        this.getReceiverTypes().forEach(function (item) {
            if (item.value === value) {
                result = item
                return
            }
        })

        return result;
    }

    getReceiverFormats() {
        return [
            {"label": "Автоопределение", "value": "auto"},
            {"label": "JSON", "value": "json"},
            {"label": "XML", "value": "xml"},
            {"label": "Текст", "value": "text"},
        ]
    }
}

export default new Dictionary()