package main

import (
	"fmt"
)

type Nodes []string

// String returns the string representation of a node var
func (n *Nodes) String() string {
	return fmt.Sprintf("%s", *n)
}

// set appends the node to the node list
func (n *Nodes) Set(node string) error {
	*n = append(*n, node)
	return nil
}
