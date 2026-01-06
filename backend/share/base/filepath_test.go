package base

import "testing"

func TestParentDir(t *testing.T) {
	paths := []string{"/path1/f1/a.txt/", "/path2/f2", "/root"}
	expects := []string{"/path1/f1", "/path2", ""}
	for i, path := range paths {
		dir := ParentDir(path)
		if expect := expects[i]; dir != expect {
			t.Fatalf("expect:%s, but get:%s\n", expect, dir)
		}
	}
}
