import React from 'react';
import {PageHeader, Spin} from "antd";
import ServerApi from "../Services/ServerApi";
import AlertMessage from "../Components/AlertMessage/AlertMessage";
import TeamAddForm from "../Forms/TeamAddForm";

class TeamAddPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            item: null,
            loading: false,
            message: {error: null, success: null}
        }

        this.handleSubmit = this.handleSubmit.bind(this)
    }

    async saveTeam(team)
    {
        this.setState({loading: true})
        let response = await ServerApi.team().add(team)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({message: {success: "Команда создана"}, loading: false})
        }
    }

    handleSubmit(values) {
        this.saveTeam(values)
    };

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Новая команда"
            >
                <div className={"app-form-container"}>
                    <Spin spinning={this.state.loading}>
                        <AlertMessage message={this.state.message} />
                        <TeamAddForm onSubmit={this.handleSubmit}/>
                    </Spin>
                </div>
            </PageHeader>
        )
    }
}

export default TeamAddPage