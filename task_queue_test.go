package virgo

import (
	"fmt"
	"testing"
)

func TestTaskQueue(t *testing.T) {
	q := newTaskQueue(4)
	for i := 0; i < 4; i++ {
		q.push(&task{taskType: int32(i)})
	}

	for i := 0; i < 3; i++ {
		fmt.Println(q.pop().taskType)
	}
	for i := 4; i < 20; i++ {
		q.push(&task{taskType: int32(i)})
	}

	for t := q.pop(); t != nil; t = q.pop() {
		fmt.Println(t.taskType)
	}
}
