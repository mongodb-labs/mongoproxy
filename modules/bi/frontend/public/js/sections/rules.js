var React = require('react')

var TabbedArea = require('react-bootstrap').TabbedArea
var TabPane = require('react-bootstrap').TabPane

var RulePanel = require('./rulePanel')
var tabs = [];
for (var i = 0; i < window.config.Rules.length; i++) {
	tabs.push(
		<TabPane key={i} eventKey={i} tab="Rule">
					<RulePanel rule = {window.config.Rules[i]} index={i} />
				</TabPane>
	)
}


var Rules = React.createClass({
	render: function() {
		return (
			<TabbedArea defaultActiveKey={0}>
				{tabs}
			</TabbedArea>
		)
	}
})

module.exports = Rules;
