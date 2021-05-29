import React from 'react';
import {Alert} from 'antd';

class AlertMessage extends React.Component {
    render() {
        return (
            <span>
                {'error' in this.props.message && this.props.message.error &&
                <Alert
                    message="Error"
                    description={this.props.message.error}
                    type="error"
                    closable
                />
                }
                {'success' in this.props.message && this.props.message.success &&
                <Alert
                    message={this.props.message.success}
                    type="success"
                    closable
                />
                }
            </span>
        )
    }
}

export default AlertMessage