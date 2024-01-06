package grpc_server

import "sync"

type brodBool struct {
	l   sync.Mutex
	val bool
}

func newBroadcastBool() *brodBool {
	b := &brodBool{
		l:   sync.Mutex{},
		val: false,
	}
	b.l.Lock()

	return b
}

func (b *brodBool) Set(status bool) {
	if status == true {
		if b.val == true {
			return
		}
		b.val = true
		b.l.Unlock()
	} else {
		if b.val == false {
			return
		}
		b.l.Lock()
		b.val = false
	}
}

func (b *brodBool) IsTrue() bool {
	return b.val
}

func (b *brodBool) WaitForTrue() {
	b.l.Lock()
	defer b.l.Unlock()
}
