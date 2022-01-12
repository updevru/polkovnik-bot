import React from 'react';
import {Form, Input, Select, Button, Checkbox} from 'antd';
import Dictionary from "../Services/Dictionary";

const { TextArea } = Input;

class ReceiverEditForm extends React.Component {
    constructor(props) {
        super(props);

        let type = "send_team_message"
        if (props.value && "type" in props.value) {
            type = props.value.type
        }

        this.state = {
            type: type
        }

        this.onFinish = this.onFinish.bind(this)
        this.onChangeType = this.onChangeType.bind(this)
    }

    onFinish(values: any) {
        this.props.onSubmit(values)
    }

    onChangeType(type) {
        this.setState({type: type})
    }

    getValues() {
        if (this.props.value) {
            return this.props.value;
        }

        return this.state
    }

    sendTeamMessageRender() {
        return (<span>
            <Form.Item
                label="Шаблон для формирования сообщения"
                name={["settings", "message"]}
                rules={[{ required: true, message: 'Не должно быть пустым' }]}
            >
                <TextArea rows={10} />
            </Form.Item>
        </span>)
    }

    render() {
        return (
            <Form
                layout="vertical"
                name="basic"
                initialValues={this.getValues()}
                onFinish={this.onFinish}
            >
                <Form.Item
                    label="Название"
                    name="title"
                    rules={[{ required: true, message: 'Не должно быть пустым' }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    name="active"
                    valuePropName="checked"
                >
                    <Checkbox>Активен</Checkbox>
                </Form.Item>

                <Form.Item
                    label="Тип задачи"
                    name="type"
                    rules={[{ required: true, message: 'Не должно быть пустым' }]}
                >
                    <Select options={Dictionary.getReceiverTypes()} onChange={this.onChangeType} />
                </Form.Item>

                <Form.Item
                    label="Формат воходящих данных"
                    name="format"
                    initialValue="auto"
                    rules={[{ required: true, message: 'Не должно быть пустым' }]}
                >
                    <Select options={Dictionary.getReceiverFormats()} />
                </Form.Item>

                {this.state.type && this.state.type === "send_team_message" && this.sendTeamMessageRender()}

                <Form.Item>
                    <Button type="primary" htmlType="submit">
                        Сохранить
                    </Button>
                </Form.Item>
            </Form>
        )
    }
}

export default ReceiverEditForm