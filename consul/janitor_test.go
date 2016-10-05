package consul

import "testing"
import "reflect"

type testpair struct {
	input    []string
	expected []string
}

var tests = []testpair{
	{
		[]string{
			"containers/alice/cntrInfodUrl",
			"containers/bob/hostinfo/hostname",
		},
		[]string{
			"alice",
			"bob",
		},
	},
	{
		[]string{
			"containers/ahluzhqe/cntrInfodUrl",
		},
		[]string{
			"ahluzhqe",
		},
	},
	{
		[]string{},
		[]string{},
	},
	{
		[]string{
			"containers",
			"containers/bob/hostinfo/hostname",
		},
		[]string{
			"bob",
		},
	},
}

func TestDecodeInstanceIds(t *testing.T) {
	for _, pair := range tests {

		actual := DecodeInstanceIds(pair.input)

		if !reflect.DeepEqual(actual, pair.expected) {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", actual,
			)
		}

	}
}
