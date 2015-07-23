'use strict';

var React = require('react')
var c3 = require("c3");

// TimeseriesChart is a line chart with time on the x axis, and whose data is received
// in the props.data field.
var TimeseriesChart = React.createClass({
    _renderChart: function(data) {
        // save reference to our chart to the instance
        var self = this;
        this.chart = c3.generate({
            data: {
                x: "time",
                xFormat: '%Y-%m-%d %H:%M:%S',
                columns: (data || [])
            },
            type: "spline",
            point: {
                show: false
            },
            transition: {
                duration: 100
            },
            axis: {
                x: {
                    type: 'timeseries',
                    tick: {
                        // https://github.com/mbostock/d3/wiki/Time-Formatting#wiki-format
                        format: "%Y-%m-%d %H:%M:%S", 
                        count: 5,
                    },

                }
            }
        });

        $("#" + self.props.panelID).append(this.chart.element);
    },

    componentDidMount: function() {
        this._renderChart(this.props.data);
    },

    componentWillReceiveProps: function(newProps) {
        if (newProps.unload == true) {
            this.chart.load({
                columns: newProps.data,
                unload: true,
            });
        }
        else {
            this.chart.load({
                columns: newProps.data
            });
        }
        

    },

    unloadChart: function() {
        this.chart.unload();
    },

    render: function() {
        return (
            <div className="row" id={this.props.panelID}></div>
        )
    }
});

module.exports = TimeseriesChart
