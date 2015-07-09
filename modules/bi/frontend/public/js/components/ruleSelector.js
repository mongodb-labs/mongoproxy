var React = require('react')
var Multiselect = require('react-bootstrap-multiselect');

var RuleSelector = React.createClass({
	render: function() {
		return <Multiselect data={this.props.rules} multiple/>
	}
})

module.exports = RuleSelector
