package grpc_server

type brodBool struct {
	lock chan bool
	val  bool
}

func newBroadcastBool() *brodBool {
	b := &brodBool{
		lock: make(chan bool, 1),
		val:  false,
	}
	b.lock <- true

	return b
}

func (b *brodBool) Set(status bool) {
	if status == true {
		if b.val == true {
			return
		}
		b.val = true
		<-b.lock
	} else {
		if b.val == false {
			return
		}
		b.lock <- true
		b.val = false
	}
}

func (b *brodBool) IsTrue() bool {
	return b.val
}

func (b *brodBool) WaitForTrue() {
	b.lock <- true
	defer func() { <-b.lock }()
}

func (b *brodBool) GetC() chan struct{} {
	b.lock <- true
	defer func() { <-b.lock }()

	// create temporary channel for selection
	var c chan struct{} = make(chan struct{}, 1)
	c <- struct{}{}
	return c
}
