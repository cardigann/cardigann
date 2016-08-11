import React, { Component } from 'react';
import { HelpBlock, Col, Modal, Button, Form, FormGroup, FormControl, ControlLabel} from 'react-bootstrap';

class Login extends Component {
  static defaultProps = {
    authUrl: "/xhr/auth",
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
    this.postAuth(this.state.passphrase);
  }
  handleAuthError = (err) => {
    this.setState({
      helpBlock: err,
      validationState: "error",
    });
  }
  postAuth = (passphrase) => {
    let handleAuth = (data) => this.props.onAuthenticate(data.token);
    let handleAuthError = (data) => this.handleAuthError(data.error);

    fetch(this.props.authUrl, {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        method: "POST",
        body: JSON.stringify({ passphrase: passphrase })
    })
    .then(function(res){
      if(res.ok) {
        console.log("ok response");
        res.json().then(handleAuth);
      } else {
        console.log("error response");
        res.json().then(handleAuthError);
      }
    })
    .catch(function(res){
        this.setState({
          helpBlock: "Network connection error",
          validationState: "error",
        });
    })
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