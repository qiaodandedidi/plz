package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/dump"
)

func Test_map(t *testing.T) {
	t.Run("map int to int", test.Case(func(ctx *countlog.Context) {
		must.JsonEqual(`{
		"__root__": {
			"type": "map[int]int",
			"data": {
				"__ptr__": "{ptr1}"
			}
		},
		"{ptr1}": {
			"count": 2,
			"flags": 0,
			"B": 0,
			"noverflow": 0,
			"hash0": "{ANYTHING}",
			"buckets": {"__ptr__":"{ptr2}"},
			"oldbuckets": {"__ptr__":"0"},
			"nevacuate": 0,
			"extra": {"__ptr__":"0"}
		},
		"{ptr2}": [{
			"tophash": "{ANYTHING}",
			"keys": [9,8,0,0,0,0,0,0],
			"elems": [7,6,0,0,0,0,0,0]
		}]}`, dump.Var{map[int]int{
			9: 7,
			8: 6,
		}}.String())
	}))
	t.Run("map string to string", test.Case(func(ctx *countlog.Context) {
		must.JsonEqual(`{
		"__root__": {
			"type": "map[string]string",
			"data": {
				"__ptr__": "{ptr1}"
			}
		},
		"{ptr1}": {
			"count": 2,
			"flags": 0,
			"B": 0,
			"noverflow": 0,
			"hash0": "{ANYTHING}",
			"buckets": {"__ptr__":"{ptr2}"},
			"oldbuckets": {"__ptr__":"0"},
			"nevacuate": 0,
			"extra": {"__ptr__":"0"}
		},
		"{ptr2}": [{
			"tophash": "{ANYTHING}",
			"keys": [
				{"data":{"__ptr__":"{key1}"},"len":1},
				{"data":{"__ptr__":"{key2}"},"len":1},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0}
			],
			"elems": [
				{"data":{"__ptr__":"{elem1}"},"len":1},
				{"data":{"__ptr__":"{elem2}"},"len":1},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0},
				{"data":{"__ptr__":"0"},"len":0}
			]
		}],
		"{key1}":"a",
		"{key2}":"c",
		"{elem1}":"b",
		"{elem2}":"d"
		}`, dump.Var{map[string]string{
			"a": "b",
			"c": "d",
		}}.String())
	}))
}