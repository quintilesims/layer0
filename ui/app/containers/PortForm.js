var React = require('react');
var PropTypes = React.PropTypes;
var Bootstrap =
  require('react-bootstrap');
var Modal = Bootstrap.Modal;
var Button = Bootstrap.Button;
var FormGroup = Bootstrap.FormGroup;
var ControlLabel = Bootstrap.ControlLabel;
var Form = Bootstrap.Form;
var FormControl = Bootstrap.FormControl;
var HelpBlock = Bootstrap.HelpBlock;
var InputGroup = Bootstrap.InputGroup;
var OverlayTrigger = Bootstrap.OverlayTrigger;
var Tooltip = Bootstrap.Tooltip;
var layer0 = require('../utils/layer0');

const MIN_PORT = 1
const MAX_PORT = 65535

function isInt(value) {
  return !isNaN(value) && parseInt(Number(value)) == value &&
    !isNaN(parseInt(value, 10));
}

const hostPortTooltip = (
  <Tooltip id='host'>Host Port</Tooltip>
);

const containerPortTooltip = (
  <Tooltip id='container'>Container Port</Tooltip>
);

var CreateLoadBalancerContainer = React.createClass({
  PropTypes: {
    onChange: PropTypes.function,
    certificates: PropTypes.array.isRequired,
    hostPort: PropTypes.number.isRequired,
    containerPort: PropTypes.number.isRequired,
    protocol: PropTypes.string.isRequired,
    certificateID: PropTypes.string,
  },
  getCertificate: function() {
    if (this.props.protocol == 'SSL' || this.props.protocol == 'HTTPS') {
      return (
        <FormGroup>
        <FormControl componentClass='select' onChange={this.handleCertificateChange}>
          <option value=''>Select Certificate...</option>
	  {this.props.certificates.map(function(certificate, i){
            return (
	      <option key={i}
              value={certificate.certificate_id}>{certificate.certificate_id}
	      </option>
	      )
          })}
        </FormControl>
      </FormGroup>
      )
    }
  },
  handlePortChange(isHostPort, e) {
    value = e.target.value
    isValid = isInt(value)

    if (isValid) {
      value = parseInt(e.target.value)
   
      if (value < MIN_PORT) {
	      value = MIN_PORT
      }

      if (value > MAX_PORT) {
	value = MAX_PORT
      }
    }

    port = {
      hostPort: isHostPort ? {
        isValid: isValid,
        value: value,
      } : this.props.hostPort,
      containerPort: isHostPort ? this.props.containerPort : {
        isValid: isValid,
        value: value
      },
      protocol: this.props.protocol,
      certificateID: this.props.certificateID,
    }

    this.props.onChange(port)
  },
  handleProtocolChange(e){
    port = {
      hostPort: this.props.hostPort,
      containerPort: this.props.containerPort,
      protocol: e.target.value,
      certificateID: this.props.certificateID,
    }

     this.props.onChange(port)
  },
  handleCertificateChange(e){
    port = {
      hostPort: this.props.hostPort,
      containerPort: this.props.containerPort,
      protocol: this.props.protocol,
      certificateID: { value: e.target.value, isValid: (e.target.value.length > 0) },
    }

    this.props.onChange(port)
  },
  render: function() {
    return (
      <Form componentClass='fieldset' inline>
    <FormGroup validationState={ this.props.hostPort.isValid ? 'success' : 'error' }>
      <OverlayTrigger placement="bottom" overlay={hostPortTooltip}>
        <InputGroup>
          <FormControl 
	    type='number' 
	    value={this.props.hostPort.value} 
	    min={MIN_PORT} 
	    max={MAX_PORT} 
	    onChange={this.handlePortChange.bind(this, true)} />
          <FormControl.Feedback />
        </InputGroup>
      </OverlayTrigger>
    </FormGroup>
    {' : '}
    <FormGroup validationState={this.props.containerPort.isValid ? 'success' : 'error' }>
      <OverlayTrigger placement="bottom" overlay={containerPortTooltip}>
        <InputGroup>
          <FormControl 
	    type='number' 
	    value={this.props.containerPort.value} 
            min={MIN_PORT} 
	    max={MAX_PORT}
            onChange={this.handlePortChange.bind(this, false)} />
          <FormControl.Feedback />
        </InputGroup>
      </OverlayTrigger>
    </FormGroup>
    {' / '}
    <FormGroup>
      <FormControl componentClass='select' onChange={this.handleProtocolChange}>
        <option value='TCP'>TCP</option>
        <option value='SSL'>SSL</option>
        <option value='HTTP'>HTTP</option>
        <option value='HTTPS'>HTTPS</option>
      </FormControl>
    </FormGroup>
    {' '} { this.getCertificate() }
    {this.props.children}
  </Form>
    )
  }
});
module.exports = CreateLoadBalancerContainer;
