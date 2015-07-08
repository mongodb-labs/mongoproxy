window.jQuery = window.$ = require('jquery')
require('bootstrap')

var React = require('react')
var Appbar = require('./sections/appbar')

var Rules = require('./sections/rules')


var App = React.createClass({
	render: function() {
		return <div>
			<Appbar />
			<div className="container">
				<h1>Hello World!</h1>
				<hr />
				<Rules />
			</div>
		</div>;
	}
})

React.render(<App />, document.getElementById("app"));
