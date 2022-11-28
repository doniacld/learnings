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
	return false
}

func main() {
	t1 := tree.New(10)
	fmt.Println("t1: ", t1)
	ch := make(chan int)
	go Walk(t1, ch)

	for value := range ch {
		fmt.Println(value)
	}

	// My first naive solution without the close channel
	//for i := 0; i < 10; i++ {
	//	value, ok := <-ch
	//	if !ok {
	//		break
	//	}
	//	fmt.Println(value)
	//}
}
