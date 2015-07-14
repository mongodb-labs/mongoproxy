'use strict';

window.jQuery = window.$ = require('jquery');
require('bootstrap/js/button');
require('./vendor/jquery.timer');

var _ = require('lodash');

var React = require('react');
var UniqeIdMixin = require('unique-id-mixin');

var ConnectionConfig = require('./sections/connectionConfig');
var RuleConfig = require('./sections/ruleConfig');
var AddRuleButton = require('./components/addRule');

// initialize the application
var App = React.createClass({
	getInitialState: function() {
		return {
			rulePanels: [],
			key: 0,
		}
	},
	componentDidMount: function() {
		var r = [];
		for (var i = 0; i < window.config.rules.length; i++) {
			
			r.push(
				<RuleConfig onClick={this.handleRemoveRule} key={i} rule={window.config.rules[i]} />
			)
		}	
		this.setState({
			rulePanels: r,
			key: r.length,
		})

	},
	handleRemoveRule: function(i) {
		// BUG: Can only remove rules in order.
		var r = this.state.rulePanels;
		var index = _.findIndex(r, function(rule) {
  			return rule.key === i;
		});
		r.splice(index, 1);
		this.setState({
			rulePanels: r
		})
	},
	addRule: function() {
		var r = this.state.rulePanels;
		r.push(
			<RuleConfig onClick={this.handleRemoveRule} key={this.state.key} rule={{}} />
		)
		this.setState({
			rulePanels: r,
			key: this.state.key + 1
		})
	},
	render: function() {
		return <div>
			<div className="container">
				<ConnectionConfig />
				{this.state.rulePanels}
				<AddRuleButton onClick={this.addRule} />
			</div>
		</div>;
	}
})

React.render(<App />, document.getElementById("app"));
