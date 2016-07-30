import React, { Component } from 'react';
import { PageHeader } from 'react-bootstrap';
import './App.css';

import AddIndexer from "./AddIndexer";
import IndexerList from "./IndexerList";
import ConfigModal from './ConfigModal';

class App extends Component {
  static defaultProps = {
    indexers: [],
    enabledIndexers: [],
  }
  state = {
    enabledIndexers: this.props.enabledIndexers,
    configure: false
  }
  isEnabled = (indexer) => {
    let enabled = new Set(this.state.enabledIndexers);
    return enabled.has(indexer.id);
  }
  handleSaveIndexer = (indexer, config) => {
    if (!this.isEnabled(indexer)) {
      this.setState({
        enabledIndexers: this.state.enabledIndexers.concat([indexer.id]),
        configure: false
      });
    }
  }
  handleAddIndexer = (selected) => {
    let handleClose = () => {
      this.setState({configure: false});
    }
    this.setState({
      configure: <ConfigModal key={selected.id}
        indexer={selected} show={true} onSave={this.handleSaveIndexer} onClose={handleClose} />
    })
  }
  render() {
    let isEnabled = this.isEnabled;
    let addableIndexers = this.props.indexers.filter((x) => !isEnabled(x));
    let enabledIndexers = this.props.indexers.filter((x) => isEnabled(x));

    return (
      <div className="App container-fluid">
        <PageHeader>Cardigann <small>Proxy</small></PageHeader>
        <div className="App__body">
          <AddIndexer
            indexers={addableIndexers}
            onIndexerAdd={this.handleAddIndexer} />
          {this.state.configure}
          <IndexerList
            indexers={enabledIndexers}
            onSaveIndexer={this.handleSaveIndexer} />
        </div>
      </div>
    );
  }
}

export default App;
