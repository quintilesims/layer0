var axios = require('axios');


// todo: not sure what best practice is in regards to updating these values
// via the admin page. Perhaps we export Getter/Setter functions 
const ROOT_URL = 'https://l0-zpdev-api-122190750.us-west-2.elb.amazonaws.com';
const AUTH_TOKEN = 'bGF5ZXIwOkFGWm40TTFRU0g=';

axios.defaults.baseURL = localStorage.getItem('endpoint') || ROOT_URL;
axios.defaults.headers.common['Authorization'] = 'Basic ' + (localStorage.getItem('token') || AUTH_TOKEN);

function createLoadBalancer(request) {
  return axios.post('/loadbalancer', request)
  .then(function(response) {
    return response.data
  })
}

function listCertificates() {
  return axios.get('/certificate').
  then(function(response) {
      return response.data
    })
    .catch(function(err) {
      handleError('List Certificates', err)
    });
}


function listEnvironments() {
  return axios.get('/environment').
  then(function(response) {
      return response.data
    })
    .catch(function(err) {
      handleError('List Environments', err)
    });
}

function listLoadBalancers() {
  return axios.get('/loadbalancer').
  then(function(response) {
      return response.data
    })
    .catch(function(err) {
      handleError('List LoadBalancers', err)
    });
}

function listServices() {
  return axios.get('/service').
  then(function(response) {
      return response.data
    })
    .catch(function(err) {
      handleError('List Services', err)
    });
}

function listTasks() {
  return axios.get('/task').
  then(function(response) {
      return response.data
    })
    .catch(function(err) {
      handleError('List Tasks', err)
    });
}

function filterByEnvironment(entities, environmentID) {
  return entities.filter(function(entity) {
    return entity.environment_id == environmentID;
  });
}

function handleError(name, err) {
  alert('An unexpected error occured when attempting ' + name + ': \n' + err.data);
  console.warn('Error in ' + name + '! =>', err);
}

var layer0 = {
  setEndpoint: function(url) {
    localStorage.setItem('endpoint', url)
    axios.defaults.baseURL = url;
  },
  getEndpoint: function() {
    return axios.defaults.baseURL;
  },
  setToken: function(token) {
    localStorage.setItem('token', token);
    axios.defaults.headers.common['Authorization'] = 'Basic ' + token;
  },
  getToken: function() {
    return axios.defaults.headers.common['Authorization'].replace('Basic ', '');
  },
  createLoadBalancer: createLoadBalancer,
  filterByEnvironment: filterByEnvironment,
  listCertificates: listCertificates,
  listEnvironments: listEnvironments,
  listLoadBalancers: listLoadBalancers,
  listServices: listServices,
  listTasks: listTasks,
}

module.exports = layer0
