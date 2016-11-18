var React = require('react');

var Page = React.createClass({
  render: function(){
    return (
      <div id='page-wrapper' >
        <div className='container-fluid' >
          {this.props.children}
        </div>
      </div>
    )
  }
});

module.exports = Page;

