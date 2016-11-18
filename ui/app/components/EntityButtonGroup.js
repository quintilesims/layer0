var React = require('react');
var PropTypes = React.PropTypes;
var ReactRouterBootstrap = require('react-router-bootstrap');
var LinkContainer = ReactRouterBootstrap.LinkContainer;
var Bootstrap = require('react-bootstrap');
var Button = Bootstrap.Button;
var ButtonGroup = Bootstrap.ButtonGroup;

var EntityButtonGroup = React.createClass({
  propTypes: {
    environmentID: PropTypes.string.isRequired,
    services: PropTypes.bool,
    tasks: PropTypes.bool,
    loadBalancers: PropTypes.bool
  },
  getLink: function(entity) {
    return '#/dashboard/' + this.props.environmentID + '/' + entity
  },
  render: function() {
    return (
      <ButtonGroup justified>
        <Button href={this.getLink('services')} active={this.props.services}>Services</Button>
        <Button href={this.getLink('tasks')} active={this.props.tasks}>Tasks</Button>
        <Button href={this.getLink('loadBalancers')} active={this.props.loadBalancers}>Load Balancers</Button>
      </ButtonGroup>
    )
  },
});

module.exports = EntityButtonGroup;
