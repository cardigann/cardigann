import React, { Component } from 'react';
import { Table, ButtonToolbar, Button, Panel } from 'react-bootstrap';

class StatefulButton extends Component {
  static defaultProps = {
    bsStyle: "primary",
    bsSize: "xsmall",
    activeLabel: "Saving...",
    onClick: (e) => {},
    disabled: false,
    active: false,
  }
  state = {
    active: this.props.active,
  }
  componentWillReceiveProps(newProps) {
    this.setState({
      active: typeof(newProps.active) !== undefined ? newProps.active : this.state.active,
    });
  }
  handleClick = (e) => {
    this.setState({active: true});
    this.props.onClick(e);
  }
  render() {
    return (
      <Button
        bsStyle={this.props.bsStyle}
        bsSize={this.props.bsSize}
        disabled={this.props.disabled || this.state.active}
        onClick={!this.state.active ? this.handleClick : null}>
        {this.state.active ? this.props.activeLabel : this.props.children}
      </Button>
    );
  }
}

class FeedLink extends Component {
  render() {
    return (
      <span className="FeedLink">{this.props.feedHref}</span>
    );
  }
}

class IndexerListRow extends Component {
  static defaultProps = {
    editing: false,
    testing: false,
  }
  static propTypes = {
   indexer: React.PropTypes.object.isRequired,
  }
  state = {
    config: {},
    status: "OK",
    editing: this.props.editing,
    testing: this.props.testing,
  }
  handleEditClick = () => {
    this.setState({editing: true});
    this.props.onEdit(this.props.indexer, (config) => {
      this.setState({editing: false});
    });
  }
  handleTestClick = () => {
    this.setState({
      status: "Testing",
      testing: true,
    });
    this.props.onTest(this.props.indexer, (ok) => {
      this.setState({
        status: ok ? "OK" : "Failed",
        testing: false,
      });
    })
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
        <td className="col-md-2">{this.state.status}</td>
        <td className="col-md-2">
          <ButtonToolbar>
            <StatefulButton
              onClick={this.handleEditClick}
              active={this.state.editing}
              activeLabel="Editing..."
              disabled={this.state.testing}>Edit</StatefulButton>
            <StatefulButton
              onClick={this.handleTestClick}
              active={this.state.testing}
              activeLabel="Testing..."
              disabled={this.state.editing}>Test</StatefulButton>
            <StatefulButton
              bsSize="xsmall"
              bsStyle="danger"
              disabled={this.state.testing || this.state.editing}>Disable</StatefulButton>
          </ButtonToolbar>
        </td>
      </tr>
    );
  }
}

class IndexerList extends Component {
  render() {
    let indexerNodes = this.props.indexers.map((indexer) => {
      return (
        <IndexerListRow
          indexer={indexer}
          key={indexer.id}
          onSave={this.props.onSave}
          onEdit={this.props.onEdit}
          onTest={this.props.onTest}
        />
      );
    });

    if (indexerNodes.length === 0) {
      return <Panel>No indexers</Panel>;
    }

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