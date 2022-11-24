import React from 'react';
import {Button, Divider, Form, Input, InputNumber, Select, Row, Col} from "antd";
import Dictionary from "../Services/Dictionary";
import {Weekend, WeekendDataValue, WeekendFormValue} from "../Components/Weekend/Weekend";

class TeamSettingsForm extends React.Component {

    constructor(props) {
        super(props);
        this.onFinish = this.onFinish.bind(this)
    }

    onFinish(values: any) {
        let data = WeekendDataValue(values);

        if (data.issue_tracker_settings && 'password' in data.issue_tracker_settings && !data.issue_tracker_settings.password) {
            delete data.issue_tracker_settings.password;
        }

        if (data.notify_channel_settings && 'token' in data.notify_channel_settings && !data.notify_channel_settings.token) {
            delete data.notify_channel_settings.token;
        }

        this.props.onSubmit(data)
    }

    render() {
        let passwordRules = [];
        if (!this.props.value) {
            passwordRules = [{ required: true, message: 'Пароль не должен быть пустым' }];
        }

        let tokenRules = [];
        if (!this.props.value) {
            tokenRules = [{ required: true, message: 'Токен не должено быть пустым' }];
        }

        return (
            <Form
                layout="vertical"
                name="basic"
                initialValues={WeekendFormValue(this.props.value)}
                onFinish={this.onFinish}
            >
                <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
                    <Col span={12}>
                        <Divider>Основные</Divider>
                        <Form.Item
                            label="Название команды"
                            name="title"
                            rules={[{ required: true, message: 'Название команды не должно быть пустым' }]}
                        >
                            <Input />
                        </Form.Item>

                        <Form.Item
                            label="Минимальное время которое нужно списать по задачам (сек.)"
                            name="min_work_log"
                            rules={[{ required: true, message: 'Название команды не должно быть пустым' }]}
                        >
                            <InputNumber  />
                        </Form.Item>

                        <Divider>Система ведения задач</Divider>

                        <Form.Item
                            label="Тип"
                            name="issue_tracker_type"
                            rules={[{ required: true}]}
                        >
                            <Select options={Dictionary.getTrackerTypes()} />
                        </Form.Item>

                        <Form.Item
                            label="Адрес"
                            name={["issue_tracker_settings", 'url']}
                            rules={[{ required: true, message: 'Адрес не должен быть пустым' }]}
                        >
                            <Input />
                        </Form.Item>

                        <Form.Item
                            label="Пользователь"
                            name={["issue_tracker_settings", 'username']}
                            rules={[{ required: true, message: 'Имя пользователя не должено быть пустым' }]}
                        >
                            <Input />
                        </Form.Item>

                        <Form.Item
                            label="Пароль"
                            name={["issue_tracker_settings", 'password']}
                            rules={passwordRules}
                        >
                            <Input.Password />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Divider>Канал оповещения</Divider>

                        <Form.Item
                            label="Тип"
                            name="notify_channel_type"
                            rules={[{ required: true}]}
                        >
                            <Select options={Dictionary.getChannelTypes()} />
                        </Form.Item>

                        <Form.Item
                            label="ID канала"
                            name="notify_channel_channel_id"
                            rules={[{ required: true, message: 'ID канала не должено быть пустым' }]}
                        >
                            <Input />
                        </Form.Item>

                        <Form.Item
                            label="Токен бота"
                            name={["notify_channel_settings", 'token']}
                            rules={tokenRules}
                        >
                            <Input.Password />
                        </Form.Item>

                        <Divider>Выходные</Divider>
                        <Weekend button_add_title={"Добавить выходные"}/>
                    </Col>
                </Row>

                <Form.Item>
                    <Button type="primary" htmlType="submit">
                        Сохранить
                    </Button>
                </Form.Item>
            </Form>
        )
    }
}

export default TeamSettingsForm