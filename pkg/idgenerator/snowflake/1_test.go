package snowflake

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	node, _ := NewNode(1)
	fmt.Println(node.Generate().String())
}
