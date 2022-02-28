package v1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCompany_Validate(t *testing.T) {
	type TestCase struct {
		data     string
		expected string
		actual   string
	}
	path := []string{"KloverCloud CI-CD\\ctl", "KloverCloud CI-CD\\ctl\\"}
	exp := []string{"KloverCloud CI-CD\\ctl\\", "KloverCloud CI-CD\\ctl\\"}

	var testdata []TestCase

	for i := 0; i < len(testdata); i++ {
		testcase := TestCase{
			data: path[i],
			expected: exp[i],
		}
		testdata = append(testdata, testcase)
	}

	for i := 0; i < len(testdata); i++ {
		testdata[i].actual = GetCfgPath(path[i])
		if !reflect.DeepEqual(testdata[i].expected, testdata[i].actual) {
			fmt.Println(testdata[i].actual, i)
			assert.ElementsMatch(t, testdata[i].expected, testdata[i].actual)
		}
	}
}

