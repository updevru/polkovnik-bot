import React from 'react';
import {PageHeader, Spin, Form, Input, Button, InputNumber, Select, Divider, Alert} from 'antd';
import ServerApi from "../Services/ServerApi";
import AlertMessage from "../Components/AlertMessage/AlertMessage";
import TaskEditForm from "../Forms/TaskEditForm";

class TaskEditPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            item: null,
            loading: true,
            message: {error: null, success: null}
        }

        this.saveTask = this.saveTask.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    componentDidMount() {
        this.loadItem(this.getTeamId(), this.getTaskId())
    }

    async loadItem(teamId, taskId) {
        let response = await ServerApi.task(teamId).get(taskId)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({item: response, loading: false})
        }
    }

    getTaskId() {
        return this.props.match.params.taskId;
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    async saveTask(task)
    {
        this.setState({loading: true})
        let response = await ServerApi.task(this.getTeamId()).edit(this.getTaskId(), task)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({message: {success: "Задача изменена"}, loading: false})
        }

    }

    handleSubmit(values) {
        this.saveTask(values)
    };

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Редактирование задачи"
            >
                <Spin spinning={this.state.loading}>
                    <div className={"app-form-container"}>
                        <AlertMessage message={this.state.message} />
                        {this.state.item &&
                            <TaskEditForm value={this.state.item} onSubmit={this.handleSubmit}/>
                        }
                    </div>
                </Spin>
            </PageHeader>
        )
    }
}

export default TaskEditPage