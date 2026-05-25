package internal

import (
	"os"
	"strings"
)

func ResolveENV(from ...map[string]string) func(string) string {
	return func(s string) string {
		for _, envs := range from {
			if value, ok := envs[s]; ok {
				return value
			}
		}
		return os.Getenv(s)
	}
}

func MapToEnvSlice(m map[string]string) []string {
	var env []string
	for k, v := range m {
		env = append(env, k+"="+v)
	}
	return env
}

func EnvSliceToMap(s []string) map[string]string {
	var result = map[string]string{}
	for _, m := range s {
		KV := strings.Split(m, "=")
		if len(KV) != 2 {
			continue
		}
		result[KV[0]] = KV[1]
	}
	return result
}
