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
                json: (data || [])
            }
        });
    },

    componentDidMount: function() {
        this._renderChart(this.props.data);
    },

    componentWillReceiveProps: function(newProps) {
        console.log("Received props")
        this.chart.load({
            json: newProps.data
        }); // or whatever API you need
    },

    render: function() {
        return (
            <div className="row" id={this.props.panelID}></div>
        )
    }
});

module.exports = TimeseriesChart
