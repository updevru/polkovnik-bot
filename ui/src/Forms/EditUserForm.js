import React from 'react';
import {Form, Input, Button, Checkbox, Select} from 'antd';
import {Weekend, WeekendDataValue, WeekendFormValue} from "../Components/Weekend/Weekend";
import {Option} from "antd/es/mentions";

class EditUserForm extends React.Component {
    constructor(props) {
        super(props);
        this.onFinish = this.onFinish.bind(this)
    }

    onFinish(values: any) {
        this.props.onSubmit(WeekendDataValue(values))
    }

    render() {
        return (
            <Form
                layout="vertical"
                name="basic"
                initialValues={WeekendFormValue(this.props.value)}
                onFinish={this.onFinish}
            >
                <Form.Item
                    name="active"
                    valuePropName="checked"
                >
                    <Checkbox>Активен</Checkbox>
                </Form.Item>

                <Form.Item
                    label="Имя"
                    name="name"
                    rules={[{ required: true, message: 'Не должно быть пустым' }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="Логин в системе ведения задач (Jira)"
                    name="login"
                    rules={[{ required: true, message: 'Не должно быть пустым' }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="Ник"
                    name="nickname"
                    rules={[{ required: true, message: 'Не должно быть пустым' }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="Пол пользователя"
                    name="gender"
                    rules={[{required: true}]}
                >
                    <Select
                        defaultValue="male"
                    >
                        <Option value="male">Мужской</Option>
                        <Option value="female">Женский</Option>
                    </Select>
                </Form.Item>

                <Weekend />

                <Form.Item>
                    <Button type="primary" htmlType="submit">
                        Сохранить
                    </Button>
                </Form.Item>
            </Form>
        )
    }
}

export default EditUserForm