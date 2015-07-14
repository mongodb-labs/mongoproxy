'use strict';

function convertToStringArray(input){
	var oldGranularities = input.timeGranularity;
	if (!oldGranularities) {
		return input;
	}
	var newGranularities = [];
	if (oldGranularities.Month) {
		newGranularities.push("M");
	}
	if (oldGranularities.Day) {
		newGranularities.push("D");
	}
	if (oldGranularities.Hour) {
		newGranularities.push("h");
	}
	if (oldGranularities.Minute) {
		newGranularities.push("m");
	}
	if (oldGranularities.Second) {
		newGranularities.push("s");
	}

	input.timeGranularity = newGranularities;
	return input;
}

function convertToBooleans(input) {
	var newGranularities = {
		Month: false,
		Day: false,
		Hour: false,
		Minute: false,
		Second: false
	};
	var oldGranularities = input.timeGranularity;
	if (!oldGranularities) {
		return input;
	}

	for (var i = 0; i < oldGranularities.length; i++) {
		switch (oldGranularities[i]) {
			case "M":
				newGranularities.Month = true;
			break;
			case "D":
				newGranularities.Day = true;
			break;
			case "h":
				newGranularities.Hour = true;
			break;
			case "m":
				newGranularities.Minute = true;
			break;
			case "s":
				newGranularities.Second = true;
			break;
		}
	}

	input.timeGranularity = newGranularities;
	return input;

}

module.exports = {
	convertToBooleans: convertToBooleans,
	convertToStringArray: convertToStringArray
};
