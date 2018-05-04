import React, { Component } from 'react';
import { HelpBlock, Col, Modal, Button, Form, FormGroup, FormControl, ControlLabel} from 'react-bootstrap';
import xhrUrl from './xhr';

class Login extends Component {
  static defaultProps = {
    authUrl: "./xhr/auth",
    onAuthenticate: () => {},
  }
  state = {
    show: true,
    passphrase: "",
    helpBlock: "",
    validationState: null,
  }
  handlePassphraseChange = (e) => {
    this.setState({passphrase: e.target.value})
  }
  handleSubmit = (e) => {
    e.preventDefault();
    this.authenticate(this.state.passphrase);
  }
  handleAuthError = (err) => {
    console.error("Auth error:", err);
    this.setState({
      helpBlock: err,
      validationState: "error",
    });
  }
  authenticate = (passphrase) => {
    fetch(xhrUrl(this.props.authUrl), {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        method: "POST",
        body: JSON.stringify({ passphrase: passphrase })
    })
    .then((res) => {
      if (!res.ok) {
        throw Error("Failed XHR request: "+res.statusText);
      }
      return res;
    })
    .then((res) => {
      res.json().then((data) => {
        if (!data.hasOwnProperty("error")) {
          this.props.onAuthenticate(data.token);
        } else {
          this.handleAuthError(data.error);
        }
      });
      return res;
    })
    .catch((err) => {
      this.handleAuthError(err.message)
    });
  }
  render() {
    return (
      <Modal show={this.state.show}>
        <Modal.Header closeButton>
          <Modal.Title>Login</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form horizontal onSubmit={this.handleSubmit}>
            <FormGroup controlId="formHorizontalPassword" validationState={this.state.validationState}>
              <Col componentClass={ControlLabel} sm={2}>
                Password
              </Col>
              <Col sm={10}>
                <FormControl type="password" placeholder="Passphrase" onChange={this.handlePassphraseChange} />
                <HelpBlock>{this.state.helpBlock}</HelpBlock>
              </Col>
            </FormGroup>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button bsStyle="primary" onClick={this.handleSubmit}>Login</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}

export default Login;