var React = require('react');
var CreateLoadBalancerContainer = require('./CreateLoadBalancerContainer');
var Page = require('../components/Page');
var PageHeader = require('../components/PageHeader');
var EntityBreadcrumb = require('../components/EntityBreadcrumb');
var EntityButtonGroup = require('../components/EntityButtonGroup');
var EntityTable = require('../components/EntityTable');
var Bootstrap = require('react-bootstrap');
var Row = Bootstrap.Row;
var Col = Bootstrap.Col;
var ButtonGroup = Bootstrap.ButtonGroup;
var DropdownButton = Bootstrap.DropdownButton;
var MenuItem = Bootstrap.MenuItem;
var Button = Bootstrap.Button;
var Panel = Bootstrap.Panel;
var Loader = require('react-loader');

var layer0 = require('../utils/layer0');

var LoadBalancerDashboardContainer = React.createClass({
  getInitialState: function() {
    return {
      isLoading: true,
      loadBalancers: [],
      showCreateModal: false,
      selectedRow: -1,
    }
  },
  componentDidMount: function() {
    layer0.listLoadBalancers()
      .then(function(loadBalancers) {
        return layer0.filterByEnvironment(loadBalancers, this.props.routeParams.environmentID)
      }.bind(this))
      .then(function(loadBalancers) {
        loadBalancers = loadBalancers.map(function(loadBalancer) {
          ports = loadBalancer.ports.map(function(port) {
            return port.host_port + ':' + port.container_port + '/' + port.protocol
          })

          return {
            id: loadBalancer.load_balancer_id,
            fields: [
              loadBalancer.load_balancer_id,
              loadBalancer.load_balancer_name,
              loadBalancer.url,
              loadBalancer.is_public ? 'True' : 'False',
              ports,
            ]
          }
        })

        this.setState({
          isLoading: false,
          loadBalancers: loadBalancers
        });
      }.bind(this))
  },
  openCreateModal: function() {
    this.setState({
      showCreateModal: true
    });
  },
  closeCreateModal: function() {
    this.setState({
      showCreateModal: false
    });
  },
  handleRowSelect: function(row){
    this.setState({
      selectedRow: row,
    });
  },
  render: function() {
    environmentID = this.props.routeParams.environmentID
    return (
      <Page>
        <Row>
          <Col lg={12}>
            <PageHeader text='Dashboard' subtext={environmentID} />
            <EntityBreadcrumb environmentID={environmentID} loadBalancers />
          </Col>
        </Row>
        <Row>
          <Col lg={12}>
            <EntityButtonGroup environmentID={environmentID} loadBalancers />
          </Col>
        </Row>
        <CreateLoadBalancerContainer 
          show={this.state.showCreateModal}
          onHide={this.closeCreateModal} 
        />
        <Loader loaded={!this.state.isLoading}>
          <Row>
            <Col lg={12}>
              <Panel>
                <ButtonGroup>
                  <Button onClick={this.openCreateModal}>Create New</Button>
                  <DropdownButton title='Actions' id='bg-nested-dropdown'>
                    <MenuItem eventKey='delete'>Delete</MenuItem>
                    <MenuItem disabled={this.state.selectedRow == -1}>Update Ports</MenuItem>
                  </DropdownButton>
                </ButtonGroup>
                <hr />
                <EntityTable 
                  headers={['ID', 'Name', 'URL', 'Public', 'Ports']}
                  entities={this.state.loadBalancers} 
                  onRowSelect={this.handleRowSelect}
                />
              </Panel>
            </Col>
          </Row>    
        </Loader>
      </Page>
    )
  }
});

module.exports = LoadBalancerDashboardContainer;
