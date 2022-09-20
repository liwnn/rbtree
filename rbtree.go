package rbtree

const (
	DefaultFreeListSize = 32
)

type Item interface {
	Less(than Item) bool
}

type color int8

// enum
const (
	RED   color = 0
	BLACK color = 1
)

type node struct {
	color color
	item  Item
	left  *node
	right *node
	p     *node
}

type FreeList struct {
	freelist []*node
}

func NewFreeList(size int) *FreeList {
	return &FreeList{freelist: make([]*node, 0, size)}
}

func (f *FreeList) newNode() (n *node) {
	index := len(f.freelist) - 1
	if index < 0 {
		return new(node)
	}
	n = f.freelist[index]
	f.freelist[index] = nil
	f.freelist = f.freelist[:index]
	return
}

func (f *FreeList) freeNode(n *node) (out bool) {
	n.item = nil
	n.left = nil
	n.right = nil
	n.p = nil
	if len(f.freelist) < cap(f.freelist) {
		f.freelist = append(f.freelist, n)
		out = true
	}
	return
}

// RBTree is red-black tree
type RBTree struct {
	root     *node
	nil      *node
	freelist *FreeList
	length   int
}

func New() *RBTree {
	t := &RBTree{
		nil: &node{
			color: BLACK,
		},
		freelist: NewFreeList(DefaultFreeListSize),
	}
	t.root = t.nil
	return t
}

/*
   x               y
  / \             / \
 a   y    ->     x	 c
    / \         / \
   b   c       a   b
*/
func (t *RBTree) leftRotate(x *node) {
	y := x.right

	// y的左节点改成x的右节点
	x.right = y.left
	if y.left != t.nil {
		y.left.p = x
	}

	// x 改成y的左节点
	y.left = x
	if x.p == t.nil {
		t.root = y
	} else if x.p.left == x {
		x.p.left = y
	} else {
		x.p.right = y
	}
	y.p = x.p
	x.p = y
}

/*
    y  	       x
   / \        / \
  x	  c  ->  a   y
 / \            / \
a   b          b   c
*/
func (t *RBTree) rightRotate(y *node) {
	x := y.left

	y.left = x.right
	if x.right != t.nil {
		x.right.p = y
	}

	x.right = y
	if y.p == t.nil {
		t.root = x
	} else if y.p.left == y {
		y.p.left = x
	} else {
		y.p.right = x
	}
	x.p = y.p
	y.p = x
}

func (t *RBTree) Insert(item Item) {
	if item == nil {
		panic("nil item is not allowed in RBTree")
	}

	insertLeft := true
	y := t.nil
	for x := t.root; x != t.nil; {
		y = x
		if item.Less(x.item) {
			x = x.left
			insertLeft = true
		} else if x.item.Less(item) {
			x = x.right
			insertLeft = false
		} else {
			x.item = item
			return
		}
	}

	z := t.freelist.newNode()
	z.item = item
	z.p = y
	if y == t.nil {
		t.root = z
	} else if insertLeft {
		y.left = z
	} else {
		y.right = z
	}
	z.left = t.nil
	z.right = t.nil
	z.color = RED
	t.insertFixup(z)

	t.length++
}

func (t *RBTree) insertFixup(z *node) {
	for z.p.color == RED {
		if z.p == z.p.p.left { // z的父节点是左节点
			y := z.p.p.right
			if y.color == RED { // case 1(a): z的叔节点是红
				z.p.color = BLACK
				y.color = BLACK
				z.p.p.color = RED
				z = z.p.p
			} else {
				if z == z.p.right { // case 2: z叔节点是黑色且z是是右孩子
					z = z.p
					t.leftRotate(z)
				}
				// case 3: z叔节点是黑色且z是左孩子
				z.p.color = BLACK
				z.p.p.color = RED
				t.rightRotate(z.p.p)
			}
		} else if z.p == z.p.p.right { // z的父节点是右节点
			y := z.p.p.left
			if y.color == RED { // case 1(b): z叔节点是红
				z.p.color = BLACK
				y.color = BLACK
				z.p.p.color = RED
				z = z.p.p
			} else {
				if z == z.p.left {
					z = z.p
					t.rightRotate(z)
				}
				z.p.color = BLACK
				z.p.p.color = RED
				t.leftRotate(z.p.p)
			}
		}
	}
	t.root.color = BLACK
}

// v替换u
func (t *RBTree) transplant(u *node, v *node) {
	if u.p == t.nil {
		t.root = v
	} else if u.p.left == u {
		u.p.left = v
	} else {
		u.p.right = v
	}
	v.p = u.p
}

func (t *RBTree) Search(item Item) Item {
	n := t.search(t.root, item)
	return n.item
}

