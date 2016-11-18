var React = require('react');
var PropTypes = React.PropTypes;
var Bootstrap = require('react-bootstrap');
var BSPageHeader = Bootstrap.PageHeader;

var PageHeader = React.createClass({
  propTypes: {
    text: PropTypes.string.isRequired,
    subtext: PropTypes.string
  },
  getDefaultProps: function() {
    return {
      subtext: ''
    }
  },
  render: function() {
    return (
      <BSPageHeader>{this.props.text} <small>{this.props.subtext}</small></BSPageHeader>
    )
  }
});

module.exports = PageHeader;
