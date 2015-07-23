'use strict';

var _ = require('lodash');
var moment = require('moment');

var Controller = require('../ajax/controller');
var g = require('./granularities');

// gets the data in tabular form (that can be consumed by a TimeseriesChart) for the rule at
// rules[index] and the granularity, and on success calls the callback, and on failure calls the 
// error functions.
// Currently, gets data in a range spanning from 60 units of time before the present time, up
// to the present time.
function getCurrentMetric(rule, granularity, callback, error) {
	var endTime = moment();
	var startTime = moment();

	// helper function to fill out the rest of the array with the proper times
	// so that the graph draws properly.
	function processMetric(data, callback) {
		var graphData = [];
		var graphTime = [];
		var roundedTime = startTime.clone().startOf(g.getProperGranularity(granularity));

		var i = 0;

		while (!roundedTime.isAfter(endTime) && i < data.data.length) {
			var dataObj = data.data[i];
			
			// The front end is a bit unstable, because server's days-in-month conversion is inexact
			// We round the times we get from the server so that the chart will display properly.
			var time = moment(dataObj.time);

			time = g.roundToGranularity(time, granularity);

			if (time.diff(roundedTime) == 0) {
				graphData.push(dataObj.value);
				i++;
			}
			else {
				graphData.push(0);
			}
			graphTime.push(roundedTime.format("YYYY-MM-DD HH:mm:ss"))
			roundedTime.add(1, g.getProperGranularity(granularity))
		}

		while (graphData.length < range) {
			graphTime.push(roundedTime.format("YYYY-MM-DD HH:mm:ss"));
			graphData.push(0);
			roundedTime.add(1, g.getProperGranularity(granularity));
		}
		callback({
			data: graphData,
			time: graphTime
		});
	}

	// TODO: eventually be able to expand the range
	var range = 60;
	startTime.subtract(range, g.getProperGranularity(granularity));

	if (rule.ValueType) {
		Controller.getTabularMetricValue(rule.index, granularity, 
			startTime, endTime, rule.ValueType,
			function(data) {
				processMetric(data, callback);
			}, error);
	}
	else {
		Controller.getTabularMetric(rule.index, granularity, 
			startTime, endTime, function(data) {
			processMetric(data, callback);
		}, error);

	}

}

module.exports = getCurrentMetric;
