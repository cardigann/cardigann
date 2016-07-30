import React, { Component } from 'react';
import { Col, Modal, Button, Form, FormGroup, FormControl, ControlLabel} from 'react-bootstrap';

class ConfigModal extends Component {
  static defaultProps = {
    onClose: () => {},
    onSave: () => {},
  }
  state = {
    show: this.props.show,
  }
  componentWillReceiveProps(newProps) {
    if (typeof(newProps.show) !== undefined) {
      this.setState({show: newProps.show});
    }
  }
  handleClose = () => {
    this.props.onClose(this.props.indexer);
    this.setState({show: false});
  }
  handleSave = () => {
    this.props.onSave(this.props.indexer);
    this.setState({show: false}, function() {
      console.log("post save", this.state);
    });
  }
  render() {
    return (
      <Modal show={this.state.show} onHide={this.handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>Configuration <small>for {this.props.indexer.name}</small></Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form horizontal>
            <FormGroup controlId="formHorizontalEmail">
              <Col componentClass={ControlLabel} sm={2}>
                Email
              </Col>
              <Col sm={10}>
                <FormControl type="email" placeholder="Email" />
              </Col>
            </FormGroup>
            <FormGroup controlId="formHorizontalPassword">
              <Col componentClass={ControlLabel} sm={2}>
                Password
              </Col>
              <Col sm={10}>
                <FormControl type="password" placeholder="Password" />
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