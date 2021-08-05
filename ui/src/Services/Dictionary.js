class Dictionary {
    getTaskTypes() {
        return [
            {"value": "check_work_log", "label": "Проверка списания времени"},
            {"value": "send_team_message", "label": "Отправка сообщения команде"}
        ]
    }

    getTaskType(value) {
        let result
        this.getTaskTypes().forEach(function (item) {
            if (item.value === value) {
                result = item
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
}

export default new Dictionary()