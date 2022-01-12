import React from 'react';
import {PageHeader, Spin} from 'antd';
import ServerApi from "../Services/ServerApi";
import AlertMessage from "../Components/AlertMessage/AlertMessage";
import ReceiverEditForm from "../Forms/ReceiverEditForm";

class ReceiverAddPage extends React.Component{

    constructor(props) {
        super(props);
        this.state = {
            item: {},
            loading: false,
            message: {error: null, success: null}
        }

        this.saveReceiver = this.saveReceiver.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    getTeamId() {
        return this.props.match.params.teamId;
    }

    async saveReceiver(receiver)
    {
        this.setState({loading: true})
        let response = await ServerApi.receiver(this.getTeamId()).add(receiver)
        if ('error' in response) {
            this.setState({message: {error: response.error}, loading: false})
        } else {
            this.setState({message: {success: "Приемник добавлен"}, loading: false})
        }

    }

    handleSubmit(values) {
        this.saveReceiver(values)
    };

    render() {
        return (
            <PageHeader
                onBack={() => window.history.back()}
                title="Новый приемник"
            >
                <Spin spinning={this.state.loading}>
                    <div className={"app-form-container"}>
                        <AlertMessage message={this.state.message} />
                        <ReceiverEditForm onSubmit={this.handleSubmit}/>
                    </div>
                </Spin>
            </PageHeader>
        )
    }
}

export default ReceiverAddPage