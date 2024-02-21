package main

import (
	"fmt"
)

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

func main() {
	h := newItem(5)
	h.appendElement(10)
	h.printList()
	h.appendElement(20)
	h.printList()
	h = h.Reverse()
	h.printList()
}
