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