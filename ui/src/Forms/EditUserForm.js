import React from 'react';
import {Form, Input, Button} from 'antd';

class EditUserForm extends React.Component {
    constructor(props) {
        super(props);
        this.onFinish = this.onFinish.bind(this)
    }

    onFinish(values: any) {
        this.props.onSubmit(values)
    }

    render() {
        return (
            <Form
                layout="vertical"
                name="basic"
                initialValues={this.props.value}
                onFinish={this.onFinish}
            >
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