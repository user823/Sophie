package log

import "go.uber.org/zap"

type environment map[string]zap.Field

// 后加入的key 会替代原先的key
func (e environment) addValues(fields ...zap.Field) {
	for _, f := range fields {
		e[f.Key] = f
	}
}

func newEnvironment() environment {
	return make(map[string]zap.Field)
}
