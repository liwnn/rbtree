# RBtree
This Go package provides an implementation of red-black tree

## Usage
All you have to do is to implement a comparison `function Less() bool` for your Item which will be store in the red-black tree, here are some examples.

A case for `int` items.
``` go
package main

import (
    "github.com/liwnn/rbtree"
)

func main() {
	tr := rbtree.New()

	// Insert some values
	for i := 0; i < 10; i++ {
		tr.Insert(rbtree.Int(i))
	}

	// Get the value of the key
	item := tr.Search(rbtree.Int(5))
	if item != nil {
		fmt.Println(item)
	}

	// Delete the key
	if tr.Delete(rbtree.Int(4)) {
		fmt.Println("Deleted", 4)
	}

	// Traverse the tree
	for it := tr.NewAscendIterator(); it.Valid(); it.Next() {
		fmt.Println(it.Value().(rbtree.Int))
	}
}
```

A case for `struct` items:
``` go
package main

import (
    "github.com/liwnn/rbtree"
)

type KV struct {
	Key   int
	Value int
}

func (kv KV) Less(than rbtree.Item) bool {
	return kv.Key < than.(KV).Key
}

func main() {
	tr := rbtree.New()

	// Insert some values
	for i := 0; i < 10; i++ {
		tr.Insert(KV{Key: i, Value: 100 + i})
	}

	// Get the value of the key
	item := tr.Search(KV{Key: 1})
	if item != nil {
		fmt.Println(item.(KV))
	}

	// Delete the key
	if tr.Delete(KV{Key: 4}) {
		fmt.Println("Deleted", 4)
	}

	// Traverse the list
	for it := tr.NewAscendIterator(); it.Valid(); it.Next() {
		fmt.Println(it.Value().(KV))
	}
}
```
