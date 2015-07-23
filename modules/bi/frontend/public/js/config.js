'use strict';

// Application file for the configuration editor.
// TODO: Eventually split this entirely from the main, as it doesn't share
// any components with the dashboard view.

window.jQuery = window.$ = require('jquery');
require('bootstrap/js/button');
require('./vendor/jquery.timer');
require('sweetalert');

var _ = require('lodash');

var React = require('react');

var ConnectionConfig = require('./sections/connectionConfig');
var RuleConfig = require('./sections/ruleConfig');
var AddRuleButton = require('./components/addRule');

var Panel = require('react-bootstrap').Panel;
var Button = require('react-bootstrap').Button;

var Controller = require('./ajax/controller');

var tg = require('./utils/convertTimeGranularity');

// initialize the application
var App = React.createClass({
	getInitialState: function() {
		return {
			// HTML elements for individual rules, to be rendered
			rulePanels: [],

			// counter to ensure unique keys for rule panels
			key: 0,

			// the current configuration displayed in the application
			configuration: window.config,
		}
	},
	componentDidMount: function() {
		var r = [];
		for (var i = 0; i < window.config.rules.length; i++) {
			
			r.push(
				<RuleConfig onChange={this.handleRuleChange} onClick={this.handleRemoveRule} 
					key={i} keyIndex={i} rule={window.config.rules[i]} />
			)
		}	
		this.setState({
			rulePanels: r,
			key: r.length,
		})

	},
	handleRemoveRule: function(i) {
		var r = this.state.rulePanels;
		var newConfig = _.extend({}, this.state.configuration);
		var index = _.findIndex(r, function(rule) {
  			return rule.key === i;
		});

		// remove the rule both from the HTML elements and from the state
		r.splice(index, 1);
		newConfig.rules.splice(index, 1);

		this.setState({
			rulePanels: r,
			configuration: newConfig
		})

	},
	addRule: function() {
		var r = this.state.rulePanels;
		r.push(
			<RuleConfig onChange={this.handleRuleChange} onClick={this.handleRemoveRule} 
				key={this.state.key} rule={{}} keyIndex={this.state.key}/>
		)
		this.setState({
			rulePanels: r,
			key: this.state.key + 1
		})
	},

	// update a particular rule whenever the rule's information changes
	handleRuleChange: function(component) {

		// react doesn't let us see the state of all components when only
		// one is updated, so we have to find the correct rule in the current
		// state and modify that one
		var newConfig = _.extend({}, this.state.configuration);
		var index = _.findIndex(this.state.rulePanels, function(rule) {
  			return rule.key.toString() === component.props.keyIndex.toString();
		});

		newConfig.rules[index] = component.state.rule;

		this.setState({
			configuration: _.extend(this.state.configuration, newConfig)
		});

	},

	// update the state whenever the connection information changes
	handleConnectionChange: function(component) {
		var newConfig = {
			connection: component.state.connection
		}
		this.setState({
			configuration: _.extend(this.state.configuration, newConfig)
		});

	},

	handleSaveConfig: function() {
		// save the configuration
		var savedConfig = _.extend({}, this.state.configuration);
		for (var i = 0; i < savedConfig.rules.length; i++) {
			savedConfig.rules[i] = tg.convertToStringArray(savedConfig.rules[i]);
		}
		Controller.postConfiguration(savedConfig, function(data) {
			sweetAlert("Success", "Configuration successfully saved.", "success");
		}, function(error) {
			sweetAlert("Error", "Configuration was not able to be saved.", "error");
		});

	},

	render: function() {
		return <div>
			<div className="container">
				<Button onClick={this.handleSaveConfig} block bsStyle="primary" 
					bsSize="large">Save Configuration</Button>
				<ConnectionConfig onChange={this.handleConnectionChange} />
				<Panel>
					<h2>Rules</h2>
					<hr />
					{this.state.rulePanels}

					<AddRuleButton onClick={this.addRule} />
				</Panel>
			</div>
		</div>;
	}
})

React.render(<App />, document.getElementById("app"));
