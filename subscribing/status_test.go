package subscribing_test

import (
	"testing"

	"github.com/cohen990/exactlyOnce/subscribing"
	"github.com/cohen990/exactlyOnce/testAll"
	"github.com/matryer/is"
)

func TestStatusToString(testing *testing.T) {
	type testInput struct {
		input    subscribing.Status
		expected string
	}

	testCases := map[string]testInput{
		"failed":   {input: subscribing.Failed, expected: "Failed"},
		"received": {input: subscribing.Received, expected: "Received"},
	}

	testAll.Of(testing, testCases, func(testCase testInput) {
		is := is.New(testing)

		is.Equal(testCase.input.String(), testCase.expected)
	})
}
