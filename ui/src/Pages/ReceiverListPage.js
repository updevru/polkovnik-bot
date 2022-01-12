import React from 'react';
import {Table, Space, PageHeader, Button, Modal, Spin} from 'antd';
import ServerApi from "../Services/ServerApi";
import {Link} from "react-router-dom";
import {DeleteOutlined, EditOutlined, PlusOutlined} from "@ant-design/icons";
import Dictionary from "../Services/Dictionary";

const { confirm } = Modal;

class ReceiverListPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            list: [],
            loading: true
        }
    }

    componentDidMount() {
        this.loadList(this.getTeamId())
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    async loadList(teamId) {
        let response = await ServerApi.receiver(teamId).list()
        if ('result' in response) {
            this.setState({list: response.result, loading: false})
        } else {
            this.setState({list: [], loading: false})
        }
    }

    deleteReceiver(receiverId)
    {
        this.setState({loading: true})
        let response = ServerApi.receiver(this.getTeamId()).delete(receiverId)
        if ('error' in response) {
            this.setState({message: {error: response.error}})
        } else {
            this.setState({message: {success: "Приемник удален"}})
        }
        this.loadList(this.getTeamId())
    }

    showConfirm(userId) {
        let self = this
        confirm({
            title: 'Точно удалить приемник?',
            icon: <DeleteOutlined />,
            onOk() {
                self.deleteReceiver(userId)
            }
        });
    }

    getColumns() {
        return [
            {
                title: 'Название',
                key: 'title',
                dataIndex: 'title',
            },
            {
                title: 'Действие',
                key: 'type',
                render: (text, record) => (
                    <span>
                        {Dictionary.getReceiverType(record.type).label}
                    </span>
                )
            },
            {
                title: 'Активен',
                key: 'active',
                render: (text, record) => (
                    <span>
                        {record.active === true && "Да"}
                        {record.active === false && "Нет"}
                    </span>
                )
            },
            {
                title: 'Действия',
                key: 'action',
                render: (text, record) => (
                    <Space size="middle">
                        <Button type="primary">
                            <Link to={"/team/" + this.getTeamId() + "/receivers/edit/" + record.id}><EditOutlined /> Редактировать</Link>
                        </Button>
                        <Button type="primary" danger onClick={() => this.showConfirm(record.id)}>
                            <DeleteOutlined /> Удалить
                        </Button>
                    </Space>
                ),
            }
        ]
    }

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Приемники"
                extra={[
                    <Button type="primary">
                        <Link to={"/team/" + this.getTeamId() + "/receivers/add"}><PlusOutlined /> Добавить</Link>
                    </Button>
                ]}
            >
                <Spin spinning={this.state.loading}>
                    <Table columns={this.getColumns()} dataSource={this.state.list} />
                </Spin>
            </PageHeader>
        )
    }
}

export default ReceiverListPage