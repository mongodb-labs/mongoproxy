'use strict';

// Application file for the primary dashboard view.

window.jQuery = window.$ = require('jquery');
require('bootstrap/js/button');
require('./vendor/jquery.timer');

var React = require('react');

var GraphPanel = require('./sections/graphPanel');

// add extra fields to the rules
var addLabelToRules = require('./utils/addLabelToRules');
addLabelToRules(window.config.Rules);

var populateRuleValues = require('./utils/populateRuleValues');

populateRuleValues(function(err) {
	if (err) {
		console.log(err);
	}

	// initialize the application. Currently starts up a single graph panel.
	var App = React.createClass({
		render: function() {
			return <div>
			<div className="container">
				<GraphPanel panelID='panel' rules={window.config.Rules}/>
			</div>
		</div>;
		}
	})

	React.render(<App />, document.getElementById("app"));

})
