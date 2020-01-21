package crdt

import (
	"testing"
)

func TestGCounter(t *testing.T) {
	g1 := MakeGCounter(ID(1))
	g2 := MakeGCounter(ID(2))
	for i := 0; i < 6; i++ {
		g1.Increment()
	}
	for i := 0; i < 3; i++ {
		g2.Increment()
	}
	g3 := g1.Merge(&g2)

	if g3.Value() != 9 {
		t.Errorf("Merged value incorrect!\ng1:\t\t%v\ng2:\t\t%v\ng1.Merge(g2):\t%v",
		g1, g2, g3)
	}
}
