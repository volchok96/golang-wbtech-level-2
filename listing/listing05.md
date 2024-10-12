Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
ok
```

### Объяснение:
Программа выводит "ok", потому что функция test возвращает nil. Переменная err типа интерфейса error присваивает себе этот nil, и при проверке if err != nil условие не выполняется, так как внутри интерфейса нет реального значения ошибки, только nil.
