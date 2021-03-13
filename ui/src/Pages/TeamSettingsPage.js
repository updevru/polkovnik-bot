import React from 'react';
import {PageHeader, Form, Input, Button, InputNumber, Select, Divider } from "antd";

const trackerTypes = [
    {"title": "Jira", "value": "jira"},
];

const channelTypes = [
    {"title": "Telegram", "value": "telegram"},
];

class TeamSettingsPage extends React.Component{

    componentDidMount() {
        console.log("SettingsPage mount")
    }

    onFinish(values: any) {
        console.log('Success:', values);
    };

    onFinishFailed(errorInfo: any) {
        console.log('Failed:', errorInfo);
    };

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Настройки"
            >
                <div className={"app-form-container"}>
                <Form
                    layout="vertical"
                    name="basic"
                    initialValues={{ title: "BIS", minWorklog: 14400 }}
                    onFinish={this.onFinish}
                    onFinishFailed={this.onFinishFailed}
                >
                    <Form.Item
                        label="Название команды"
                        name="title"
                        rules={[{ required: true, message: 'Название команды не должно быть пустым' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item
                        label="Минимальное время которое нужно списать по задачам (сек.)"
                        name="minWorklog"
                        rules={[{ required: true, message: 'Название команды не должно быть пустым' }]}
                    >
                        <InputNumber  />
                    </Form.Item>

                    <Divider>Система ведения задач</Divider>

                    <Form.Item
                        label="Тип"
                        name="IssueTracker_type"
                        rules={[{ required: true}]}
                    >
                        <Select options={trackerTypes} />
                    </Form.Item>

                    <Form.Item
                        label="Адрес"
                        name="IssueTracker_url"
                        rules={[{ required: true, message: 'Адрес не должен быть пустым' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item
                        label="Пользователь"
                        name="IssueTracker_username"
                        rules={[{ required: true, message: 'Имя пользователя не должено быть пустым' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item
                        label="Пароль"
                        name="IssueTracker_password"
                        rules={[{ required: true, message: 'Пароль не должен быть пустым' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Divider>Канал оповещения</Divider>

                    <Form.Item
                        label="Тип"
                        name="Channel_type"
                        rules={[{ required: true}]}
                    >
                        <Select options={channelTypes} />
                    </Form.Item>

                    <Form.Item
                        label="ID канала"
                        name="Channel_id"
                        rules={[{ required: true, message: 'ID канала не должен быть пустым' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item
                        label="Токен бота"
                        name="Channel_bot_token"
                        rules={[{ required: true, message: 'Имя пользователя не должено быть пустым' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit">
                            Сохранить
                        </Button>
                    </Form.Item>
                </Form>
                </div>
            </PageHeader>
        )
    }
}

export default TeamSettingsPage