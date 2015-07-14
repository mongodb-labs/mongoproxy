'use strict';

var React = require('react');
var Panel = require('react-bootstrap').Panel;
var JSONEditor = require('react-json');

var Button = require('react-bootstrap').Button;

var _ = require('lodash');

var tg = require('../utils/convertTimeGranularity');

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
var RuleConfig = React.createClass({
	getInitialState: function() {
		return {
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
	handleClick: function() {

		this.props.onClick(this.props.key);
		// bubble up to parent
	},
	render: function() {
		return <Panel>
			<JSONEditor value={ this.state.rule } settings={ this.state.settings }/>
			<hr />
			<Button onClick={this.handleClick} >Delete Rule</Button>
		</Panel>
		

	}
})

module.exports = RuleConfig
