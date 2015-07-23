'use strict';

var React = require('react');
var Panel = require('react-bootstrap').Panel;
var JSONEditor = require('react-json');

var Button = require('react-bootstrap').Button;

var _ = require('lodash');

var tg = require('../utils/convertTimeGranularity');

// RuleConfig is an editor for an individual rule in the BI Module configuration.
var RuleConfig = React.createClass({
	getInitialState: function() {
		var defaultRule = {
				origin: "db.originCollection",
				prefix: "db.metricCollection",
				timeGranularity: {
					Month: true,
					Day: true,
					Hour: true,
					Minute: true,
					Second: true
				},
				valueField: "fieldName",
				timeField: ""
			}
		return {
			defaultRule: defaultRule,
			rule: _.extend(defaultRule,
				tg.convertToBooleans(this.props.rule)),
			settings: {
				form: true,
				fields: {
					origin: {
						type: 'string',
						placeholder: "db.originCollection",
						editing: true
					},
					timeField: {
						type: 'string',
						placeholder: "optional",
						editing: 'always'
					}
				}
			}
		}
	},
	componentDidMount: function() {
		
		this.props.onChange(this, this.state.rule)
	},
	handleClick: function() {

		this.props.onClick(this.props.key);
		// bubble up to parent
	},
	getRule: function() {
		return this.state.rule;
	},
	handleChange: function(value) {
		this.setState({
			rule: value
		});

		// bubble up to parent
		this.props.onChange(this, value);
	},
	render: function() {
		return <Panel>
			<JSONEditor onChange={this.handleChange} value={ this.state.rule } 
				settings={ this.state.settings }/>
			<hr />
			<Button bsStyle='danger' onClick={this.handleClick} >Delete Rule</Button>
		</Panel>
		

	}
})

module.exports = RuleConfig
