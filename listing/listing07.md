Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
1
3
4
3
5
7
8
6
0
0
...
```

Программа будет работать некорректно и зациклится, поэтому она не завершится и будет вечно ждать данные из закрытых каналов. Но перед этим она выведет значения из каналов a и b в произвольном порядке.

### Объяснение:
- Функция asChan: Эта функция принимает набор значений, запускает горутину и отправляет значения в канал с паузами между отправками (случайная задержка до 1000 миллисекунд). После завершения отправки значений канал закрывается.

- Функция merge: Функция запускает горутину, которая читает данные из двух каналов a и b и передаёт их в один канал c. Она использует select, чтобы поочерёдно считывать данные из этих каналов. Однако функция merge не закрывает канал c, а цикл продолжается бесконечно, даже после того как оба исходных канала будут закрыты.

- Основная программа (main): В main создаются два канала a и b, с различными последовательностями значений. Затем сливаются через merge. Программа выводит данные из объединённого канала c в цикле for v := range c. Поскольку канал c никогда не закрывается, цикл будет зациклен.

### Проблема:
Когда оба канала a и b будут закрыты, select в функции merge продолжит пытаться читать из закрытых каналов, и каждый раз будет возвращать нулевые значения. Так как канал c не закрывается, программа зациклится на выводе нулевых значений.

#### Что выводит программа:
Программа выводит числа от 1 до 8 в произвольном порядке (поскольку горутины выполняются с задержками), а затем зацикливается, выводя пустые строки или нулевые значения.

### Как исправить:
Чтобы программа корректно завершалась, нужно закрыть канал c после того, как оба исходных канала будут полностью обработаны.