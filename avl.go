// Public domain

package avl

// Bounded recursive AVL tree adapted from Julienne Walker's presentation at
// http://eternallyconfuzzled.com/tuts/datastructures/jsw_tut_avl.aspx.

// The Key interface must be supported by data stored in the AVL tree.
type Key interface {
	Less(Key) bool
	Equal(Key) bool
	Greater(Key) bool
}

// Tree is an AVL tree.  The zero value is an empty tree.
type Tree struct{ root *node }

// Insert inserts a node into the tree.
//
// The node is inserted even if other nodes with the same key already exist.
func (t *Tree) Insert(node Key) {
	t.root, _ = t.root.insertR(node)
}

// Remove a single node from an AVL tree.
//
// If key does not exist, the method has no effect.
func (t *Tree) Remove(key Key) {
	t.root, _ = t.root.removeR(key)
}

// Find finds a node in the tree by key.Eq.
//
// It returns nil if there is no match.
func (t *Tree) Find(key Key) Key {
	for node := t.root; node != nil; {
		switch {
		case key.Less(node.Key):
			node = node.link[0]
		case key.Greater(node.Key):
			node = node.link[1]
		default:
			return node.Key
		}
	}
	return nil
}

// node is a node in an AVL tree.
type node struct {
	Key              // anything comparable with Less and Eq.
	link    [2]*node // children, indexed by "direction", 0 or 1.
	balance int8     // balance factor
}

// A little readability function for returning the opposite of a direction,
// where a direction is 0 or 1.  Go inlines this.
// Where JW writes !dir, this code has opp(dir).
func opp(dir int8) int8 {
	return 1 - dir
}

// single rotation
func (root *node) single(dir int8) *node {
	save := root.link[opp(dir)]
	root.link[opp(dir)] = save.link[dir]
	save.link[dir] = root
	return save
}

// double rotation
func (root *node) double(dir int8) *node {
	save := root.link[opp(dir)].link[dir]

	root.link[opp(dir)].link[dir] = save.link[opp(dir)]
	save.link[opp(dir)] = root.link[opp(dir)]
	root.link[opp(dir)] = save

	save = root.link[opp(dir)]
	root.link[opp(dir)] = save.link[dir]
	save.link[dir] = root
	return save
}

// adjust balance factors after double rotation
func (root *node) adjustBalance(dir, bal int8) {
	n := root.link[dir]
	nn := n.link[opp(dir)]
	switch nn.balance {
	case 0:
		root.balance = 0
		n.balance = 0
	case bal:
		root.balance = -bal
		n.balance = 0
	default:
		root.balance = 0
		n.balance = bal
	}
	nn.balance = 0
}

func (root *node) insertBalance(dir int8) *node {
	n := root.link[dir]
	bal := 2*dir - 1
	if n.balance == bal {
		root.balance = 0
		n.balance = 0
		return root.single(opp(dir))
	}
	root.adjustBalance(dir, bal)
	return root.double(opp(dir))
}

func (root *node) insertR(data Key) (*node, bool) {
	if root == nil {
		return &node{Key: data}, false
	}
	var dir int8
	if root.Less(data) {
		dir = 1
	}
	var done bool
	if root.link[dir], done = root.link[dir].insertR(data); done {
		return root, true
	}
	root.balance += 2*dir - 1
	switch root.balance {
	case 0:
		return root, true
	case 1, -1:
		return root, false
	}
	return root.insertBalance(dir), true
}

func (root *node) removeBalance(dir int8) (*node, bool) {
	n := root.link[opp(dir)]
	bal := 2*dir - 1
	switch n.balance {
	case -bal:
		root.balance = 0
		n.balance = 0
		return root.single(dir), false
	case bal:
		root.adjustBalance(opp(dir), -bal)
		return root.double(dir), false
	}
	root.balance = -bal
	n.balance = bal
	return root.single(dir), true
}

func (root *node) removeR(key Key) (*node, bool) {
	if root == nil {
		return nil, false
	}
	if key.Equal(root.Key) {
		switch {
		case root.link[0] == nil:
			return root.link[1], false
		case root.link[1] == nil:
			return root.link[0], false
		}
		heir := root.link[0]
		for heir.link[1] != nil {
			heir = heir.link[1]
		}
		root.Key = heir.Key
		key = heir.Key
	}
	var dir int8
	if key.Greater(root.Key) {
		dir = 1
	}
	var done bool
	if root.link[dir], done = root.link[dir].removeR(key); done {
		return root, true
	}
	root.balance += 1 - 2*dir
	switch root.balance {
	case 1, -1:
		return root, true
	case 0:
		return root, false
	}
	return root.removeBalance(dir)
}
