var React = require('react');
var Panel = require('react-bootstrap').Panel;

var GranularityToggle = require('../components/granularityToggle')
var RuleSelector = require('../components/ruleSelector')
var TimeseriesChart = require('../components/timeseriesChart')

var getMetrics = require('../utils/metricsToChart')

var GraphPanel = React.createClass({

	getInitialState: function() {
		return {
			timeGranularities: ["M", "D", "h", "m", "s"],
			data: [],
		};
	},
	componentDidMount: function() {
		var self = this;

		setInterval(function() {
			getMetrics(self.props.rules, "s", 0, function(data) {
				var graphData = [];
				console.log(data);
				graphData[0] = data.time;
				graphData[1] = data.data;

				graphData[0].unshift('time');
				graphData[1].unshift('price');

				console.log(graphData);
				self.setState({
					data: graphData
				});
			});
		}, 1000);

	},
	render: function() {
		return (
			<Panel>
			<RuleSelector rules={this.props.rules}/>
			<GranularityToggle panelID={this.props.panelID} granularities={this.state.timeGranularities}/>
			<TimeseriesChart data={this.state.data} panelID={this.props.panelID}/>
		</Panel>
		)
	}
})

module.exports = GraphPanel
