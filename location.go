package main

import (
	"fmt"
	"slices"
	"strings"
)

func (l *location) ways() string {
	ways := make([]string, 0, len(l.exits))

	for _, way := range l.exits {
		ways = append(ways, way.label)
	}

	return " можно пройти - " + strings.Join(ways, ", ")
}

// возвращает предметы в помещении
// в каждом элементе - предметы на одном месте
//
// example:
// [<все предметы на столе>],
// [<все предметы на стуле>],
// [<все предметы где-то ещё>]
func (l *location) getPositionddItems() [][]*item {
	m := make(map[position][]*item, len(l.items))

	res := make([][]*item, 0)

	for _, item := range l.items {
		m[item.position] = append(m[item.position], item)
	}

	for _, item := range l.items {
		items, ok := m[item.position]
		if ok {
			res = append(res, items)
			delete(m, item.position)
		}
	}

	return res
}

func (l *location) canPutOn() []string {
	res := make([]string, 0, len(l.items))

	for _, item := range l.items {
		if item.pick == "надеть" {
			res = append(res, item.name)
		}
	}

	return res
}

func (l *location) canPutOnItem(thing string) bool {
	items := l.canPutOn()

	for _, item := range items {
		if item == thing {
			return true
		}
	}

	return false
}

func (l *location) punOnItem(thing string) string {
	if !l.canPutOnItem(thing) {
		return ""
	}

	var pos int
	for i, item := range l.items {
		if item.name == thing {
			pos = i
		}
	}

	l.items = append(l.items[:pos], l.items[pos+1:]...)

	return fmt.Sprintf("вы надели: %s", thing)
}

// получить полную строку с позицией вещей
//
// example:
// "на столе: ключи, конспекты, на стуле: рюкзак"
func (l *location) getFullItemsPosition() string {
	allItems := l.getPositionddItems()

	if len(allItems) == 0 {
		return ""
	}

	itemsByPosition := make([]string, 0, len(allItems))

	for _, items := range allItems {
		// "на столе"
		onWhat := items[0].position.onWhat

		// [ключи конспекты рюкзак]
		itemsStrings := make([]string, 0, len(items))
		for _, item := range items {
			itemsStrings = append(itemsStrings, item.name)
		}

		itemsByPosition = append(itemsByPosition, fmt.Sprintf("%s: %s", onWhat, strings.Join(itemsStrings, ", ")))
	}

	return strings.Join(itemsByPosition, ", ")
}

func (l *location) itemsText() string {
	itemsPosition := l.getFullItemsPosition()

	return itemsPosition
}

func (l *location) canTake() []string {
	res := make([]string, 0, len(l.items))

	for _, item := range l.items {
		if item.pick == Take {
			res = append(res, item.name)
		}
	}

	return res
}

func (l *location) canTakeItem(thing string) bool {
	items := l.canTake()

	return slices.Contains(items, thing)
}

func (l *location) takeItem(thing string) (string, *item) {
	if !l.canTakeItem(thing) {
		return "", nil
	}

	pos := slices.IndexFunc(l.items, func(item *item) bool { return item.name == thing })

	if pos < 0 {
		return "", nil
	}

	i := l.items[pos]

	l.items = slices.Delete(l.items, pos, pos+1)

	return fmt.Sprintf("предмет добавлен в инвентарь: %s", thing), i
}
