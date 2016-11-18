var React = require('react');
var Page = require('../components/Page');
var PageHeader = require('../components/PageHeader');
var EntityBreadcrumb = require('../components/EntityBreadcrumb');
var EntityButtonGroup = require('../components/EntityButtonGroup');
var Bootstrap = require('react-bootstrap');
var Row = Bootstrap.Row;
var Col = Bootstrap.Col;

var layer0 = require('../utils/layer0');

var EnvironmentDashboardContainer = React.createClass({
  getInitialState: function() {
    return {
      isLoading: true,
      serviceRows: []
    }
  },
  render: function () {
    environmentID = this.props.routeParams.environmentID

    return (
      <Page>

        <Row>
          <Col lg={12}>
            <PageHeader text='Dashboard' subtext={environmentID} />
            <EntityBreadcrumb environmentID={environmentID} />
          </Col>
        </Row>

        <Row>
          <Col lg={12}>
            <EntityButtonGroup environmentID={environmentID} />
          </Col>
        </Row>

      </Page>
    )
  }
});

module.exports = EnvironmentDashboardContainer;
