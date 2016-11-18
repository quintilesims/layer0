var React = require('react');
var TopNavContainer = require('./TopNavContainer');

var MainContainer = React.createClass({
  render: function () {
    return (
      <div id='wrapper'>
        <TopNavContainer />
        {this.props.children}
      </div>
    )
  }
});

module.exports = MainContainer;
