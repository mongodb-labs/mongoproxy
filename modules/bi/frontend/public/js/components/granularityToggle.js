'use strict';

var React = require('react')
var Multiselect = require('react-bootstrap-multiselect');
var _ = require('lodash');

// GranularityToggle is a select component that lets users choose the time
// granularity to display on a graph.
var GranularityToggle = React.createClass({
	getInitialState: function() {
		return {
			data: [
				{
					value: "M",
					label: "Month"
				},
				{
					value: "D",
					label: "Day"
				},
				{
					value: "h",
					label: "Hour"
				},
				{
					value: "m",
					label: "Minute",
					selected: true,
				},
				{
					value: "s",
					label: "Second"
				}
			],
			selected: "m"
		};
	},

	handleChange: function(element, checked) {
		
		this.setState({selected: element.val()})
		
		// bubble to parent
		this.props.onChange(this);
	},
	
	render: function() {
		return <Multiselect data={this.state.data} onChange={this.handleChange}/>
	}
});

module.exports = GranularityToggle
