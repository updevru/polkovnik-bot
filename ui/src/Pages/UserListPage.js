import React from 'react';
import { Table, PageHeader, Spin, Space, Button, Modal  } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import ServerApi from "../Services/ServerApi";
import {Link} from "react-router-dom";
import AlertMessage from "../Components/AlertMessage/AlertMessage";

const { confirm } = Modal;

class UserListPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            list: [],
            loading: true,
            message: {error: null, success: null}
        }

        this.deleteUser = this.deleteUser.bind(this)
        this.showConfirm = this.showConfirm.bind(this)
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    componentDidMount() {
        this.loadList(this.getTeamId())
    }

    getColumns() {
        return [
            {
                title: 'Имя',
                dataIndex: 'name',
                key: 'name',
            },
            {
                title: 'Логин',
                dataIndex: 'login',
                key: 'login',
            },
            {
                title: 'Ник',
                dataIndex: 'nickname',
                key: 'nickname',
            },
            {
                title: 'Действия',
                key: 'action',
                render: (text, record) => (
                    <Space size="middle">
                        <Button type="primary">
                            <Link to={"/team/" + this.getTeamId() + "/users/edit/" + record.id}><EditOutlined /> Редактировать</Link>
                        </Button>
                        <Button type="primary" danger onClick={() => this.showConfirm(record.id)}>
                            <DeleteOutlined /> Удалить
                        </Button>
                    </Space>
                ),
            }
        ]
    }

    async loadList(teamId) {
        let response = await ServerApi.user(teamId).list()
        if ('result' in response) {
            this.setState({list: response.result, loading: false})
        } else {
            this.setState({list: [], loading: false})
        }
    }

    deleteUser(userId)
    {
        this.setState({loading: true})
        let response = ServerApi.user(this.getTeamId()).delete(userId)
        if ('error' in response) {
            this.setState({message: {error: response.error}})
        } else {
            this.setState({message: {success: "Пользователь удален"}})
        }
        this.loadList(this.getTeamId())
    }

    showConfirm(userId) {
        let self = this
        confirm({
            title: 'Точно удалить пользователя?',
            icon: <DeleteOutlined />,
            onOk() {
                self.deleteUser(userId)
            }
        });
    }

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Пользователи"
                extra={[
                    <Button type="primary">
                        <Link to={"/team/" + this.getTeamId() + "/users/add"}><PlusOutlined /> Добавить</Link>
                    </Button>
                ]}
            >
                <AlertMessage message={this.state.message} />
                <Spin spinning={this.state.loading}>
                    <Table columns={this.getColumns()} dataSource={this.state.list} />
                </Spin>
            </PageHeader>
        )
    }
}

export default UserListPage