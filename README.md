# WB Tech: Level 2 (Golang)
## Описание проекта
Этот репозиторий содержит решения задач второго уровня от WB Tech с использованием языка программирования Go (Golang). Каждое задание представляет собой файл с полностью оформленным решением, комментариями к коду и объяснением подхода.

## Структура проекта
Каждая задача оформлена в виде отдельного модуля или функции. Решения включают все необходимые комментарии, разъяснения и тесты для подтверждения корректности работы программ. В решениях применяются паттерны проектирования, а также учитываются требования по использованию встроенных библиотек и возможностей языка Go.

## Задания
### Паттерны проектирования
Реализованы паттерны:

- Фасад
- Строитель
- Посетитель
- Команда
- Цепочка вызовов
- Фабричный метод
- Стратегия
- Состояние

### Пример программы с использованием NTP
Программа печатает точное время, используя NTP-библиотеку, с корректной обработкой ошибок и выводом их в STDERR.

### Распаковка строки
Программа для распаковки строки с поддержкой escape-последовательностей и обработкой ошибок для некорректного ввода.

### Утилита sort
Утилита для сортировки строк в файле с поддержкой различных флагов, таких как сортировка по числовым значениям, обратный порядок, уникальные строки и другие параметры.

### Поиск анаграмм
Функция для поиска множеств анаграмм в массиве слов на русском языке, с приведением к нижнему регистру и сортировкой.

### Утилита grep
Утилита для фильтрации строк с поддержкой опций: количество строк до и после совпадений, игнорирование регистра, точное совпадение и другие.

### Утилита cut
Утилита для выборки колонок из строк, разбитых по разделителю с поддержкой различных опций.

### Or channel
Функция, объединяющая один или более done-каналов в единый канал, с поддержкой неизвестного количества каналов.

### UNIX-шелл
Собственная реализация Unix-оболочки с поддержкой команд: смена директории, вывод пути, вывод аргументов, убийство процесса, список процессов, а также поддержка конвейеров (pipes).

### Утилита wget
Утилита для скачивания сайтов целиком.

### Telnet-клиент
Простейший Telnet-клиент для подключения к TCP-серверам и обмена данными через сокет.

### HTTP-сервер для календаря
Реализован HTTP-сервер для работы с календарём, поддерживающий методы для создания, обновления, удаления и получения событий.

## Запуск
Для запуска каждого задания в локальной среде необходимо:

1. Установить Go.
2. Склонировать репозиторий.
```sh
git clone https://github.com/volchok96/golang-wbtech-level-2
```
3.Перейти в директорию соответствующего задания.
```sh
cd golang-wbtech-level-2
```
3. Выполнить команду go run <путь_к_файлу/имя_файла.go> для запуска программы.
```sh
go run develop/dev01/task.go
```

Каждое задание оформлено как отдельный Go-модуль или файл и готово к запуску в соответствии с инструкциями в комментариях к коду.

## Тестирование
Каждое решение в develop/ директории сопровождается тестами. Тесты можно запустить с помощью команды:
```sh
cd develop && go test ./... && cd ..
```

## Линтеры
Проект проходит проверки go vet и golint.
```sh
go vet ./...
golint ./...
```

## Лицензия
Этот проект лицензирован под MIT License.

## Контакты

Если у вас есть вопросы или предложения, пожалуйста, свяжитесь со мной:
- Email: kzakharova96@yandex.com
- TG: https://t.me/volchok_96

## Благодарности
Спасибо всем, кто внес свой вклад в развитие этого проекта!

### Идея проекта:
- https://tech.wildberries.ru/

### Техническое сопровождение:
- https://21-school.ru/




