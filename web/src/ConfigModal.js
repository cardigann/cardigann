import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Col, Modal, Button, Form, FormGroup, FormControl, ControlLabel }  from 'react-bootstrap';

class ConfigForm extends Component {
  state = {
    values: {},
  }
  getValues = () => {
    let values = {};
    Object.keys(this.refs).forEach((ref) => {
      values[ref] = ReactDOM.findDOMNode(this.refs[ref]).value
    });
    return values;
  }
  render() {
    let fields = this.props.fields.map((field) => {
      return (
        <FormGroup controlId={"formHorizontal" + field.name} key={field.name}>
          <Col componentClass={ControlLabel} sm={2}>{field.label}</Col>
            <Col sm={10}>
              <FormControl
                type={field.type}
                placeholder={field.placeholder}
                defaultValue={field.value}
                ref={field.name} />
          </Col>
        </FormGroup>
      );
    });
    return <Form horizontal>
      <FormGroup controlId="formHorizontalUrl">
        <Col componentClass={ControlLabel} sm={2}>
          URL
        </Col>
        <Col sm={10}>
          <FormControl type="text" placeholder="URL" defaultValue={this.props.url} ref="url" />
        </Col>
      </FormGroup>
      {fields}
    </Form>;
  }
}

class ConfigModal extends Component {
  static defaultProps = {
    onClose: () => {},
    onSave: () => {},
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
    let vals = this.refs.form.getValues();
    vals.enabled = "true";
    this.props.onSave(this.props.indexer, vals, () => {
      this.setState({show: false});
    });
  }
  buildFields = () => {
    return this.props.indexer.settings.map((s) => {
      if (typeof(this.state.config[s.name]) !== undefined) {
        s.value = this.state.config[s.name];
      }
      s.placeholder = s.placeholder || s.label
      return s;
    });
  }
  render() {
    return (
      <Modal show={this.state.show} onHide={this.handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>Configuration <small>for {this.props.indexer.name}</small></Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <ConfigForm fields={this.buildFields()} url={this.state.config.url} ref="form" />
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
