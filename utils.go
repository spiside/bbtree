package bbtree

type Int int

func (li Int) More(than Comparable) bool {
	return li > than.(Int)
}
