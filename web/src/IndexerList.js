import React, { Component } from 'react';
import { Table, ButtonToolbar, Button, FormControl } from 'react-bootstrap';
import ConfigModal from './ConfigModal';

class EditButton extends Component {
  state = {
    isSaving: false,
  }
  handleClick = (e) => {
    this.props.onClick(e);
  }
  render() {
    let isSaving = this.state.isSaving;
    return (
      <Button
        bsStyle="primary"
        bsSize="xsmall"
        disabled={isSaving}
        onClick={!isSaving ? this.handleClick : null}>
        {isSaving ? 'Saving...' : 'Edit'}
      </Button>
    );
  }
}

class FeedLink extends Component {
  render() {
    return (
      <FormControl type="text" value={this.props.feedHref} className="input-sm" readOnly />
    );
  }
}

class IndexerListRow extends Component {
  state = {
    showEditModal: false,
  }
  handleEditClick = () => {
    this.setState({showEditModal: true});
  }
  render() {
    return (
      <tr>
        <td className="col-md-2">{this.props.indexer.name}</td>
        <td className="col-md-6">
          <FeedLink
            feedHref={this.props.indexer.feeds.torznab}
            label="torznab" />
        </td>
        <td className="col-md-2">Working</td>
        <td className="col-md-2">
          <ButtonToolbar>
            <EditButton
              onClick={this.handleEditClick} />
            <ConfigModal
              show={this.state.showEditModal}
              indexer={this.props.indexer}
              onSave={this.props.onSaveIndexer} />
            <Button
              bsSize="xsmall">Test</Button>
            <Button
              bsSize="xsmall"
              bsStyle="danger">Delete</Button>
          </ButtonToolbar>
        </td>
      </tr>
    );
  }
}

class IndexerList extends Component {
  render() {
    let onSaveIndexer = this.props.onSaveIndexer;
    let indexerNodes = this.props.indexers.map(function(indexer) {
      return (
        <IndexerListRow
          indexer={indexer}
          key={indexer.id}
          onSaveIndexer={onSaveIndexer}
        />
      );
    });
    return (
      <div>
        <Table striped bordered condensed hover>
          <thead>
            <tr>
              <th className="col-md-2">Indexer</th>
              <th className="col-md-6">Feeds</th>
              <th className="col-md-2">State</th>
              <th className="col-md-2">Actions</th>
            </tr>
          </thead>
          <tbody>
            {indexerNodes}
          </tbody>
        </Table>
      </div>
    );
  }
}

export default IndexerList;