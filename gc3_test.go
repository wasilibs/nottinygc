package nottinygc_test

import (
	"strings"
	"testing"
)

// Some arbitrary stress tests. It's unclear what is a correct test for GC in practice.

func TestStress(t *testing.T) {
	for i := 0; i < 10000; i++ {
		a := strings.Repeat("a", 100000)
		b := strings.Repeat("b", 100000)
		c := strings.Repeat("c", 100000)

		if strings.Count(a, "a") != 100000 {
			t.Fatal("corrupted heap")
		}
		if strings.Count(b, "b") != 100000 {
			t.Fatal("corrupted heap")
		}
		if strings.Count(c, "c") != 100000 {
			t.Fatal("corrupted heap")
		}
	}
}
