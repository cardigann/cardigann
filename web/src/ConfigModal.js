import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Col, Modal, Button, Form, FormGroup, FormControl, ControlLabel} from 'react-bootstrap';

class ConfigModal extends Component {
  static defaultProps = {
    onClose: () => {},
    onSave: () => {},
  }
  static propTypes = {
   indexer: React.PropTypes.object.isRequired,
   config: React.PropTypes.object.isRequired,
  }
  state = {
    config: this.props.config,
    show: this.props.show,
  }
  componentWillReceiveProps(newProps) {
    this.setState({
      show: typeof(newProps).show !== undefined ? newProps.show : this.state.show,
      config: typeof(newProps).config !== undefined ? newProps.config : this.state.config,
    });
  }
  handleClose = () => {
    this.props.onClose(this.props.indexer);
    this.setState({show: false});
  }
  handleSave = () => {
    this.props.onSave(this.props.indexer, {
      url: ReactDOM.findDOMNode(this.refs.url).value,
      username: ReactDOM.findDOMNode(this.refs.username).value,
      password: ReactDOM.findDOMNode(this.refs.password).value,
      enabled: "true"
    }, () => {});
  }
  render() {
    return (
      <Modal show={this.state.show} onHide={this.handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>Configuration <small>for {this.props.indexer.name}</small></Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form horizontal>
            <FormGroup controlId="formHorizontalUrl">
              <Col componentClass={ControlLabel} sm={2}>
                URL
              </Col>
              <Col sm={10}>
                <FormControl type="text" placeholder="URL" defaultValue={this.state.config.url} ref="url" />
              </Col>
            </FormGroup>
            <FormGroup controlId="formHorizontalUsername">
              <Col componentClass={ControlLabel} sm={2}>
                Username
              </Col>
              <Col sm={10}>
                <FormControl type="username" placeholder="Username" defaultValue={this.state.config.username} ref="username" />
              </Col>
            </FormGroup>
            <FormGroup controlId="formHorizontalPassword">
              <Col componentClass={ControlLabel} sm={2}>
                Password
              </Col>
              <Col sm={10}>
                <FormControl type="password" placeholder="Password" defaultValue={this.state.config.password} ref="password" />
              </Col>
            </FormGroup>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button bsStyle="primary" onClick={this.handleSave}>Save and Close</Button>
          <Button onClick={this.handleClose}>Cancel</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}

export default ConfigModal;