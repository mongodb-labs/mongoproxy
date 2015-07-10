var React = require('react');
var Panel = require('react-bootstrap').Panel;
var async = require('async');

var GranularityToggle = require('../components/granularityToggle')
var RuleSelector = require('../components/ruleSelector')
var TimeseriesChart = require('../components/timeseriesChart')

var getMetrics = require('../utils/metricsToChart')

var clearChartTimeout;
var chartInterval;

function setChart(graphPanel, unload) {
	var graphData = [];
	async.forEachOf(graphPanel.state.rules, function(rule, key, callback) {

		i = rule.index || key;
		getMetrics(graphPanel.props.rules, graphPanel.state.granularity, i, function(data) {

			if (!graphData.length) {
				graphData[0] = data.time;
				graphData[0].unshift('time');
			}
			gData = data.data;

			// label for the rule
			gData.unshift(graphPanel.props.rules[key].ValueField)

			graphData.push(gData)
			callback();
		});

	}, function(err) {
		graphPanel.setState({
			data: graphData,
			unload: unload,
		});
	});
}

var GraphPanel = React.createClass({

	getInitialState: function() {
		return {
			granularity: "m",
			data: [],
			rules: [],
			unload: false,
		};
	},
	componentDidMount: function() {
		var self = this;

		chartInterval = $.timer(function() {
			setChart(self, false);
		}, 1000, true);

	},
	handleGranularityToggle: function(timeToggle) {
		chartInterval.pause();
		clearTimeout(clearChartTimeout);

		this.setState({
			granularity: timeToggle.state.selected
		});

		this.refs.chart.unloadChart();
		var self = this;

		// needs to be after c3's animation cycle
		clearChartTimeout = setTimeout(function() {
			setChart(self, false)
			chartInterval.play();
		}, 500);
	},
	handleRuleChange: function(ruleSelector) {
		chartInterval.pause();
		clearTimeout(clearChartTimeout);
		var newRules = [];
		var s = ruleSelector.state.selected;

		for (var i = 0; i < this.props.rules.length; i++) {
			if (s[i]) {
				var rule = this.props.rules[i];
				rule.index = i;
				newRules.push(rule);
			}
		}
		this.setState({
			rules: newRules,
		})
		this.refs.chart.unloadChart();
		var self = this;

		// needs to be after c3's animation cycle
		clearChartTimeout = setTimeout(function() {
			setChart(self, false)
			chartInterval.play();
		}, 500);
	},
	render: function() {
		return (
			<Panel>
			<RuleSelector onChange={this.handleRuleChange} rules={this.props.rules}/>
			<GranularityToggle onChange={this.handleGranularityToggle} panelID={this.props.panelID} ref="timeToggle" granularities={this.state.timeGranularities}/>
			<TimeseriesChart unload={this.state.unload} ref="chart" data={this.state.data} panelID={this.props.panelID}/>
		</Panel>
		)
	}
})

module.exports = GraphPanel
