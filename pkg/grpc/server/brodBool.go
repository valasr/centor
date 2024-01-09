package grpc_server

type brodBool struct {
	C   chan struct{}
	val bool
}

func newBroadcastBool() *brodBool {
	b := &brodBool{
		C:   make(chan struct{}, 1),
		val: false,
	}

	go func() {
		for {
			if b.val {
				b.C <- struct{}{}
			}
		}
	}()

	return b
}

func (b *brodBool) Set(status bool) {
	if status == true {
		if b.val == true {
			return
		}
		b.val = true
		b.C <- struct{}{}
	} else {
		if b.val == false {
			return
		}
		<-b.C
		b.val = false
	}
}

func (b *brodBool) IsTrue() bool {
	return b.val
}

func (b *brodBool) WaitForTrue() {
	<-b.C
}
