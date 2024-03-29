import React from 'react';
import {Table, Tag, Space, PageHeader, Button, Modal, Spin} from 'antd';
import ServerApi from "../Services/ServerApi";
import {Link} from "react-router-dom";
import {DeleteOutlined, EditOutlined, PlusOutlined, HistoryOutlined, CaretRightOutlined} from "@ant-design/icons";
import Dictionary from "../Services/Dictionary";

const { confirm } = Modal;

class TaskListPage extends React.Component{

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
        let response = await ServerApi.task(teamId).list()
        if ('result' in response) {
            this.setState({list: response.result, loading: false})
        } else {
            this.setState({list: [], loading: false})
        }
    }

    deleteTask(taskId)
    {
        this.setState({loading: true})
        let response = ServerApi.task(this.getTeamId()).delete(taskId)
        if ('error' in response) {
            this.setState({message: {error: response.error}})
        } else {
            this.setState({message: {success: "Задача удалена"}})
        }
        this.loadList(this.getTeamId())
    }

    showConfirm(userId) {
        let self = this
        confirm({
            title: 'Точно удалить задачу?',
            icon: <DeleteOutlined />,
            onOk() {
                self.deleteTask(userId)
            }
        });
    }

    async runTask(taskId)
    {
        this.setState({loading: true})
        let response = await ServerApi.runTask(this.getTeamId(), taskId)
        if ('error' in response) {
            Modal.error({title: 'Произошла ошибка', content: response.error})
        } else {
            Modal.info({title: 'Задача запущена'})
        }
        this.setState({loading: false})
    }

    /**
     * Добавляет ведущий ноль к минутам если нужно
     * @param minutes
     * @returns {string|*}
     */
    renderMinutes(minutes)
    {
        if (minutes < 10) {
            return  "0" + minutes
        }

        return minutes;
    }

    getColumns() {
        return [
            {
                title: '',
                key: 'run',
                render: (text, record) => (
                    <Space size="middle">
                        <Button shape="circle" onClick={() => this.runTask(record.id)} icon={<CaretRightOutlined />} />
                    </Space>
                ),
            },
            {
                title: 'Задание',
                key: 'type',
                render: (text, record) => (
                    <span>
                        {Dictionary.getTaskType(record.type).label}
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
                title: 'Расписание',
                key: 'type',
                render: (text, record) => (
                    <div>
                        {record.schedule_weekdays.map((day, i) => {

                            return (<Tag>{day}</Tag>)
                        })}
                        в <Tag>{record.schedule_hour}:{this.renderMinutes(record.schedule_minute)}</Tag>
                    </div>
                )
            },
            {
                title: 'Последний запуск',
                dataIndex: 'last_run_time',
                key: 'last_run_time',
            },
            {
                title: 'Действия',
                key: 'action',
                render: (text, record) => (
                    <Space size="middle">
                        <Button type="">
                            <Link to={"/team/" + this.getTeamId() + "/tasks/" + record.id + "/history"}><HistoryOutlined /></Link>
                        </Button>
                        <Button type="primary">
                            <Link to={"/team/" + this.getTeamId() + "/tasks/edit/" + record.id}><EditOutlined /> Редактировать</Link>
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
                title="Задачи"
                extra={[
                    <Button type="primary">
                        <Link to={"/team/" + this.getTeamId() + "/tasks/add"}><PlusOutlined /> Добавить</Link>
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

export default TaskListPage