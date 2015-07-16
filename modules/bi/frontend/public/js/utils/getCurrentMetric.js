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
			var graphData = [];
			var graphTime = [];
			var roundedTime = startTime.clone().startOf(g.getProperGranularity(granularity));

			var i = 0;

			while(roundedTime.isBefore(endTime) && i < data.data.length) {
				var dataObj = data.data[i];
				var time = moment(dataObj.time);
				if (time.diff(roundedTime) == 0) {
					graphData.push(dataObj.value);
					i ++;
				} else {
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
		}, error);

}

module.exports = getCurrentMetric;
