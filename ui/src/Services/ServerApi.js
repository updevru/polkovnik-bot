class ServerApi {

    /**
     * @param restUrl
     */
    constructor(restUrl) {
        this.restUrl = restUrl;
    }

    team() {
        return new ApiResource(this.restUrl, 'team');
    }

    teamSettings(teamId) {
        return new ApiResource(this.restUrl, 'team/' + teamId + '/settings');
    }

    user(teamId) {
        return new ApiResource(this.restUrl, 'team/' + teamId + '/users');
    }

    task(teamId) {
        return new ApiResource(this.restUrl, 'team/' + teamId + '/tasks');
    }

    receiver(teamId) {
        return new ApiResource(this.restUrl, 'team/' + teamId + '/receivers');
    }

    async runTask(teamId, taskId) {
        return fetch(
            this.restUrl + '/team/' + teamId + '/tasks/' + taskId + '/run',
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
            }
        ).then(response => response.json());
    }

    taskHistory(teamId, taskId) {
        return new ApiResource(this.restUrl, 'team/' + teamId + '/tasks/' + taskId + '/history');
    }

    async sendTeamMessage(teamId, text) {
        return fetch(
            this.restUrl + '/team/' + teamId + '/sendMessage',
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
                body: JSON.stringify({"Text": text})
            }
        ).then(response => response.json());
    }


    /**
     * Форматирует дату для API в формат 2020-04-16T05:26:33
     * @returns string|null
     * @param date
     */
    dateFormat(date) {
        if (date instanceof Date) {
            return date.getFullYear() + '-' + date.getMonth() + '-' + date.getDate() +
                'T' + date.getHours() + ':' + date.getMinutes() + ':' + date.getSeconds();
        }
        return null;
    }
}

class ApiResource {

    constructor(restUrl, name) {
        this.restUrl = restUrl;
        this.name = name;
    }

    /**
     * Возвращает уникальный ID объекта
     * @param id
     * @returns {string}
     */
    id(id) {
        let result = '/' + this.name;
        if (id) {
            result += '/' + id;
        }
        return result;
    }

    /**
     * Возвращает полный URL до ресурса
     * @param id
     * @param query
     * @returns {URL}
     */
    url(id, query) {
        let url = new URL(this.restUrl + this.id(id));

        if (query) {
            url.search = new URLSearchParams(query).toString();
        }

        return url;
    }

    /**
     * Получение одного объекта по его ID
     * @param id
     * @returns {Promise<any>}
     */
    async get(id) {
        return fetch(this.url(id)).then(response => response.json());
    }

    /**
     * Получение списка объектов
     * @param parameters
     * @param page
     * @returns {Promise<any>}
     */
    async list(parameters, page) {
        parameters = parameters || {};
        if (page) {
            parameters['page'] = page;
        }

        return fetch(this.url('', parameters)).then(response => response.json());
    }

    /**
     * Создание нового объекта
     * @param parameters
     * @returns {Promise<any>}
     */
    async add(parameters) {
        return fetch(
            this.url(''),
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
                body: JSON.stringify(parameters)
            }
        ).then(response => response.json());
    }

    /**
     * Обновление существующего объекта
     * @param id
     * @param parameters
     * @returns {Promise<any>}
     */
    async edit(id, parameters) {
        return fetch(
            this.url(id),
            {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/merge-patch+json;charset=utf-8'
                },
                body: JSON.stringify(parameters)
            }
        ).then(response => response.json());
    }

    /**
     * Удаление объекта по его ID
     * @param id
     * @returns {Promise<any>}
     */
    async delete(id) {
        return fetch(this.url(id), {method: 'DELETE'}).then(response => response.json());
    }
}

export default new ServerApi(document.location.origin + '/api');