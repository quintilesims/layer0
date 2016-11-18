var React = require('react');
var PropTypes = React.PropTypes;
var PortForm = require('./PortForm');
var Bootstrap = require('react-bootstrap');
var Modal = Bootstrap.Modal;
var Button = Bootstrap.Button;
var FormGroup = Bootstrap.FormGroup;
var ControlLabel = Bootstrap.ControlLabel;
var Form = Bootstrap.Form;
var FormControl = Bootstrap.FormControl;
var HelpBlock = Bootstrap.HelpBlock;
var InputGroup = Bootstrap.InputGroup;
var Checkbox = Bootstrap.Checkbox;
var Radio = Bootstrap.Radio;
var layer0 = require('../utils/layer0');

function newPort() {
  return {
    hostPort: {
      value: 80,
      isValid: true
    },
    containerPort: {
      value: 80,
      isValid: true
    },
    protocol: 'TCP',
    certificateID: {
      value: '',
      isValid: false
    },
  }
}

var CreateLoadBalancerContainer = React.createClass({
  getInitialState: function() {
    return {
      environments: [],
      certificates: [],
      isPublic: true,
      name: {
        value: '',
        isValid: false,
      },
      environmentID: {
        value: '',
        isValid: false,
      },
      ports: [newPort()],
    }
  },
  componentDidMount: function() {
    layer0.listEnvironments()
      .then(function(environments) {
        this.setState({
          environments: environments,
        });
      }.bind(this))

    layer0.listCertificates().then(function(certificates) {
      this.setState({
        certificates: certificates,
      });
    }.bind(this))
  },
  handleNameChange(e) {
    isValid = (e.target.value.length > 0)

    this.setState({
      name: {
        value: e.target.value,
        isValid: isValid
      },
    });
  },
  handleEnvironmentChange(e) {
    this.setState({
      environmentID: {
        value: e.target.value,
        isValid: (e.target.value.length > 0)
      },
    });
  },
  handleAddPort: function() {
    this.state.ports.push(newPort())

    this.setState({
      ports: this.state.ports
    });
  },
  handleDeletePort: function(index) {
    this.state.ports.splice(index, 1)

    this.setState({
      ports: this.state.ports
    });
  },
  handlePortChange: function(index, port) {
    this.state.ports[index] = port

    this.setState({
      ports: this.state.ports,
    });
  },
  getPorts: function() {
    ports = this.state.ports.map(function(port, i) {
      return (
        <PortForm key={i} 
	  onChange={this.handlePortChange.bind(this, i)} 
	  certificates={this.state.certificates}
          hostPort={port.hostPort} 
	  containerPort={port.containerPort} 
	  protocol={port.protocol}
          certificateID={port.certificateID}>

	  {' '}
	  { i > 0 ? <Button onClick={this.handleDeletePort.bind(this, i)}>Delete</Button> : null } 
	  {' '}
	  { i == this.state.ports.length-1 ? <Button onClick={this.handleAddPort}>+</Button> : null }
	  </PortForm>
      )
    }.bind(this))

    return ports
  },
  isFormValid: function() {
    if (!this.state.name.isValid) {
      return false
    }

    if (!this.state.environmentID.isValid) {
      return false
    }

    for (i = 0; i < this.state.ports.length; i++) {
      port = this.state.ports[i]

      if (!port.hostPort.isValid) {
        return false
      }

      if (!port.containerPort.isValid) {
        return false
      }

      if (port.protocol == 'SSL' || port.protocol == 'HTTPS') {
        if (!port.certificateID.isValid) {
          return false
        }
      }
    }

    return true
  },
  handleSubmit: function(e) {
    e.preventDefault()

    ports = this.state.ports.map(function(port) {
      return {
        host_port: port.hostPort.value,
        container_port: port.containerPort.value,
        protocol: port.protocol,
        certificate_id: port.certificateID.value,
      }
    })

    request = {
      environment_id: this.state.environmentID.value,
      load_balancer_name: this.state.name.value,
      is_public: this.state.isPublic,
      ports: ports,
    }

    console.log(request)

    return layer0.createLoadBalancer(request)
      .then(function(response) {
        console.log(response)
        this.props.onHide()
      }.bind(this))
      .catch(function(err) {
        alert(err.data.message)
        console.log(err)
      })
  },
  render: function() {
    return (
      <Modal bsSize='large' show={this.props.show} onHide={this.props.onHide}>
  <form onSubmit={this.handleSubmit}>
    <Modal.Header closeButton>
      <Modal.Title>Create Load Balancer</Modal.Title>
    </Modal.Header>
    <Modal.Body>

      <FormGroup 
      validationState={this.state.environmentID.isValid ? 'success' : 'error' }>
        <ControlLabel>Environment</ControlLabel>
        <FormControl componentClass='select' onChange={this.handleEnvironmentChange}>
          <option value='' >Select Environment...</option>
          {this.state.environments.map(function(environment, i){
            return (
              <option key={i} value={environment.environment_id}>
              {environment.environment_id}
              </option>
             ) 
	  })}
 	  </FormControl>
  	</FormGroup>

      <FormGroup validationState={this.state.name.isValid ? 'success' : 'error' }>
        <ControlLabel>Name</ControlLabel>
        <FormControl type='text' value={this.state.name.value} placeholder='Enter name'
          onChange={this.handleNameChange} />
        <FormControl.Feedback />
      </FormGroup>

      <ControlLabel>Accessiblity</ControlLabel>
      <FormGroup>
        <Radio inline checked={this.state.isPublic} onChange={()=> this.setState({ isPublic: true }) }>Public</Radio>
        <Radio inline checked={!this.state.isPublic} onChange={()=> this.setState({ isPublic: false }) }>Private</Radio>
      </FormGroup>

      <ControlLabel>Ports</ControlLabel>
      {this.getPorts()}

    </Modal.Body>
    <Modal.Footer>
      <Button bsStyle='primary' type='submit' disabled={!this.isFormValid()}>Create</Button>
      <Button onClick={this.props.onHide}>Cancel</Button>
    </Modal.Footer>
  </form>
</Modal>
    )
  }
});
module.exports = CreateLoadBalancerContainer;
