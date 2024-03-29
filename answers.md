# Вопросы к собеседованию

## 1. Go — императивный или декларативный? А в чем разница?

Декларативный описывает "что нужно делать", но не описывает как, например SQL или map, filter в Python.
 
Императивный говорит "как делать".

__Пример:__
```python
a = filter(lambda x: x>2, [1, 3, 0]) #- декларативный подход
```
То же самое циклом - императивный

## 2. Что такое type switch?

Инструмент для выбора типа в рантайме. Пример из [A tour of a Go](https://go.dev/tour/methods/16)

```go
func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}
```

## 3. Как сообщить компилятору, что наш тип реализует интерфейс?

Реализовать все методы интерфейса. 

Если птица, плавает, летает и кричит как утка, значит это утка. - __duck typing__

## 4. Как работает append?
`Append` увеличивает `capacity` по следующим правилам:
1. Если требуемая емкость cap больше двойного размера старой емкости `old.cap`, то требуемая емкость cap будет использована в качестве новой newcap .

2. В противном случае, если старая емкость `old.cap` меньше 1024. Конечной емкостью newcap будет увеличение в 2 раза старой емкости `(old.cap)`, то есть newcap = doublecap

3. Если оба предыдущих условия не выполнены, а длина старого среза больше или равна 1024, окончательная емкость newcap начинается со старой емкости old.cap и циклически увеличивается на 1/4 от исходной, где `newcap = old.cap`, для `{newcap + = newcap / 4}` до тех пор, пока конечной емкостью newcap не станет емкость большая требуемой емкости cap, то есть newcap >= cap

[Источник](https://habr.com/ru/post/660827/)

## 5. Какое у slice zero value? Какие операции над ним возможны?

1. `nil`
2. `len()`, `cap()`
3. `append`
4. `make`
5. получение элемента, запиь в него, slice - [::]
6. `range`

## 6. Как устроен тип `map` в Go
[Разбор](https://habr.com/ru/post/457728/)

Map - это указатель на структуру `hmap` - хештаблица.

```go
// A header for a Go map.
type hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}
```
![Hash table](https://habrastorage.org/r/w1560/getpro/habr/post_images/979/e11/792/979e11792b1b87cc2a2548ebd3bd1743.png)

На картинке схематичное изображение структуры в памяти — есть хэдер hmap, указатель на который и есть map в Go (именно он создается при объявлении с помощью var, но не инициализируется, из-за чего падает программа при попытке вставки). Поле buckets — хранилище пар ключ-значение, таких «ведер» несколько, в каждом лежит 8 пар. Сначала в «ведре» лежат слоты для дополнительных битов хэшей (e0..e7 названо e — потому что extra hash bits). Далее лежат ключи и значения как сначала список всех ключей, потом список всех значений.

По хэш функции определяется в какое «ведро» мы кладем значение, внутри каждого «ведра» может лежать до 8 коллизий, в конце каждого «ведра» есть указатель на дополнительное, если вдруг предыдущее переполнилось. 

### 7. Каков порядок перебора map?
Порядок перебора не гарантируется. (чаще всего это будет разная последовательность). Связано с устройством мапы, процессом переноса данных из одного в другой бакет.

__как получить одно случайное значение из map?__
```go
for _, v := range myMap {
	fmt.Println(v)
	break
}
```
### 8. Что будет, если читать из закрытого канала?
Мы получим zero value того типа, которого объявлен канал. 

Канал возвращает одно значение при чтении, поэтому мы не можем получить ошибку.

### 9. Что будет, если писать в закрытый канал?
Так как мы не можем получить ошибку, то будет паника.

### 10. Как вы отсортируете массив структур по алфавиту по полю Name?
При помощи стандартной сортировки из пакета `sort` и метода `sort.Slice` (`func Slice(x any, less func(i, j int) bool)`) или `SliceStable`:
```go
type data struct {
	Name string
	age  int
}

func main() {
	a := []data{{"c", 10}, {"b", 15}, {"aa", 25}}
	sort.Slice(a, func(i, j int) bool {
		return a[i].Name < a[j].Name
	})
	fmt.Println(a)
}
```
### 11. Что такое сериализация? Зачем она нужна?

Сериализация позволяет нам сохранять состояние объекта и воссоздавать объект в новом месте. Сериализация охватывает как хранение объекта, так и обмен данными. Поскольку объекты состоят из нескольких компонентов, сохранение или доставка всех частей обычно требует значительных усилий по написанию кода, поэтому сериализация — это стандартный способ зафиксировать объект в совместно используемом формате.

[варианты сериализации в Go](https://developpaper.com/several-ways-of-serialization-and-deserialization-of-golang/)

### 12. Сколько времени в минутах займет у вас написание процедуры обращения односвязного списка?
10 мин
```go
type item struct {
	data int
	next *item
}

func (head *item) appendElement(element int) {
	p := head
	for p.next != nil {
		p = p.next
	}
	p.next = newItem(element)
}

func newItem(data int) *item {
	return &item{
		data: data,
	}
}

func (head *item) printList() {
	for p := head; p != nil; p = p.next {
		fmt.Printf("%d ", p.data)
	}
	fmt.Println()
}

func (head *item) Reverse() *item {
	var prev *item = nil
	current := head
	var next *item = nil
	for current != nil {
		next = current.next
		current.next = prev
		prev = current
		current = next
	}
	return prev
}
```

### 13. Где следует поместить описание интерфейса: в пакете с реализацией или в пакете, где этот интерфейс используется? Почему?

`tight coupling` - Тесная связь — это метод связи, при котором аппаратные и программные компоненты сильно зависят друг от друга. Он используется для обозначения состояния/назначения взаимосвязи между двумя или более вычислительными экземплярами в интегрированной системе. Жесткая связь также известна как высокая связь и сильная связь.

Оба варианта рабочие. Мы разделяем интерфейс и реализацию для уменьшения связности приложения.

### 14. Предположим, ваша функция должна возвращать детализированные Recoverable и Fatal ошибки. Как это реализовано в пакете net? Как это надо делать в современном Go?

__что нового и важного в плане обработки ошибок появилось в Go 1.13?__
в 1.13 в пакете `errors` появились функции:
* `errors.Unwrap()`
* `errors.Is()`
* `errors.As()`
* `fmt.Errorf("error 1: %w", error2)`

Что позволет глубже сопоставлять ошибки и отслеживать цепочки ошибок.

### 15. Главный недостаток стандартного логгера?

Главный недостаток стандартного логгера в Go заключается в его относительной простоте и ограниченности функционала. Хотя стандартный пакет log предоставляет базовые возможности для журналирования, такие как вывод сообщений с временем и возможностью перенаправления вывода, он имеет несколько существенных недостатков для более сложных или масштабируемых приложений:

* __Отсутствие уровней логирования:__ В стандартном логгере нет встроенной поддержки уровней логирования (например, DEBUG, INFO, WARN, ERROR), которые позволяют более гибко управлять выводом логов в зависимости от их важности и контекста использования.
* __Ограниченная настройка формата вывода:__ Хотя можно изменить префикс и флаги, определяющие, что включается в каждое лог-сообщение (например, дата, время), возможности кастомизации формата сообщения довольно ограничены по сравнению с более продвинутыми библиотеками логирования.
* __Невозможность фильтрации или перенаправления логов:__ В стандартной библиотеке нет встроенных средств для фильтрации логов по уровням или их перенаправления в различные назначения (например, файлы, сетевые сервисы) без дополнительной реализации со стороны разработчика.
* __Отсутствие структурированного логирования:__ Стандартный логгер не поддерживает структурированное логирование (например, в формате JSON), которое облегчает анализ логов, особенно в масштабируемых и распределенных системах.
* 
Из-за этих недостатков разработчики часто прибегают к использованию сторонних библиотек логирования, таких как logrus, zap или zerolog, которые предлагают более продвинутые возможности, включая уровни логирования, структурированный вывод, а также гибкую настройку формата и направления вывода логов.

### 17. Какой у вас любимый линтер?
 - govet
 - errcheck
 - funlen

__какое отношение линтеры имеют к CI? Зачем нужен CI в процессе разработки?__

Чтобы привести кодовую базу к единому стандарту. 


### 18. Можно ли использовать один и тот же буфер []byte в нескольких горутинах?

???


### 19. Какие типы мьютексов предоставляет stdlib?

Мььютексы доступны в пакете `sync`
* `sync.Mutex`
* `sync.RWMutex`

Обычный мютекс блокирует объект для доступа из одной горутины.

`RWMutex` позволяет развести процессы чтения и записи и в момент чтения не блокировать доступ к объекту, блокировать объект только на запись
`func (rw *RWMutex) Lock()` - Lock locks rw for writing
`func (rw *RWMutex) RLock()` - RLock locks rw for reading. 

__что именно и от чего защищает мьютекс?__

От состояния гонки.

### 20. Что такое lock-free структуры данных, и есть ли в Go такие?
[Безблокировочное программирование](https://habr.com/ru/company/wunderfund/blog/322094/)

В Go не приветствуется подход с общим доступам к данным.

`Share memory by communicating; don't communicate by sharing memory. `

 Но он возможен. Для этого есть в пакете sync мьютексы, для ограничения доступа к данным, есть sync.Map - структура с возможностью конкурентного чтения и записи. 

Внутри sync.Map доступ реализован при помощи мьютексов.

Есть пакет atomic, который позволяет выполнять "атомарные" (примитивные) операции для конкурентной работы с данными.

### 21. Способы поиска проблем производительности на проде?

Не было опыта 😔

### 22. Стандартный набор метрик prometheus в Go -программе?

Для сервисов можно исполььзовать RED-метрики:
* Request - число запросов
* Error - количество ошибок
* Duration - время задержки

### 23. Как встроить стандартный профайлер в свое приложение?

???

### 24. Overhead от стандартного профайлера?
???

### 25. Почему встраивание — не наследование?
В Go, встраивание структур часто используется как средство для достижения эффекта, похожего на наследование в объектно-ориентированных языках программирования, но на самом деле это не наследование в традиционном понимании. Вот основные причины, по которым встраивание в Go отличается от наследования:

Плоская структура: В отличие от наследования, которое образует иерархическую структуру между базовыми и производными классами, встраивание в Go создает плоскую структуру. Это означает, что методы и свойства встроенной структуры становятся частью внешней структуры без создания явной иерархии.
Композиция вместо расширения: Встраивание в Go скорее является формой композиции, чем наследования. Вместо того чтобы "быть" подтипом, структура "имеет" встроенную структуру, делая все её поля и методы доступными. Это соответствует принципу "предпочитай композицию наследованию", который поощряется во многих парадигмах программирования для более гибкой и управляемой архитектуры.
Нет переопределения методов: В наследовании производные классы могут переопределять методы базового класса. В Go, если встроенная структура и внешняя структура имеют методы с одинаковым именем, оба метода сохраняются и становятся доступными через внешнюю структуру, но это не переопределение в традиционном смысле.
Типы и интерфейсы: Go использует интерфейсы для достижения полиморфизма. Вместо того чтобы полагаться на иерархию классов, как в наследовании, полиморфизм в Go достигается за счет реализации интерфейсов. Это позволяет объектам разных типов быть использованными в одном и том же контексте, если они реализуют один и тот же интерфейс, без необходимости явного наследования.
Простота и ясность: Использование встраивания вместо наследования подчеркивает простоту и ясность в дизайне программ на Go. Это предотвращает некоторые сложности и ограничения, связанные с наследованием, такие как хрупкая база классов и проблемы с множественным наследованием.
Таким образом, встраивание в Go предоставляет гибкий и мощный механизм для организации кода, который отличается от классического наследования и нацелен на упрощение проектирования и повышение читаемости кода.
* __S__ _The single-responsibility principle:_ "There should never be more than one reason for a class to change." In other words, every class should have only one responsibility
* __O__ _The open–closed principle_: "Software entities ... should be open for extension, but closed for modification."
* __L__ _The Liskov substitution principle:_ "Functions that use pointers or references to base classes must be able to use objects of derived classes without knowing it."
* __I__ _The interface segregation principle:_ "Many client-specific interfaces are better than one general-purpose interface."
* __D__ _The dependency inversion principle:_ "Depend upon abstractions, not concretions."

Принцип Лисков 
```
если S является подтипом T, тогда объекты типа T в программе могут быть замещены объектами типа S без каких-либо изменений желательных свойств этой программы
```

Встраивание - это не наследование, а композиция.
Аггрегация - это ссылка на объект внутри структуры.

### 26. Какие средства обобщенного программирования есть в Go?

Пустые интерфейсы `interface{}`

### 27. Какие технологические преимущества языка Go вы можете назвать?

* Нативная конкурентность
* Пакет net, который позволяет легко делать веб-сурвисы из коробки
* Кроссплатформенность.
* Высокая производительность.
* гарбэдж коллектор

### 28. Какие технологические недостатки языка Go вы можете назвать?

* Отсутствие привычных стандартных функций типа abs, для встроенных типов.
* гарбэдж коллектор
* отсутсвие нативного приведения сравнимых типов, например нельзя сравнить int, int64
