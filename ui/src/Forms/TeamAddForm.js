import React from 'react';
import {Button, Form, Input} from "antd";

class TeamAddForm extends React.Component {

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
                    label="Название команды"
                    name="title"
                    rules={[{ required: true, message: 'Название команды не должно быть пустым' }]}
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

export default TeamAddForm