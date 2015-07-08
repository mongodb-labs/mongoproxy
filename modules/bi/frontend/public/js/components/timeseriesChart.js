var React = require('react')
var c3 = require("c3");

// this.props.granularities - an array of the time granularities

var TimeseriesChart = React.createClass({
	_renderChart: function (data) {
        // save reference to our chart to the instance
        var self = this;
        console.log(self.props)
        this.chart = c3.generate({
            bindto: "#" + self.props.chartID,
            data: {
              json: data
          	}
        });
    },

    componentDidMount: function () {
        this._renderChart(this.props.data);
    },

    componentWillReceiveProps: function (newProps) {
        this.chart.load({
            json: newProps.data
        }); // or whatever API you need
    },

    render: function () {
        return (
            <div className="row" id={this.props.chartID}></div>
        )
    }
});

module.exports = TimeseriesChart
