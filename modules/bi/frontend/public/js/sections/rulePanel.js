var React = require('react');
var moment = require('moment');

var Chart = require('../components/timeseriesChart');
var Controller = require('../ajax/controller');
var GranularityToggle = require('../components/granularityToggle')

var minutesToChart = require('../utils/minutesToChart');


var data = {
	"price": [34, 234, 44, 87, 0, 0, 50]
}

var RulePanel = React.createClass({

	getInitialState: function() {
		var data = {}
		data[this.props.rule.ValueField] = [];
    	return {
    		data: data
    	};
  	},

	componentDidMount: function () {
		var endTime = moment();
		var startTime = moment();
		startTime.year(endTime.year() - 1)
		var self = this;
        Controller.getMetric(this.props.index, "m", this.props.rule.ValueField, startTime, endTime, function(data) {
        	if (!data) {
        		console.log(data);
        		return;
        	}
        	for (var i = 0; i < data.length; i++) {
        		data[i].start = moment(data[i].start)
        	}
        	

			dataArray = minutesToChart(data, data[0].start.clone().add(10, 'minutes'), 100);
			var data = {}
			data[self.props.rule.ValueField] =dataArray;
			self.setState({
				data: data
			})

			console.log(dataArray)
        })
    },

	render: function() {
		return <div>
			{this.props.rule.OriginDatabase}.{this.props.rule.OriginCollection}
			<GranularityToggle rule={this.props.rule} granularities={this.props.rule.TimeGranularities} />
			<hr/>
			<Chart chartID={this.props.rule.ValueField} data={this.state.data}/>
		</div>
	}

});

module.exports = RulePanel;
