import React from 'react';
import {Card, List, Spin, Row, Col, PageHeader, Button} from 'antd';
import {UserOutlined, BarsOutlined, SettingOutlined, PlusOutlined, ApiOutlined} from '@ant-design/icons';
import {Link} from "react-router-dom";
import ServerApi from "../Services/ServerApi";
import SendMessage from "../Components/SendMessage/SendMessage";

const teamMenu = [
    {title: "Пользователи", icon: <UserOutlined />, url: function (team){
        return '/team/' + team.id + '/users'
    }},
    {title: "Задачи", icon: <BarsOutlined />, url: function (team){
            return '/team/' + team.id + '/tasks'
    }},
    {title: "Прием webhook", icon: <ApiOutlined />, url: function (team){
            return '/team/' + team.id + '/receivers'
        }},
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
            <PageHeader
                title="Команды"
                extra={[
                    <Button type="primary">
                        <Link to={"/team/add"}><PlusOutlined /> Добавить</Link>
                    </Button>
                ]}
            >
            <Spin spinning={this.state.loading}>
                <Row gutter={[16, 24]}>
                    {this.state.teams.map((team, i) => {
                        return <Col className="gutter-row" span={6}>
                            <Card title={team.title} bodyStyle={{padding: 0}} extra={<SendMessage teamId={team.id}/>}>
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
                        </Col>
                    })}
                </Row>
            </Spin>
            </PageHeader>
        )
    }
}

export default HomePage