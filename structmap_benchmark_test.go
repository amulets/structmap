package structmap_test

import (
	"testing"

	"github.com/amulets/structmap"
	"github.com/amulets/structmap/behavior/cast"
	"github.com/amulets/structmap/behavior/name"
)

type Benchmark struct {
	String string
	Number int
	Slice  []string
	Map    map[string]string
}

func BenchmarkDecode(b *testing.B) {
	sm := structmap.New(structmap.WithBehaviors(name.Noop))

	from := map[string]interface{}{
		"String": "MyString",
		"Number": 1000,
		"Slice":  []string{"1", "2", "3"},
		"Map": map[string]string{
			"A": "B",
		},
	}

	var to Benchmark
	for i := 0; i < b.N; i++ {
		if err := sm.Decode(from, &to); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkDecodeWithCast(b *testing.B) {
	sm := structmap.New(structmap.WithBehaviors(name.Noop, cast.ToType()))

	from := map[string]interface{}{
		"String": 2000,
		"Number": "1000",
		"Slice":  []int{1, 2, 3},
		"Map": map[interface{}]interface{}{
			"A": "B",
		},
	}

	var to Benchmark
	for i := 0; i < b.N; i++ {
		if err := sm.Decode(from, &to); err != nil {
			b.Error(err)
		}
	}
}
