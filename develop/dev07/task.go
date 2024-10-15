package main

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

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
