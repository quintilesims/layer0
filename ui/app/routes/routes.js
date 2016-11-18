var React = require('react');
var ReactRouter = require('react-router');
var Router = ReactRouter.Router;
var Route = ReactRouter.Route;
var IndexRoute = ReactRouter.IndexRoute;
var Redirect = ReactRouter.Redirect;
var hashHistory = ReactRouter.hashHistory;
var MainContainer = require('../containers/MainContainer');
var HomeDashboardContainer = require('../containers/HomeDashboardContainer');
var EnvironmentDashboardContainer = require('../containers/EnvironmentDashboardContainer');
var ServiceDashboardContainer = require('../containers/ServiceDashboardContainer');
var TaskDashboardContainer = require('../containers/TaskDashboardContainer');
var LoadBalancerDashboardContainer = require('../containers/LoadBalancerDashboardContainer');

var routes = (
  <Router history={hashHistory}>
    <Redirect from='/' to='/dashboard' />
    <Route path='/' component={MainContainer} >
      <Route path='/dashboard' component={HomeDashboardContainer} />
      <Route path='/dashboard/:environmentID' component={EnvironmentDashboardContainer} />
      <Route path='/dashboard/:environmentID/services' component={ServiceDashboardContainer} />
      <Route path='/dashboard/:environmentID/tasks' component={TaskDashboardContainer} />
      <Route path='/dashboard/:environmentID/loadbalancers' component={LoadBalancerDashboardContainer} />
    </Route>
  </Router>
)

module.exports = routes