func (t *RBTree) search(x *node, item Item) *node {
	for x != t.nil {
		if item.Less(x.item) {
			x = x.left
		} else if x.item.Less(item) {
			x = x.right
		} else {
			break
		}
	}
	return x
}

func (t *RBTree) Delete(item Item) (removeItem Item) {
	if item == nil {
		return nil
	}
	n := t.search(t.root, item)
	if n == t.nil {
		return nil
	}
	removeItem = n.item
	t.delete(n)
	t.freelist.freeNode(n)
	return
}

func (t *RBTree) delete(z *node) {
	var y = z
	yOriginalColor := y.color
	var x *node
	if z.left == t.nil {
		x = z.right
		t.transplant(z, z.right)
	} else if z.right == t.nil {
		x = z.left
		t.transplant(z, z.left)
	} else {
		y = t.minimum(z.right)
		yOriginalColor = y.color
		x = y.right
		if y.p == z {
			x.p = y // t.nil
		} else {
			t.transplant(y, y.right)
			y.right = z.right
			y.right.p = y
		}
		t.transplant(z, y)
		y.left = z.left
		y.left.p = y
		y.color = z.color
	}
	if yOriginalColor == BLACK {
		t.deleteFixup(x)
	}
	t.length--
}

func (t *RBTree) minimum(x *node) *node {
	for x.left != t.nil {
		x = x.left
	}
	return x
}

func (t *RBTree) maximum(x *node) *node {
	for x.right != t.nil {
		x = x.right
	}
	return x
}

func (t *RBTree) deleteFixup(x *node) {
	for x != t.root && x.color == BLACK {
		if x == x.p.left {
			w := x.p.right
			if w.color == RED { // case 1: x的兄弟节点w是红色
				w.color = BLACK
				x.p.color = RED
				t.leftRotate(x.p)
				w = x.p.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				// case 2: x的兄弟节点w是黑色的, 而且w的两个孩子都是黑色
				w.color = RED
				x = x.p
			} else {
				if w.right.color == BLACK {
					// case 3: x的兄弟节点w是黑色的, w的左孩子是红色, w的右孩子是黑色
					w.left.color = BLACK
					w.color = RED
					t.rightRotate(w)
					w = x.p.right
				}
				// case 4: x的兄弟节点w是黑色的, w的左孩子黑色, w的右孩子是红色
				w.color = x.p.color
				x.p.color = BLACK
				w.right.color = BLACK
				t.leftRotate(x.p)
				x = t.root
			}
		} else {
			w := x.p.left
			if w.color == RED {
				w.color = BLACK
				x.p.color = RED
				t.rightRotate(x.p)
				w = x.p.left
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.p
			} else {
				if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					t.leftRotate(w)
					w = x.p.left
				}
				w.color = x.p.color
				x.p.color = BLACK
				w.left.color = BLACK
				t.rightRotate(x.p)
				x = t.root
			}
		}
	}
	x.color = BLACK
}

func (t *RBTree) predecessor(x *node) *node {
	if x.left != t.nil {
		return t.maximum(x.left)
	}
	y := x.p
	for y != t.nil && x == y.left {
		x = y
		y = y.p
	}
	return y
}

func (t *RBTree) successor(x *node) *node {
	if x.right != t.nil {
		return t.minimum(x.right)
	}
	y := x.p
	for y != t.nil && x == y.right {
		x = y
		y = y.p
	}
	return y
}

func (t *RBTree) Len() int {
	return t.length
}

func (t *RBTree) morris(cur *node, f func(Item)) {
	// 如果cur无左孩子，cur向右移动（cur=cur.right）
	// 如果cur有左孩子，找到cur左子树上最右的节点，记为mostright
	// 如果mostright的right指针指向空，让其指向cur，cur向左移动（cur=cur.left）
	// 如果mostright的right指针指向cur，让其指向空，cur向右移动（cur=cur.right）
	var mostRight *node
	for cur != t.nil {
		mostRight = cur.left
		if mostRight != t.nil {
			for mostRight.right != t.nil && mostRight.right != cur {
				mostRight = mostRight.right
			}
			if mostRight.right == t.nil {
				mostRight.right = cur
				cur = cur.left
				continue
			} else {
				mostRight.right = t.nil
			}
		}
		f(cur.item)
		cur = cur.right
	}
}

func (t *RBTree) Ascend(f func(Item)) {
	t.morris(t.root, f)
}

func (t *RBTree) NewAscendIterator() Iterator {
	return Iterator{t: t, x: t.minimum(t.root)}
}

type Int int

func (a Int) Less(b Item) bool {
	return a < b.(Int)
}
