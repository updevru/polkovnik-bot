import React from 'react';
import {Typography, Table} from "antd";

const { Title } = Typography;

class ReceiverHelpPage extends React.Component{

    renderObjectTable() {
        let columns = [
            {
                title: 'Параметр',
                dataIndex: 'param',
                key: 'param'
            },
            {
                title: 'Описание',
                dataIndex: 'text',
                key: 'text'
            }
        ];

        let data = [
            {
                param: 'Method',
                text: 'Каким методом был отправлен запрос - GET, POST, PUT, DELETE и т.д.'
            },
            {
                param: 'Body',
                text: 'Преобразованное тело запроса в объект со всеми свойствами и значениями. Поддерживаются форматы: JSON, XML. Если формат не удалось определить, то тут будет просто текст.'
            },
            {
                param: 'Params',
                text: 'Массив параметров из URL запроса'
            },
            {
                param: 'Headers',
                text: 'Массив заголовков из запроса'
            }
        ];

        return (
            <Table columns={columns} dataSource={data} pagination={false} />
        );
    }

    renderFunctionTable() {
        let columns = [
            {
                title: 'Название',
                dataIndex: 'name',
                key: 'name'
            },
            {
                title: 'Описание',
                dataIndex: 'text',
                key: 'text'
            },
            {
                title: 'Пример',
                dataIndex: 'example',
                key: 'example'
            }
        ];

        let data = [
            {
                name: 'getValue',
                text: 'Выводит элемент в массиве если он там есть, если нет, то ничего не выводит.',
                example: '{{  getValue "test" .Params }} Отображает свойство "test" из массива .Params'
            }
        ];

        return (
            <Table columns={columns} dataSource={data} pagination={false} />
        );
    }

    render() {
        const postData = JSON.stringify({"id":"13bc8e56-f754-40ab-ab16-a21ec4444b72","name":"Nick", "lastname":"Brown", "login":"test"});
        const postTemplate  = "{{ .Method }} запрос: ->  пользователь {{ .Body.name }} с логином {{ .Body.login }} сменил имя на {{  getValue \"newName\" .Params }}";
        const postResult  = "POST запрос: -> пользователь Nick с логином test сменил имя на Jon";

        return (
            <>
            <Title level={3}>Формирование шаблона</Title>
            <p>
                В качестве движка преобразования данных запроса в текст используется <a href={"https://golangforall.com/en/post/templates.html"} target={"_blank"} rel={"noreferrer"}>шаблоны языка Go</a>.
                В шаблонах можно использовать, переменные, цыклы, условия, функции и т.д. Благодаря такому функционалу можно достаточно гибко трансформировать запрос в текст.
            </p>
            <p>В момент получения запроса система его разбирает, формирует объект и передает в шаблон.</p>
            <p>{this.renderObjectTable()}</p>
            <Title level={4}>Функции</Title>
            <p>{this.renderFunctionTable()}</p>
            <Title level={4}>Пример</Title>
            <p>
                <b>Запрос:</b><br/>
                POST http://host/receive/id?newName=Jon&old=25<br/>
                {postData}
            </p>
            <p>
                <b>Шаблон:</b><br/>
                {postTemplate}
            </p>
            <p>
                <b>Результат:</b><br/>
                {postResult}
            </p>
            </>
        )
    }
}

export default ReceiverHelpPage