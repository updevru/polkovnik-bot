import React from 'react';
import {Table, Tag, Space, PageHeader} from 'antd';
import ServerApi from "../Services/ServerApi";

const columns = [
    {
        title: 'Type',
        dataIndex: 'type',
        key: 'type',
    },
    {
        title: 'Last run time',
        dataIndex: 'last_run_time',
        key: 'last_run_time',
    },
];

class TaskListPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            list: [],
            loading: true
        }
    }

    componentDidMount() {
        this.loadList(this.props.match.params.teamId)
    }

    async loadList(teamId) {
        let response = await ServerApi.task(teamId).list()
        if ('result' in response) {
            this.setState({list: response.result, loading: false})
        } else {
            this.setState({list: [], loading: false})
        }
    }

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Задачи"
            >
                <Table columns={columns} dataSource={this.state.list} />
            </PageHeader>
        )
    }
}

export default TaskListPage