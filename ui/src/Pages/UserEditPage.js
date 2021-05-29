import React from 'react';
import {PageHeader, Spin} from 'antd';
import ServerApi from "../Services/ServerApi";
import AlertMessage from "../Components/AlertMessage/AlertMessage";
import EditUserForm from "../Forms/EditUserForm";

class UserEditPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            item: null,
            loading: true,
            message: {error: null, success: null}
        }

        this.saveUser = this.saveUser.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    componentDidMount() {
        this.loadItem(this.getTeamId(), this.getUserId())
    }

    async loadItem(teamId, userId) {
        let response = await ServerApi.user(teamId).get(userId)
        console.log("Response", response)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({item: response, loading: false})
        }
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    getUserId() {
        return this.props.match.params.userId;
    }

    async saveUser(user)
    {
        this.setState({loading: true})
        let response = await ServerApi.user(this.getTeamId()).edit(this.getUserId(), user)
        console.log('Response:', response)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({message: {success: "Пользователь отредактирован"}, loading: false})
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
                title="Редактирование пользователя"
            >
                <Spin spinning={this.state.loading}>
                    <div className={"app-form-container"}>
                        <AlertMessage message={this.state.message} />
                        {this.state.item &&
                            <EditUserForm value={this.state.item} onSubmit={this.handleSubmit}/>
                        }
                    </div>
                </Spin>
            </PageHeader>
        )
    }
}

export default UserEditPage