import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Modal, Button, Form, FormGroup, FormControl, ControlLabel }  from 'react-bootstrap';
import { BootstrapTable, TableHeaderColumn }  from 'react-bootstrap-table';
import xhrUrl from './xhr';
import spinner from './spinner.gif';
import queryString from 'query-string';

import 'react-bootstrap-table/dist/react-bootstrap-table.min.css';

class SearchForm extends Component {
  static defaultProps = {
    searching: false,
  }
  onSubmit = (e) => {
    e.preventDefault();
    this.props.onSearch({keywords: ReactDOM.findDOMNode(this.refs.keywords).value});
  }
  render() {
    return <Form inline onSubmit={this.onSubmit} className={this.props.searching?'searching':''}>
      <FormGroup controlId="formInlineName">
        <ControlLabel>Keywords</ControlLabel>
        {' '}
        <FormControl type="text" placeholder="" ref="keywords" />
      </FormGroup>
      {' '}
      <Button type="submit">Go</Button>
      {' '}
      <img src={spinner} height="50" width="50" alt="loading..." className="loading" />
    </Form>;
  }
}

class SearchModal extends Component {
  static defaultProps = {
    results: [],
  }
  static propTypes = {
   indexer: React.PropTypes.object.isRequired,
   apiKey: React.PropTypes.string.isRequired,
  }
  state = {
    apiKey: this.props.apiKey,
    show: this.props.show,
    searching: false,
    results: this.props.results,
  }
  componentWillReceiveProps(newProps) {
    this.setState({
      show: typeof(newProps).show !== undefined ? newProps.show : this.state.show,
    });
  }
  handleSearch = (query) => {
    this.setState({searching: true});
    fetch(xhrUrl("/torznab/"+this.props.indexer.id+"/api?"+queryString.stringify({
      t: "search",
      format: "json",
      apikey: this.state.apiKey,
      q: query.keywords,
    })))
    .then((response) => {
      if (!response.ok) {
        return response.json().then((resp) => {
          throw Error(resp.error);
        });
      }
      return response.json()
    })
    .then((results) => {
      console.log(results);
      this.setState({
        searching: false,
        results: results.Items,
      });
    })
    .catch((err) => {
      console.error(err);
      this.setState({searching: false});
    });
  }
  handleClose = () => {
    this.props.onClose();
    this.setState({show: false});
  }
  render() {
    let titleLinkFormatter = (cell, row) => {
      return "<a href="+row.Link+">"+cell+"</a>";
    }

    let fileSizeFormatter = (cell, row) => {
      var i = Math.floor( Math.log(cell) / Math.log(1024) );
      return ( cell / Math.pow(1024, i) ).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
    };

    return (
      <Modal show={this.state.show} onHide={this.handleClose} dialogClassName="App__SearchModal">
        <Modal.Header closeButton>
          <Modal.Title>Search <small>on {this.props.indexer.name}</small></Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <SearchForm onSearch={this.handleSearch} searching={this.state.searching}/>
          <hr />
          <div>
            <BootstrapTable
              data={this.state.results}
              striped={true}
              hover={true}
              pagination={true}
              >
              <TableHeaderColumn dataField="Title" isKey={true} dataSort={true} dataFormat={titleLinkFormatter} width="800px">Title</TableHeaderColumn>
              <TableHeaderColumn dataField="Size" dataSort={true} dataFormat={fileSizeFormatter} width="80px">Size</TableHeaderColumn>
              <TableHeaderColumn dataField="Category" dataSort={true} width="200px">Category</TableHeaderColumn>
              <TableHeaderColumn dataField="Seeders" dataSort={true} width="100px">Seeders</TableHeaderColumn>
              <TableHeaderColumn dataField="Peers" dataSort={true} width="100px">Peers</TableHeaderColumn>
              <TableHeaderColumn dataField="Site" dataSort={true} width="100px">Site</TableHeaderColumn>
            </BootstrapTable>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={this.handleClose}>Close</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}

export default SearchModal;