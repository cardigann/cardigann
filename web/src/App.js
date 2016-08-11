import React, { Component } from 'react';
import { PageHeader, Button, Glyphicon } from 'react-bootstrap';
import CopyToClipboard from 'react-copy-to-clipboard';
import './App.css';

import AddIndexer from "./AddIndexer";
import IndexerList from "./IndexerList";
import ConfigModal from "./ConfigModal";
import Login from './Login';

class App extends Component {
  static defaultProps = {
    indexers: [],
    enabledIndexers: [],
    apiKey: localStorage.getItem('apiKey'),
  }
  state = {
    indexers: this.props.indexers,
    enabledIndexers: this.props.enabledIndexers,
    configure: null,
    apiKey: this.props.apiKey,
    apiKeyCopied: false,
  }
  isEnabled = (indexer) => {
    return this.state.enabledIndexers.filter((x) => x === indexer.id).length > 0;
  }
  handleSaveIndexer = (indexer, config, afterFunc) => {
    if (!this.isEnabled(indexer)) {
      this.setState({
        enabledIndexers: this.state.enabledIndexers.concat([indexer.id]),
        configure: false
      });
    }
    fetch("/xhr/indexers/"+indexer.id+"/config", {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'apitoken ' + this.state.apiKey,
        },
        method: "PATCH",
        body: JSON.stringify(config),
    })
    .then(function(res){
      if(res.ok) {
        afterFunc();
      } else {
        console.log("error response");
      }
    })
    .catch(function(res){
      console.error(res);
    })
  }
  handleAddIndexer = (selected) => {
    this.loadIndexerConfig(selected, (config) => {
      this.showConfigModal(selected, config);
    });
  }
  handleEditIndexer = (selected, afterFunc) => {
    this.loadIndexerConfig(selected, (config) => {
      this.showConfigModal(selected, config, afterFunc);
    });
  }
  handleTestIndexer = (indexer, afterFunc) => {
    fetch("/xhr/indexers/"+indexer.id+"/test", {
        headers: {
          'Accept': 'application/json',
          'Authorization': 'apitoken ' + this.state.apiKey,
        },
        method: "POST",
    })
    .then((response) => response.json())
    .then(function(data){
      if(data.ok) {
        console.log("ok response");
        afterFunc(true);
      } else {
        console.log("error response", data.error);
        afterFunc(false, data.error);
      }
    })
    .catch(function(res){
      console.error(res);
    });
  }
  handleAuthenticate = (apiKey) => {
    localStorage.setItem("apiKey", apiKey);
    this.setState({apiKey: apiKey}, () => {
      this.loadIndexers();
    });
  }
  loadIndexerConfig = (indexer, dataFunc) => {
    fetch("/xhr/indexers/"+indexer.id+"/config", {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'apitoken ' + this.state.apiKey,
        },
    })
    .then((response) => response.json())
    .then(dataFunc)
  }
  loadIndexers = () => {
    if (!this.state.apiKey) {
      console.error("No api key is set");
      return;
    }
    fetch("/xhr/indexers", {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'apitoken ' + this.state.apiKey,
        },
    })
    .then((response) => response.json())
    .then(function(indexers) {
      this.setState({
        indexers: indexers,
        enabledIndexers: indexers.filter((x) => x.enabled).map((x) => x.id),
      })
    }.bind(this));
  }
  showConfigModal = (indexer, config, afterFunc) => {
    this.setState({
      configure: <ConfigModal config={config} indexer={indexer} show={true}
        onClose={() => {
          this.setState({configure: null});
          afterFunc();
        }}
        onSave={(indexer, config, afterSaveFunc) => {
          this.handleSaveIndexer(indexer, config, () => { afterFunc(); afterSaveFunc() });
        }}
      />
    });
  }
  componentWillMount() {
    if (this.state.apiKey) {
      this.loadIndexers();
    }
  }
  render() {
    let isEnabled = this.isEnabled;
    let addableIndexers = this.state.indexers.filter((x) => !isEnabled(x));
    let enabledIndexers = this.state.indexers.filter((x) => isEnabled(x));

    if (this.state.apiKey === null) {
      return <Login onAuthenticate={this.handleAuthenticate} />;
    }

    return (
      <div className="App container-fluid">
        <PageHeader>Cardigann <small>Proxy</small></PageHeader>
        <div className="App__apiKey">
          <strong>API Key: </strong>
          <code>{this.state.apiKey}</code>
          <CopyToClipboard text={this.state.apiKey} onCopy={() => this.setState({apiKeyCopied: true})}>
            <Button bsSize="xsmall">Copy <Glyphicon glyph="copy" /></Button>
          </CopyToClipboard>
          {this.state.apiKeyCopied ? <span className="copied">Copied.</span> : null}
        </div>
        <div className="App__body">
          <AddIndexer
            indexers={addableIndexers}
            onAdd={this.handleAddIndexer} />
          <IndexerList
            indexers={enabledIndexers}
            onEdit={this.handleEditIndexer}
            onSave={this.handleSaveIndexer}
            onTest={this.handleTestIndexer} />
          {this.state.configure}
        </div>
      </div>
    );
  }
}

export default App;
