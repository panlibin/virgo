package virgo

const (
	bitsize = 32 << (^uint(0) >> 63)
	//maxint        = int(1<<(bitsize-1) - 1)
	maxintHeadBit = 1 << (bitsize - 2)
)

type taskQueue struct {
	q     []*task
	size  int
	r     int
	w     int
	mask  int
	empty bool
}

func newTaskQueue(size int) *taskQueue {
	size = ceilToPowerOfTwo(size)
	return &taskQueue{
		q:     make([]*task, size),
		size:  size,
		mask:  size - 1,
		empty: true,
	}
}

func (tq *taskQueue) pop() *task {
	if tq.empty {
		return nil
	}
	pTask := tq.q[tq.r]
	tq.r++
	if tq.r >= tq.size {
		tq.r = 0
	}
	if tq.r == tq.w {
		tq.empty = true
	}
	return pTask
}

func (tq *taskQueue) push(pTask *task) {
	if !tq.empty && tq.r == tq.w {
		tq.realloc()
	}

	tq.q[tq.w] = pTask
	tq.w++

	tq.empty = false
	if tq.w >= tq.size {
		tq.w = 0
	}
}

func (tq *taskQueue) length() int {
	if tq.empty {
		return 0
	}
	if tq.r < tq.w {
		return tq.w - tq.r
	} else {
		return tq.size - tq.r + tq.w
	}
}

func (tq *taskQueue) realloc() {
	ns := ceilToPowerOfTwo(tq.size + 1)
	nq := make([]*task, ns)
	nw := tq.length()
	if !tq.empty {
		if tq.r >= tq.w {
			copy(nq, tq.q[tq.r:])
			copy(nq[tq.size-tq.r:], tq.q[:tq.w])
		} else {
			copy(nq, tq.q[tq.r:tq.w])
		}
	}
	tq.r = 0
	tq.w = nw
	tq.size = ns
	tq.mask = ns - 1
	tq.q = nq
}

func ceilToPowerOfTwo(n int) int {
	if n&maxintHeadBit != 0 && n > maxintHeadBit {
		panic("argument is too large")
	}
	if n <= 2 {
		return 2
	}
	n--
	n = fillBits(n)
	n++
	return n
}

func fillBits(n int) int {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n
}
