var React = require('react');
var ReactRouterBootstrap = require('react-router-bootstrap');
var LinkContainer = ReactRouterBootstrap.LinkContainer;
var Bootstrap = require('react-bootstrap');
var Navbar = Bootstrap.Navbar;
var NavItem = Bootstrap.NavItem;
var NavDropdown = Bootstrap.NavDropdown;
var Nav = Bootstrap.Nav;
var MenuItem = Bootstrap.MenuItem;

var layer0 = require('../utils/layer0');

var TopNavContainer = React.createClass({
  getInitialState: function() {
    return {
      isLoading: true,
      environments: []
    }
  },
  componentDidMount: function() {
   layer0.listEnvironments()
     .then(function(environments){
       this.setState({
         isLoading: false,
         environments: environments
       });
      }.bind(this))
  },
  getEnvironmentMenuItems: function() {
    if (this.state.isLoading) {
      return <MenuItem disabled>Loading...</MenuItem>
    }

    return this.state.environments.map(function(environment, i){
      environmentID = environment.environment_id

      return (
        <LinkContainer key={i} to={'/dashboard/'+environmentID}>
          <MenuItem>{environmentID}</MenuItem>
        </LinkContainer>
      )
    }.bind(this))
  },
  render: function () {
    return (
      <Navbar>
        <Navbar.Header>
          <Navbar.Brand>
            <a href='#/'>Layer0</a>
          </Navbar.Brand>
          <Navbar.Toggle />
        </Navbar.Header>
        <Navbar.Collapse>
        <Nav>
          <NavItem href="http://docs.xfra.ims.io">Docs</NavItem>
          <NavItem href="https://layer0.xfra.ims.io/apidocs">API</NavItem>
          <NavDropdown eventKey={3} title="Environments" id="basic-nav-dropdown">
            {this.getEnvironmentMenuItems()}
            <MenuItem divider />
            <LinkContainer to='#'>
              <MenuItem>Create New</MenuItem>
            </LinkContainer>
          </NavDropdown>
       </Nav>
       <Nav pullRight>
         <LinkContainer to='/admin'>
           <NavItem>Admin</NavItem>
         </LinkContainer>
       </Nav>
     </Navbar.Collapse>
   </Navbar>
    )
  }
});

module.exports = TopNavContainer;
