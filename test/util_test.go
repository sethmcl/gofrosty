package test

import (
	"github.com/sethmcl/gofrosty/lib"
	"github.com/sethmcl/gofrosty/lib/testutil"
	"testing"
)

func TestTruncatePath(t *testing.T) {
	test := testutil.New(t)

	type TableEntry struct {
		Path     string
		Count    int
		Expected string
	}
	table := []TableEntry{
		TableEntry{"/foo/bar/biz", 0, "/foo/bar/biz"},
		TableEntry{"/foo/bar/biz", 1, "/bar/biz"},
		TableEntry{"foo/bar/biz", 1, "bar/biz"},
		TableEntry{"foo/bar/biz", 0, "foo/bar/biz"},
		TableEntry{"/foo/bar/biz", -1, "/foo/bar"},
		TableEntry{"/foo/bar/biz", -2, "/foo"},
		TableEntry{"/foo/bar/biz", -3, "/"},
		TableEntry{"/foo/bar/biz", 3, ""},
		TableEntry{"/foo/bar/biz", 2, "/biz"},
		TableEntry{"", 0, ""},
	}

	for _, entry := range table {
		actual, err := lib.TruncatePath(entry.Path, entry.Count)
		if entry.Expected == "error" {
			test.Assert(err != nil, true)
		} else {
			test.Assert(err, nil)
			test.Assert(actual, entry.Expected, entry.Path, entry.Count)
		}
	}
}

func TestSliceContains(t *testing.T) {
	test := testutil.New(t)
	strSlice := []string{"foo", "bar", "biz"}
	test.Assert(lib.SliceContains(strSlice, "foo"), true)
	test.Assert(lib.SliceContains(strSlice, "zzz"), false)
}
