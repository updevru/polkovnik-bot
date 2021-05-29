import React from 'react';
import {PageHeader, Spin} from 'antd';
import ServerApi from "../Services/ServerApi";
import AlertMessage from "../Components/AlertMessage/AlertMessage";
import TaskEditForm from "../Forms/TaskEditForm";

class TaskAddPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            item: {},
            loading: false,
            message: {error: null, success: null}
        }

        this.saveTask = this.saveTask.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    async saveTask(task)
    {
        this.setState({loading: true})
        let response = await ServerApi.task(this.getTeamId()).add(task)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({message: {success: "Задача добавлена"}, loading: false})
        }

    }

    handleSubmit(values) {
        this.saveTask(values)
    };

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Новая задача"
            >
                <Spin spinning={this.state.loading}>
                    <div className={"app-form-container"}>
                        <AlertMessage message={this.state.message} />
                        <TaskEditForm onSubmit={this.handleSubmit}/>
                    </div>
                </Spin>
            </PageHeader>
        )
    }
}

export default TaskAddPage