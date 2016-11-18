var React = require('react');
var Page = require('../components/Page');
var Bootstrap = require('react-bootstrap');
var Jumbotron = Bootstrap.Jumbotron;
var Button = Bootstrap.Button;

var HomeDashboardContainer = React.createClass({
  render: function () {
    return (
      <Page>
        <Jumbotron>
          <h1>Hello, world!</h1>
          <p>
            This is a template for a simple marketing or informational website. 
            It includes a large callout called a jumbotron and three supporting pieces of content. 
            Use it as a starting point to create something more unique.
          </p>
          <p><Button href='http://docs.xfra.ims.io' bsStyle='primary'>Learn more</Button></p>
        </Jumbotron>
      </Page>
    )
  }
});

module.exports = HomeDashboardContainer;
