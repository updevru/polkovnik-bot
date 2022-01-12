import React from 'react';
import {PageHeader, Spin} from 'antd';
import ServerApi from "../Services/ServerApi";
import AlertMessage from "../Components/AlertMessage/AlertMessage";
import ReceiverEditForm from "../Forms/ReceiverEditForm";

class ReceiverEditPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            item: null,
            loading: true,
            message: {error: null, success: null}
        }

        this.saveReceiver = this.saveReceiver.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    componentDidMount() {
        this.loadItem(this.getTeamId(), this.getReceiverId())
    }

    async loadItem(teamId, receiverId) {
        let response = await ServerApi.receiver(teamId).get(receiverId)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({item: response, loading: false})
        }
    }

    getReceiverId() {
        return this.props.match.params.receiverId;
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    async saveReceiver(receiver)
    {
        this.setState({loading: true})
        let response = await ServerApi.receiver(this.getTeamId()).edit(this.getReceiverId(), receiver)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({message: {success: "Приемник изменен"}, loading: false})
        }

    }

    handleSubmit(values) {
        this.saveReceiver(values)
    };

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Редактирование приемника"
            >
                <Spin spinning={this.state.loading}>
                    <div className={"app-form-container"}>
                        <AlertMessage message={this.state.message} />
                        {this.state.item &&
                            <ReceiverEditForm value={this.state.item} onSubmit={this.handleSubmit}/>
                        }
                    </div>
                </Spin>
            </PageHeader>
        )
    }
}

export default ReceiverEditPage