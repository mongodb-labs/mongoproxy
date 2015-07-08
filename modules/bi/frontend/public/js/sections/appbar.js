var React = require('react')

var Navbar = require('react-bootstrap').Navbar
var Nav = require('react-bootstrap').Nav
var NavItem = require('react-bootstrap').NavItem
var DropdownButton = require('react-bootstrap').DropdownButton
var MenuItem = require('react-bootstrap').MenuItem

var Appbar = React.createClass({
	render: function() {
		return <Navbar brand={<a href="#">React-Bootstrap</a>}>
		    <Nav>
		      <NavItem eventKey={1} href='#'>Link</NavItem>
		      <NavItem eventKey={2} href='#'>Link</NavItem>
		      <DropdownButton eventKey={3} title='Dropdown'>
		        <MenuItem eventKey='1'>Action</MenuItem>
		        <MenuItem eventKey='2'>Another action</MenuItem>
		        <MenuItem eventKey='3'>Something else here</MenuItem>
		        <MenuItem divider />
		        <MenuItem eventKey='4'>Separated link</MenuItem>
		      </DropdownButton>
		    </Nav>
		  </Navbar>;
	}
})

module.exports = Appbar;
