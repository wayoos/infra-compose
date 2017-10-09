package compose

import (
	"sort"
	"testing"
)

func TestAppendEnv(t *testing.T) {

	envRoot := Environment{"A=0", "ROOT_TEST=test"}
	env := Environment{"A=0", "ROOT_TEST=new", "ROOT_TEST2=new", "NEW=ENV"}

	mergedEnv := appendEnv(envRoot, env)

	expected := Environment{"A=0", "ROOT_TEST2=new", "ROOT_TEST=test", "NEW=ENV"}

	if !isEqual(expected, mergedEnv) {
		t.Error("Expected, got ", mergedEnv)
	}
}

func isEqual(a1 []string, a2 []string) bool {
	sort.Strings(a1)
	sort.Strings(a2)
	if len(a1) == len(a2) {
		for i, v := range a1 {
			if v != a2[i] {
				return false
			}
		}
	} else {
		return false
	}
	return true
}
