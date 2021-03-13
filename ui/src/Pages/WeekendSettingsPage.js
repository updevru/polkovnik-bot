import React from 'react';
import {PageHeader, Form, Input, Button, Checkbox  } from "antd";

const days = [
    { label: 'Monday', value: 'Monday' },
    { label: 'Tuesday', value: 'Tuesday' },
    { label: 'Wednesday', value: 'Wednesday' },
    { label: 'Thursday', value: 'Thursday' },
    { label: 'Friday', value: 'Friday' },
    { label: 'Saturday', value: 'Saturday' },
    { label: 'Sunday', value: 'Sunday' },
];

class WeekendSettingsPage extends React.Component{

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
                title="Расписание"
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
                        label="Рабочие дни недели"
                        name="weekdays"
                    >
                        <Checkbox.Group options={days} />
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

export default WeekendSettingsPage