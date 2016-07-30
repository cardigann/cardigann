import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';

var indexers = [
  {
    id: "bithdtv",
    name: "BIT-HDTV",
    feeds: {
      "torznab": window.location.href + "torznab/bithdtv/api"
    }
  },
  {
    id: "example",
    name: "Example",
    feeds: {
      "torznab": window.location.href + "torznab/example/api"
    }
  }
];

var enabledIndexers = [
  "example"
]

ReactDOM.render(
  <App indexers={indexers} enabledIndexers={enabledIndexers} />,
  document.getElementById('root')
);
