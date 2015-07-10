// helper function to get the proper time scale for moment.JS given a time granularity.
var getProperGranularity = function(granularity) {
	switch (granularity) {
		case "M":
			return "month"
			break;
		case "D":
			return "day"
			break;
		case "h":
			return "hour"
			break;
		case "m":
			return "minute"
			break;
		case "s":
			return "second"
			break;
		default:
			return ""
	}
}

// helper function to get the next higher time scale for moment.JS given a time granularity.
var getHigherGranularity = function(granularity) {
	switch (granularity) {
		case "M":
			return "year"
			break;
		case "D":
			return "month"
			break;
		case "h":
			return "day"
			break;
		case "m":
			return "hour"
			break;
		case "s":
			return "minute"
			break;
		default:
			return ""
	}
}

// helper function to get the number of ticks of this granularity happens in the next higher time
// granularity.
var getNumInNextGranularity = function(granularity, time) {
	switch (granularity) {
		case "M":
			return 12
			break;
		case "D":
			switch (time.month) {
				case 4:
				case 6:
				case 9:
				case 11:
					return 30;
					break;
				case 2:
					if (time.isLeapYear()) {
						return 29;
					}
					return 28;
					break;
				default:
					return 31;
			}
			break;
		case "h":
			return 24
			break;
		case "m":
			return 60
			break;
		case "s":
			return 60
			break;
		default:
			return 0
	}
}

module.exports = {
	getNumInNextGranularity: getNumInNextGranularity,
	getHigherGranularity: getHigherGranularity,
	getProperGranularity: getProperGranularity
};
