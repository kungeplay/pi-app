package aigc

import (
	"fmt"
	"testing"
)

func TestChat(t *testing.T) {
	reply := Chat("你是谁")
	fmt.Println(reply)
}
