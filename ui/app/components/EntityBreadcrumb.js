var React = require('react');
var PropTypes = React.PropTypes;
var ReactRouterBootstrap = require('react-router-bootstrap');
var LinkContainer = ReactRouterBootstrap.LinkContainer;
var Bootstrap = require('react-bootstrap');
var Breadcrumb = Bootstrap.Breadcrumb;
var Glyphicon = Bootstrap.Glyphicon;

var EntityBreadcrumb = React.createClass({
  propTypes: {
    environmentID: PropTypes.string.isRequired,
    services: PropTypes.bool,
    tasks: PropTypes.bool,
    loadBalancers: PropTypes.bool
  },
  getDashboardItem: function() {
    if (this.props.services || this.props.tasks || this.props.loadBalancers){
      return (
        <Breadcrumb.Item 
          href={'#/dashboard/'+this.props.environmentID}>
            <Glyphicon glyph='dashboard'/> Dashboard
        </Breadcrumb.Item>
      )
    }

    return <Breadcrumb.Item active><Glyphicon glyph='dashboard'/> Dashboard</Breadcrumb.Item>
  },
  render: function() {
    return (
          <Breadcrumb>
            {this.getDashboardItem()}
            {this.props.services ? <Breadcrumb.Item active><Glyphicon glyph='th'/> Services</Breadcrumb.Item> : null }
            {this.props.tasks ? <Breadcrumb.Item active><Glyphicon glyph='tasks'/> Tasks</Breadcrumb.Item> : null }
            {this.props.loadBalancers ? <Breadcrumb.Item active><Glyphicon glyph='th-list'/> Load Balancers</Breadcrumb.Item> : null }
         </Breadcrumb>
        )
  },
});

module.exports = EntityBreadcrumb
