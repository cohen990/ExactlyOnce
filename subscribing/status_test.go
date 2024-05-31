package subscribing_test

import (
	"testing"

	"github.com/cohen990/exactlyOnce/subscribing"
	"github.com/matryer/is"
)

type testInput struct {
	input    subscribing.Status
	expected string
}

func TestStatusToString(t *testing.T) {
	is := is.New(t)

	input := subscribing.Failed
	expected := "Failed"

	is.Equal(input.String(), expected)
}

func TestStatusToString2(t *testing.T) {
	testCases := map[string]struct {
		input    subscribing.Status
		expected string
	}{
		"failed":   {input: subscribing.Failed, expected: "Failed"},
		"received": {input: subscribing.Received, expected: "Received"},
	}

	for desc, testCase := range testCases {
		desc, testCase := desc, testCase
		t.Run(desc, func(t *testing.T) {
			is := is.New(t)

			is.Equal(testCase.input.String(), testCase.expected)
		})
	}
}

type testCases = map[string]testInput

func TestStatusToString3(testing *testing.T) {
	testCases := testCases{
		"failed":   {input: subscribing.Failed, expected: "Failed"},
		"received": {input: subscribing.Received, expected: "Received"},
	}
	test(testing, testCases, func(testCase testInput) {
		is := is.New(testing)

		is.Equal(testCase.input.String(), testCase.expected)
	})
}

func test(t *testing.T, testCases map[string]testInput, test func(testInput)) {
	for desc, testCase := range testCases {
		desc, testCase := desc, testCase
		t.Run(desc, func(t *testing.T) {
			test(testCase)
		})
	}
}
