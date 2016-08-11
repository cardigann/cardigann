import React, { Component } from 'react';
import { Col, Form, Panel, FormGroup, Button } from 'react-bootstrap';
import Select from 'react-select';

let buildOptions = function(indexers) {
  return indexers.map((indexer) => {
    return { value: indexer.id, label: indexer.name };
  });
}

class AddIndexer extends Component {
  static defaultProps = {
    indexers: [],
  }
  static propTypes = {
    indexers: React.PropTypes.array,
    onAdd: React.PropTypes.func,
  }
  state = {
    selected: false,
    options: buildOptions(this.props.indexers)
  }
  componentWillReceiveProps(newProps) {
    if (typeof(newProps.indexers) !== undefined) {
      this.setState({options: buildOptions(newProps.indexers)});
    }
  }
  handleSubmit = (e) => {
    e.preventDefault();
    if (this.state.selected === false) {
      return;
    }
    this.props.onAdd(this.state.selected, {});
    this.setState({selected: false});
  }
  handleSelectChange = (opt) => {
    this.setState({selected: this.props.indexers.filter((x) => x.id === opt.value)[0]});
  }
  render() {
    return (
      <div className="AddIndexer">
        <Panel header="Add Indexer">
          <Form horizontal onSubmit={this.handleSubmit}>
            <FormGroup controlId="formControlsSelect">
              <Col xs={12} md={4}>
                <Select name="form-field-name"
                  value={this.state.selected.id}
                  options={this.state.options}
                  onChange={this.handleSelectChange} />
              </Col>
              <Button type="submit" bsSize="small">Add</Button>
            </FormGroup>
          </Form>
        </Panel>
      </div>
    );
  }
}

export default AddIndexer;