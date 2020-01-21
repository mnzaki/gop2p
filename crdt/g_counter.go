package crdt

import "fmt"

/**
 * Payload
 */
type GCounter struct {
	Id       ID
	Counters map[ID]int
}
func MakeGCounter(id ID) GCounter {
	var g GCounter
	g.Id = id
	g.Counters = make(map[ID]int)
	return g
}
func (g GCounter) String() string {
	return fmt.Sprintf("GCounter{Value:%v\tID:%v\t%v}", g.Value(), g.Id, g.Counters)
}

/**
 * Update
 */
func (g *GCounter) Increment() {
	g.Counters[g.Id]++
}

/**
 * Query
 */
func (g *GCounter) Value() int {
	sum := 0
	for _, counter := range g.Counters {
		sum += counter
	}
	return sum
}

/**
 * Compare
 */
func (g *GCounter) Compare(f *GCounter) bool {
	for id, counter := range g.Counters {
		other, ok := f.Counters[id]
		if !ok || counter > other {
			return false
		}
	}
	return true
}

/**
 * Merge
 * just taking the maximum of each counter
 * which is a also a least upper bound for this datatype
 */
func (g *GCounter) Merge(f *GCounter) GCounter {
	newG := MakeGCounter(g.Id)
	for id, counter := range g.Counters {
		newG.Counters[id] = counter
	}

	for id, counter := range f.Counters {
		my, ok := g.Counters[id]
		if !ok || counter > my {
			newG.Counters[id] = counter
		} else {
			newG.Counters[id] = my
		}
	}
	return newG
}
