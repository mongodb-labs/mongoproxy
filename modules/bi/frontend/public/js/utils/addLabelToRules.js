'use strict';

// helper function that is called on application start, to add some useful
// fields to the rules
function addLabelsToRules(rules) {
	if (!rules || !rules.length) {
		return rules;
	}

	for (var i = 0; i < rules.length; i++) {
		rules[i].label = rules[i].ValueField;
		rules[i].value = i;
		rules[i].index = i;
	}

	return rules;
}

module.exports = addLabelsToRules;
