import React from 'react';
import {PageHeader, Spin} from "antd";
import TeamSettingsForm from "../Forms/TeamSettingsForm";
import ServerApi from "../Services/ServerApi";
import AlertMessage from "../Components/AlertMessage/AlertMessage";



class TeamSettingsPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            item: null,
            loading: true,
            message: {error: null, success: null}
        }

        this.handleSubmit = this.handleSubmit.bind(this)
    }

    componentDidMount() {
        this.loadItem(this.getTeamId())
    }

    async loadItem(teamId, userId) {
        let response = await ServerApi.teamSettings(teamId).list()
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({item: response, loading: false})
        }
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    async saveSettings(settings)
    {
        this.setState({loading: true})
        let response = await ServerApi.teamSettings(this.getTeamId()).add(settings)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({message: {success: "Настройки изменены"}, loading: false})
        }
    }

    handleSubmit(values) {
        this.saveSettings(values)
    };

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Настройки"
            >
                <Spin spinning={this.state.loading}>
                    <AlertMessage message={this.state.message} />
                    {this.state.item &&
                    <TeamSettingsForm value={this.state.item} onSubmit={this.handleSubmit}/>
                    }
                </Spin>
            </PageHeader>
        )
    }
}

export default TeamSettingsPage