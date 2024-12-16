package snowflakeid

import (
	"fmt"
	"os"

	"github.com/bwmarrin/snowflake"
)

var node, _ = snowflake.NewNode(int64(os.Getpid() % 1024))

func NewString() string {
	return fmt.Sprintf("%d", node.Generate())
}

func NewInt64() int64 {
	return node.Generate().Int64()
}

func NewInt() int {
	return int(node.Generate().Int64() % 2147483000)
}
