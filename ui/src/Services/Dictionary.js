class Dictionary {
    getTaskTypes() {
        return [
            {"value": "check_work_log", "label": "Проверка списания времени"},
            {"value": "send_team_message", "label": "Отправка сообщения команде"}
        ]
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
}

export default new Dictionary()