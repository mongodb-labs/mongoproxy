'use strict';

var React = require('react');
var Panel = require('react-bootstrap').Panel;
var JSONEditor = require('react-json');

var _ = require('lodash');

var defaultConnection = {
	addresses: [],
	direct: false,
	timeout: 60000,
	auth: {
		username: "username",
		password: "password",
		database: "database"
	}
}
var ConnectionConfig = React.createClass({
	getInitialState: function() {
		return {
			connection: _.extend(defaultConnection, window.config.connection)
		}
	},
	render: function() {
		var settings = {
			form: false,
			editing: 'always',
		};
		return <Panel>
			<JSONEditor value={ this.state.connection } settings={ settings }/>
		</Panel>
	}
})

module.exports = ConnectionConfig
