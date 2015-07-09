function addLabelsToRules(rules) {
	if (!rules || !rules.length) {
		return rules;
	}
	for (var i = 0; i < rules.length; i++) {
		rules[i].value = rules[i].ValueField;
	}

	return rules;
}

module.exports = addLabelsToRules;
