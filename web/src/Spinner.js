import React, { Component } from 'react';
import Halogen from 'halogen';

class Spinner extends Component {
  static defaultProps = {
    color: "#000000",
  }
  render() {
    var style = {
      display: 'inline-block',
      height: '45px',
      width: '45px',
      verticalAlign: 'center',
      align: 'center',
      marginTop: '10px',
    };

    return (
      <div style={style}><Halogen.ClipLoader color={this.props.color}/></div>
    );
  }
}


export default Spinner;