import React, { Component } from 'react';
import { Alert } from 'react-bootstrap';

class AlertDismissable extends Component {
  state = {
    visible: true
  }
  handleAlertDismiss = () => {
    this.setState({visible: false});
  }
  handleAlertShow = () => {
    this.setState({visible: true});
  }
  render() {
    if (this.state.visible) {
      return (
        <Alert bsStyle="danger" onDismiss={this.handleAlertDismiss}>
          {this.props.children}
        </Alert>
      );
    }
    return (
      <div />
    );
  }
}

export default AlertDismissable;
