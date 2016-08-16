import React, { Component } from 'react';
import { PageHeader, Button, Glyphicon } from 'react-bootstrap';
import CopyToClipboard from 'react-copy-to-clipboard';
import './App.css';

import AddIndexer from "./AddIndexer";
import IndexerList from "./IndexerList";
import ConfigModal from "./ConfigModal";
import AlertDismissable from "./AlertDismissable";
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
    errorMessage: false
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
    .catch((err) => {
      this.setState({errorMessage: err})
    });
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
        afterFunc(true);
      } else {
        this.setState({errorMessage: data.error})
        afterFunc(false, data.error);
      }
    })
    .catch((err) => {
      this.setState({errorMessage: err})
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
    .then((response) => {
      if (!response.ok) {
        return response.json().then((resp) => {
          throw Error(resp.error);
        });
      }
      return response.json()
    })
    .then((indexers) => {
      this.setState({
        indexers: indexers,
        enabledIndexers: indexers.filter((x) => x.enabled).map((x) => x.id),
      })
    })
    .catch((err) => {
      this.setState({errorMessage: err.message, errorScope: "loading indexers"})
    });
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
    let errorAlert = null;

    if (this.state.apiKey === null) {
      return <Login onAuthenticate={this.handleAuthenticate} />;
    }

    if (this.state.errorMessage) {
      errorAlert = <AlertDismissable>
        <h4>An error occurred {this.state.errorScope ? "whilst " + this.state.errorScope : ""}</h4>
        <p>{this.state.errorMessage}</p>
      </AlertDismissable>;
    }

    return (
      <div className="App container-fluid">
        <PageHeader>Cardigann <small>Proxy</small></PageHeader>
        {errorAlert}
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
