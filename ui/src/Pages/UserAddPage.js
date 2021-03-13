import React from 'react';
import {PageHeader, Spin, Form, Input, Button, InputNumber, Select, Divider, Alert} from 'antd';
import ServerApi from "../Services/ServerApi";
import AlertMessage from "../Components/AlertMessage/AlertMessage";
import EditUserForm from "../Forms/EditUserForm";

class UserAddPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            item: {},
            loading: false,
            message: {error: null, success: null}
        }

        this.saveUser = this.saveUser.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    async saveUser(user)
    {
        this.setState({loading: true})
        let response = await ServerApi.user(this.getTeamId()).add(user)
        console.log('Response:', response)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({message: {success: "Пользователь добавлен"}, loading: false})
        }

    }

    handleSubmit(values) {
        console.log('Success:', values);
        this.saveUser(values)
    };

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Новый пользователь"
            >
                <Spin spinning={this.state.loading}>
                    <div className={"app-form-container"}>
                        <AlertMessage message={this.state.message} />
                        <EditUserForm onSubmit={this.handleSubmit}/>
                    </div>
                </Spin>
            </PageHeader>
        )
    }
}

export default UserAddPage