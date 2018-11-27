package compose

import (
	"os"
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

func TestArgsToEnv(t *testing.T) {
	argsToEnv([]string{"a", "b", "c", "d"})

	t.Run("test a", testSumFunc("test ${arg.0}", "test a"))
	t.Run("test ab", testSumFunc("test ${arg.0}${arg.1}", "test ab"))
	t.Run("test abcd", testSumFunc("test ${arg.0}${arg.1}${arg.2}${arg.3}", "test abcd"))
	t.Run("test abcd", testSumFunc("test ${arg.0}${arg.1}${arg.2}${arg.3}${arg.4}", "test abcd"))
}

func testSumFunc(testValue string, expected string) func(*testing.T) {
	return func(t *testing.T) {
		result := os.ExpandEnv(testValue)
		if result != expected {
			t.Errorf("Expected %s got %s", expected, result)
		}
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
