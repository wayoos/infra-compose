package compose

import "strings"

func appendEnv(rootEnv Environment, env Environment) Environment {
	for _, e := range env {
		if !isPresent(rootEnv, e) {
			rootEnv = append(rootEnv, e)
		}
	}
	return rootEnv
}

func isPresent(env Environment, envVar string) bool {
	variableSplit := strings.Split(envVar, "=")
	variable := envVar
	if len(variableSplit) > 0 {
		variable = variableSplit[0] + "="
	}

	for _, v := range env {
		if strings.HasPrefix(v, variable) {
			return true
		}
	}
	return false
}
