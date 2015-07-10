var React = require('react');
var Panel = require('react-bootstrap').Panel;
var async = require('async');

var GranularityToggle = require('../components/granularityToggle')
var RuleSelector = require('../components/ruleSelector')
var TimeseriesChart = require('../components/timeseriesChart')

var getMetrics = require('../utils/metricsToChart')

var GraphPanel = React.createClass({

	getInitialState: function() {
		return {
			timeGranularities: ["M", "D", "h", "m", "s"],
			data: [],
			rules: []
		};
	},
	componentDidMount: function() {
		var self = this;

		setInterval(function() {
			var graphData = [];
			async.forEachOf(self.state.rules, function(rule, key, callback) {

				getMetrics(self.state.rules, "s", key, function(data) {
					if (!graphData.length) {
						graphData[0] = data.time;
						graphData[0].unshift('time');
					}
					gData = data.data;

					// label for the rule
					gData.unshift(self.state.rules[key].ValueField)

					graphData.push(gData)
					callback();
				});

			}, function(err) {
				self.setState({
					data: graphData
				});
			});
		}, 1000);

	},
	handleRuleChange: function(ruleSelector) {
		var newRules = [];
		var s = ruleSelector.state.selected;

		for (var i = 0; i < this.props.rules.length; i++) {
			if (s[i]) {
				newRules.push(this.props.rules[i]);
			}
		}

		this.setState({
			rules: newRules
		})
	},
	render: function() {
		return (
			<Panel>
			<RuleSelector onChange={this.handleRuleChange} rules={this.props.rules}/>
			<GranularityToggle panelID={this.props.panelID} granularities={this.state.timeGranularities}/>
			<TimeseriesChart data={this.state.data} panelID={this.props.panelID}/>
		</Panel>
		)
	}
})

module.exports = GraphPanel
