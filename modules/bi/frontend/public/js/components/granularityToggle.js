var React = require('react')

var ButtonGroup = require('react-bootstrap').ButtonGroup;
var Button = require('react-bootstrap').Button;

// this.props.granularities - an array of the time granularities

var GranularityToggle = React.createClass({
	getInitialState: function() {
		return {
			buttons: []
		};
	},
	componentDidMount: function() {
		var self = this;
		var b = [];
		for (var i = 0; i < self.props.granularities.length; i++) {
			
			b.push(
				<label key={this.props.panelID+i} className="btn btn-default">
                    <input type="radio" name={self.props.granularities[i]} value={i} /> {self.props.granularities[i]}
                </label> 
			)
		}	
		this.setState({
			buttons: b
		})
	},
	render: function() {
		return <form><ButtonGroup data-toggle="buttons">
			{this.state.buttons}
		</ButtonGroup></form>
	}
});

module.exports = GranularityToggle
