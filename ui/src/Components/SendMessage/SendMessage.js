import React, { useState } from 'react';
import {Modal, Button, Form, Input, Spin} from 'antd';
import {MessageOutlined} from "@ant-design/icons";
import ServerApi from "../../Services/ServerApi";
import AlertMessage from "../AlertMessage/AlertMessage";

const SendMessage = ({teamId}) => {
    const [form] = Form.useForm();
    const [visible, setVisible] = useState(false);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState({error: null, success: null});

    async function send(text) {
        setLoading(true);
        const response = await ServerApi.sendTeamMessage(teamId, text);

        if ('error' in response) {
            setMessage({error: response.error})
        } else {
            setMessage({success: "Сообщение отправлено"})
        }

        setLoading(false);
    }

    const handleOk = () => {
        form.validateFields()
        .then((values) => {
            form.resetFields();
            send(values.text);
        });
    }

    return <>
        <Button type="dashed" onClick={() => {setVisible(true)}}>
            <MessageOutlined />
        </Button>
        <Modal title="Отправка сообщения команде" visible={visible} onOk={handleOk} onCancel={() => {setVisible(false)}}>
            <Spin spinning={loading}>
                <AlertMessage message={message} />
                <Form
                    form={form}
                    layout="vertical"
                    name="basic"
                >
                    <Form.Item
                        name="text"
                        label="Сообщение"
                        rules={[{ required: true, message: 'Заполните текст сообщения' }]}
                    >
                        <Input.TextArea rows={6} />
                    </Form.Item>
                </Form>
            </Spin>
        </Modal>
    </>
}

export default SendMessage