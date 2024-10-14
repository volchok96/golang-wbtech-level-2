package main

func or(doneChannels ...<-chan interface{}) <-chan interface{} {
	switch len(doneChannels) {
	case 0:
		return nil
	case 1:
		return doneChannels[0]
	}

	mergedDone := make(chan interface{})
	go func() {
		defer close(mergedDone)
		switch len(doneChannels) {
		case 2:
			select {
			case <-doneChannels[0]:
			case <-doneChannels[1]:
			}
		default:
			select {
			case <-doneChannels[0]:
			case <-doneChannels[1]:
			case <-doneChannels[2]:
			case <-or(append(doneChannels[3:], mergedDone)...):
			}
		}
	}()
	return mergedDone
}
