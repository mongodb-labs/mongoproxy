'use strict';

var React = require('react');
var Button = require('react-bootstrap').Button;

var AddRuleButton = React.createClass({
	
	render: function() {
		return <Button onClick={this.props.onClick} block bsSize="large">Add Rule</Button>
	}
})

module.exports = AddRuleButton
