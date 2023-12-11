package log

type environment map[string]string

// 后加入的key 会替代原先的key
func (e environment) addValues(fields ...string) {
	for i := 0; i < len(fields); i += 2 {
		e[fields[i]] = fields[i+1]
	}
}

func newEnvironment() environment {
	return make(map[string]string)
}
