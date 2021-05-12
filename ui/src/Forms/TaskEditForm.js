import React from 'react';
import {Form, Input, Select, Button, Checkbox, InputNumber} from 'antd';
import Dictionary from "../Services/Dictionary";

const { TextArea } = Input;

class TaskEditForm extends React.Component {
    constructor(props) {
        super(props);

        let type = "check_work_log"
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

    checkWorkLogRender() {
        return (<span>
            <Form.Item
                label="Проекты в сиситеме задач (через запятую)"
                name="projects"
            >
                <Input/>
            </Form.Item>

            <Form.Item
                label="Изменение даты проверки"
                extra="Примеры -25h, -1.5h, 2h45m"
                name="check_date_modify"
            >
                    <Input/>
                </Form.Item>
        </span>)
    }

    sendTeamMessageRender() {
        return (<span>
            <Form.Item
                label="Сообщение которое будет отправлено"
                name="message"
                rules={[{ required: true, message: 'Не должно быть пустым' }]}
            >
                <TextArea rows={4} />
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
                    label="Тип задачи"
                    name="type"
                    rules={[{ required: true, message: 'Не должно быть пустым' }]}
                >
                    <Select options={Dictionary.getTaskTypes()} onChange={this.onChangeType} />
                </Form.Item>

                <Form.Item
                    name="active"
                    valuePropName="checked"
                >
                    <Checkbox>Активен</Checkbox>
                </Form.Item>

                <Form.Item
                    label="Выполнять по дням"
                    name="schedule_weekdays"
                >
                    <Checkbox.Group options={Dictionary.getWeekdays()} />
                </Form.Item>

                <Form.Item
                    label="Запускать в час"
                    name="schedule_hour"
                    rules={[{ required: true, message: 'Не должно быть пустым' }]}
                >
                    <InputNumber min={0} max={23} />
                </Form.Item>

                <Form.Item
                    label="Запускать в минуту"
                    name="schedule_minute"
                    rules={[{ required: true, message: 'Не должно быть пустым' }]}
                >
                    <InputNumber min={0} max={59} />
                </Form.Item>

                {this.state.type && this.state.type === "check_work_log" && this.checkWorkLogRender()}
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

export default TaskEditForm