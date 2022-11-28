package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	walk(t, ch)
	close(ch) // close the channel when the walk of all the tree is done
}

// walk the tree recursively
func walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	ch <- t.Value

	walk(t.Left, ch)
	walk(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		// no more values in the tree
		if !ok1 && !ok2 {
			return true
		}

		// not the same length
		if ok1 != ok2 {
			return false
		}

		// not the same value
		if v1 != v2 {
			return false
		}
	}
}

func main() {
	t1 := tree.New(10)
	t2 := tree.New(5)

	fmt.Println(Same(t1, t1)) // true
	fmt.Println(Same(t1, t2)) // false
	fmt.Println(Same(t2, t1)) // false
}
