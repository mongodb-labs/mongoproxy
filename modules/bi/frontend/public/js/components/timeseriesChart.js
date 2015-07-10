var React = require('react')
var c3 = require("c3");

var TimeseriesChart = React.createClass({
    _renderChart: function(data) {
        // save reference to our chart to the instance
        var self = this;
        console.log(self.props)
        this.chart = c3.generate({
            bindto: "#" + self.props.panelID,
            data: {
                x: "time",
                xFormat: '%Y-%m-%dT%H:%M:%SZ',
                columns: (data || [])
            },
            type: "spline",
            axis: {
                x: {
                    type: 'timeseries',
                    tick: {
                        //              format : "%m/%d" // https://github.com/mbostock/d3/wiki/Time-Formatting#wiki-format
                        format: "%Y-%m-%d %H:%M:%S", // https://github.com/mbostock/d3/wiki/Time-Formatting#wiki-format
                        count: 5,
                    },

                }
            }
        });
    },

    componentDidMount: function() {
        this._renderChart(this.props.data);
    },

    componentWillReceiveProps: function(newProps) {
        this.chart.load({
            columns: newProps.data
        }); // or whatever API you need
    },

    render: function() {
        return (
            <div className="row" id={this.props.panelID}></div>
        )
    }
});

module.exports = TimeseriesChart
