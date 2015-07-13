'use strict';

var _ = require('lodash');
var moment = require('moment');

var Controller = require('../ajax/controller');
var g = require('./granularities');

function getCurrentMetric(rules, granularity, index, callback, error) {
	var endTime = moment();
	var startTime = moment();
	var rule = rules[index];

	// TODO: eventually be able to expand the range
	var range = 60;
	startTime.subtract(range, g.getProperGranularity(granularity));

	Controller.getTabularMetric(index, granularity, startTime, endTime,
		function(data) {
			callback(data);
		}, error);

}

module.exports = getCurrentMetric;
