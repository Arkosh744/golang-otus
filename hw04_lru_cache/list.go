package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(item *ListItem)
	MoveToFront(item *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		item.Next = l.front
		l.front.Prev = item
		l.front = item
	}

	l.len++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		item.Prev = l.back
		l.back.Next = item
		l.back = item
	}

	l.len++

	return item
}

func (l *list) Remove(item *ListItem) {
	if item == nil {
		return
	}

	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		l.front = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		l.back = item.Prev
	}

	l.len--
}

func (l *list) MoveToFront(item *ListItem) {
	if item == nil || item == l.front {
		return
	}

	if item.Next == nil {
		l.back = item.Prev
	}

	l.Remove(item)
	l.PushFront(item.Value)
}
