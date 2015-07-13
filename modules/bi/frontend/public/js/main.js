window.jQuery = window.$ = require('jquery');
require('bootstrap/js/button');
require('./vendor/jquery.timer');

var React = require('react');
var UniqeIdMixin = require('unique-id-mixin');

var GraphPanel = require('./sections/graphPanel');

// add extra fields to the rules
var addLabelToRules = require('./utils/addLabelToRules');
addLabelToRules(window.config.Rules);

// initialize the application
var App = React.createClass({
	mixins: [ UniqeIdMixin ],
	render: function() {
		return <div>
			<div className="container">
				<h1>Hello World!</h1>
				<hr />
				<GraphPanel panelID={this.getNextUid('panel')} rules={window.config.Rules}/>
			</div>
		</div>;
	}
})

React.render(<App />, document.getElementById("app"));
