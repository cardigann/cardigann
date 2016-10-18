import React, { Component } from 'react';
import { Table, ButtonToolbar, Button, Panel } from 'react-bootstrap';
import { OverlayTrigger, Tooltip } from 'react-bootstrap';
import CopyToClipboard from 'react-copy-to-clipboard';
import xhrUrl from './xhr';

function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

class StatefulButton extends Component {
  static defaultProps = {
    bsStyle: "default",
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
      <span className="FeedLink">
        <CopyToClipboard text={this.props.feedHref}>
          <Button bsStyle="default" bsSize="xsmall" title={this.props.feedHref}>Copy {capitalizeFirstLetter(this.props.label)} Feed</Button>
        </CopyToClipboard>{' '}
      </span>
    );
  }
}

class IndexerListRow extends Component {
  static defaultProps = {
    editing: false,
    testing: false,
    disabling: false,
    allowEdit: true,
    allowDisable: true,
    allowTest: true,
    allowSearch: true,
  }
  static propTypes = {
   indexer: React.PropTypes.object.isRequired,
  }
  state = {
    config: {},
    status: "OK",
    editing: this.props.editing,
    testing: this.props.testing,
    disabling: this.props.disabling,
  }
  handleEditClick = () => {
    this.setState({editing: true});
    this.props.onEdit(this.props.indexer, (config) => {
      this.setState({editing: false});
    });
  }
  handleDisableClick = () => {
    this.setState({disabling: true});
    this.props.onDisable(this.props.indexer, (config) => {
      this.setState({disabling: false});
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
  handleSearchClick = () => {
    this.setState({
      searching: true,
    });
    this.props.onSearch(this.props.indexer, () => {
      this.setState({searching: false});
    });
  }
  render() {
    var buttons = [];

    if (this.props.allowEdit) {
      buttons.push(
        <StatefulButton
        key="edit"
        onClick={this.handleEditClick}
        active={this.state.editing}
        activeLabel="Editing..."
        disabled={this.state.testing}>Edit</StatefulButton>
      );
    }

    if (this.props.allowTest) {
      buttons.push(
        <StatefulButton
          key="test"
          onClick={this.handleTestClick}
          active={this.state.testing}
          activeLabel="Testing..."
          disabled={this.state.editing}>Test</StatefulButton>
      );
    }

    if (this.props.allowSearch) {
      buttons.push(
        <StatefulButton
          key="search"
          onClick={this.handleSearchClick}
          active={this.state.searching}
          activeLabel="Searching..."
          disabled={this.state.testing || this.state.editing}>Search</StatefulButton>
      );
    }

    if (this.props.allowDisable) {
      buttons.push(
        <StatefulButton
          key="disable"
          onClick={this.handleDisableClick}
          bsSize="xsmall"
          bsStyle="danger"
          active={this.state.disabling}
          activeLabel="Disabling..."
          disabled={this.state.testing || this.state.editing}>Disable</StatefulButton>
      );
    }

    let tooltip = (
      <Tooltip id="tooltip">
        {this.props.indexer.stats ? this.props.indexer.stats.source : "unknown"}<br />
        {this.props.indexer.stats ? this.props.indexer.stats.modtime : "n/a"}
      </Tooltip>
    );

    return (
      <tr className={this.props.className}>
        <td className="col-md-2">
          <OverlayTrigger trigger="click" placement="right" overlay={tooltip} rootClose={true}>
            <span>{this.props.indexer.name}</span>
          </OverlayTrigger>
        </td>
        <td className="col-md-6">
            {this.props.indexer.feeds.torznab ? <FeedLink
              feedHref={xhrUrl(this.props.indexer.feeds.torznab)}
              label="torznab" /> : ''}
            {this.props.indexer.feeds.potatotorrent ? <FeedLink
              feedHref={xhrUrl(this.props.indexer.feeds.potatotorrent)}
              label="potato" /> : ''}
        </td>
        <td className="col-md-1">{this.state.status}</td>
        <td className="col-md-3">
          <ButtonToolbar>{buttons}</ButtonToolbar>
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
          onSearch={this.props.onSearch}
          onDisable={this.props.onDisable}
        />
      );
    });

    if (indexerNodes.length === 0) {
      return <Panel>No indexers</Panel>;
    }

    var aggregate = {
      "id": "aggregate",
      "name": "All Indexers",
      "enabled": true,
      "feeds": {
        "torznab": xhrUrl("/torznab/aggregate")
      }
    };

    return (
      <div>
        <Table striped bordered condensed hover>
          <thead>
            <tr>
              <th className="col-md-2">Indexer</th>
              <th className="col-md-6">Feeds</th>
              <th className="col-md-1">State</th>
              <th className="col-md-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {indexerNodes}
            <IndexerListRow
              className="all-indexers"
              indexer={aggregate}
              key="aggregate"
              allowEdit={false}
              allowDisable={false}
              allowTest={false}
              onSearch={this.props.onSearch}
            />
          </tbody>
        </Table>
      </div>
    );
  }
}

export default IndexerList;