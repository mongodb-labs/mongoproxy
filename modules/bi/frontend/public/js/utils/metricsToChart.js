// Should try to develop here...
var _ = require('lodash');
var moment = require('moment');

var Controller = require('../ajax/controller')

var timeFormat = "YYYY-MM-DD HH:mm:ss";

var g = require('./granularities')

function getMetrics(rules, granularity, index, callback) {
	var endTime = moment();
	var startTime = moment();
	var rule = rules[index];

	// TODO: eventually be able to expand the range
	var range = 60;
	startTime.subtract(range, g.getProperGranularity(granularity));

	Controller.getMetric(index, granularity, startTime, endTime,
		function(data) {
			callback(data);
		});

}

module.exports = getMetrics;
