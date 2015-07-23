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

// ConnectionConfig is an editor for the connection part of BI Module configuration.
var ConnectionConfig = React.createClass({
	getInitialState: function() {
		return {
			connection: _.extend(defaultConnection, window.config.connection)
		}
	},

	handleChange: function(value) {
		// bubble up to parent
		this.setState({
			connection: value
		});
		this.props.onChange(this, value)
	},

	render: function() {
		var settings = {
			form: false,
			editing: 'always',
		};
		return <Panel>
			<h2>Connection</h2>
			<hr />
			<JSONEditor onChange={this.handleChange} value={ this.state.connection } 
				settings={ settings }/>
		</Panel>
	}
})

module.exports = ConnectionConfig
