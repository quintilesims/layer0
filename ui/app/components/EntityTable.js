var React = require('react');
var PropTypes = React.PropTypes;
var Bootstrap = require('react-bootstrap');
var Table = Bootstrap.Table;
var Radio = Bootstrap.Radio;

var EntityTable = React.createClass({
  propTypes: {
    headers: PropTypes.array.isRequired,
    rows: PropTypes.array,
    onRowSelect: PropTypes.func,
  },
  getDefaultProps: function() {
    return {
      entities: [],
      onRowSelect: function(){},
    }
  },
  getCell: function(field, i) {
    if (Array.isArray(field)) {
      return (
        <td key={i}>
          {field.map(function(f, j){
            return <p key={j}>{f}</p>
          })}
        </td>
      )
    }

    return <td key={i}>{field}</td>
  },
  handleRowSelect: function(i){
    this.props.onRowSelect(i)
  },
  getRow: function(entity, i) {
    return (
      <tr key={i}>
        <td><Radio name='selected' onChange={this.handleRowSelect.bind(this, i)}/></td>
        {entity.fields.map(function(field, j){
          return this.getCell(field, j)
        }.bind(this))}
      </tr>
    )
  },
  render: function() {
    return (
      <Table responsive striped condensed hover>
        <thead>
          <tr>
            <th>{' '}</th>
            {this.props.headers.map(function(header, i){
              return <th key={i}>{header}</th>
            })}
          </tr>
        </thead>
        <tbody>
          {this.props.entities.map(function(entity, i){
            return this.getRow(entity, i)
          }.bind(this))}
        </tbody>
      </Table>
    )
  }
});

module.exports = EntityTable;
