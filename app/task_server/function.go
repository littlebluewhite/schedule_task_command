package task_server

func commandReturn2Variables(variables map[int]map[string]string,
	comB comBuilder) map[int]map[string]string {
	for _, parserItem := range comB.parser {
		for _, to := range parserItem.To {
			variables[to.ID][to.Key] = comB.com.Return[parserItem.FromKey]
		}
	}
	return variables
}
