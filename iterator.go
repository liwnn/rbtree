package rbtree

type Iterator struct {
	t *RBTree
	x *node
}

func (it *Iterator) Valid() bool {
	return it.x != it.t.nil
}

func (it *Iterator) Next() {
	n := it.t.successor(it.x)
	it.x = n
}

func (it *Iterator) Value() Item {
	return it.x.item
}
