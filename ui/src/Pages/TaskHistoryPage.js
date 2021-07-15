import React from 'react';
import {Table, PageHeader, Spin} from 'antd';
import ServerApi from "../Services/ServerApi";
import {CheckCircleOutlined, ExclamationCircleOutlined} from "@ant-design/icons";

class TaskHistoryPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            list: [],
            loading: true,
            pagination: {
                current: 1,
                pageSize: 30,
            },
        }
        this.handleTableChange = this.handleTableChange.bind(this);
    }

    componentDidMount() {
        this.loadList(this.getTeamId(), this.getTaskId(), this.state.pagination)
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    getTaskId() {
        return this.props.match.params.taskId;
    }

    async loadList(teamId, taskId, pagination) {
        this.setState({loading: true})
        let response = await ServerApi.taskHistory(teamId, taskId).list({page: pagination.current, size: pagination.pageSize})
        if ('result' in response) {
            let result = response.result.map((item, i) => {
                item['key'] = i + 1
                return item
            });
            this.setState({list: result, loading: false, pagination: {...pagination, total: response.total}})
        } else {
            this.setState({list: [], loading: false})
        }
    }

    getColumns() {
        return [
            {
                title: 'Дата',
                dataIndex: 'date',
                key: 'date',
            },
            {
                title: 'Статус',
                key: 'status',
                render: (text, record) => (
                    <span>
                        {record.success && <CheckCircleOutlined style={{"color": "green"}}/>}
                        {record.error && <ExclamationCircleOutlined style={{"color": "red"}}/>}
                    </span>
                )
            },
        ]
    }

    handleTableChange = (pagination) => {
        this.loadList(this.getTeamId(), this.getTaskId(), pagination)
    };

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="История запусков"
            >
                <Spin spinning={this.state.loading}>
                    <Table
                        columns={this.getColumns()}
                        dataSource={this.state.list}
                        pagination={this.state.pagination}
                        onChange={this.handleTableChange}
                        expandable={{
                            expandedRowRender: record => <pre style={{ margin: 0 }}>
                                {record.logs.map((log, i) => {

                                    return (<div><span className={"log-line"}>{i+1}.</span> {log}</div>)
                                })}
                            </pre>,
                            rowExpandable: record => record.logs.length > 0,
                        }}
                    />
                </Spin>
            </PageHeader>
        )
    }
}

export default TaskHistoryPage