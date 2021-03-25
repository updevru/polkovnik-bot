import React from 'react';
import { Card, List, Spin } from 'antd';
import {UserOutlined, BarsOutlined, SettingOutlined} from '@ant-design/icons';
import {Link} from "react-router-dom";
import ServerApi from "../Services/ServerApi";

const teamMenu = [
    {title: "Пользователи", icon: <UserOutlined />, url: function (team){
        return '/team/' + team.id + '/users'
    }},
    {title: "Задачи", icon: <BarsOutlined />, url: function (team){
            return '/team/' + team.id + '/tasks'
    }},
    // {title: "Расписание", icon: <CalendarOutlined />, url: function (team){
    //     return '/team/weekend'
    // }},
    {title: "Настройки", icon: <SettingOutlined />, url: function (team){
        return '/team/' + team.id + '/settings';
    }}
];

class HomePage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            teams: [],
            loading: true
        }
    }

    componentDidMount() {
        this.loadTeams();
    }

    async loadTeams() {
        let response = await ServerApi.team().list()
        if ('result' in response) {
            this.setState({teams: response.result, loading: false})
        } else {
            this.setState({teams: [], loading: false})
        }
    }

    render() {
        return (
            <div className="site-card-wrapper">
                <Spin spinning={this.state.loading}>
                <List
                    grid={{
                        gutter: 10,
                        column: 5,
                    }}
                    dataSource={this.state.teams}
                    renderItem={team => (
                        <List.Item>
                            <Card title={team.title} bodyStyle={{padding: 0}}>
                                <List
                                    size="small"
                                    dataSource={teamMenu}
                                    renderItem={item => (
                                        <List.Item>
                                            <Link to={item.url(team)}>{item.icon} {item.title}</Link>
                                        </List.Item>
                                    )}
                                />
                            </Card>
                        </List.Item>
                    )}
                />
                </Spin>
            </div>
        )
    }
}

export default HomePage