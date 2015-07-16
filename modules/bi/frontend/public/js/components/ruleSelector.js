'use strict';

var React = require('react')
var Multiselect = require('react-bootstrap-multiselect');
var _ = require('lodash');

// RuleSelector is a multi selector that lets users choose which rules to display
// on a graph.
var RuleSelector = React.createClass({
	getInitialState: function() {
		return {
			rules: this.props.rules,
			selected: {}
		}
	},
	handleChange: function(element, checked) {
		var newSelectItems = _.extend({}, this.state.selected);
		newSelectItems[element.val()] = checked;
		this.setState({selected: newSelectItems})
		
		// bubble to parent
		this.props.onChange(this);
	},
	render: function() {
		return <Multiselect data={this.props.rules} onChange={this.handleChange} multiple/>
	}
})

module.exports = RuleSelector
