package main

import (
	"fmt"
	"slices"
	"strings"
)

var u *you

const (
	table   = "стол"
	chair   = "стул"
	onTable = "на столе"
	onChair = "на стуле"
)

const (
	room    = "комната"
	kitchen = "кухня"
	street  = "улица"
	hallway = "коридор"
)

const (
	tea      = "чай"
	backpack = "рюкзак"
	keys     = "ключи"
	paper    = "конспекты"
)

const (
	door = "дверь"
)

const (
	putOn = "надеть"
	take  = "взять"
)

type position struct {
	name   string
	onWhat string
}

type active struct {
	name     string
	trigger  string
	message  string
	isActive bool
}

type item struct {
	name     string
	position position
	pick     string
	actOn    []*active
}

type you struct {
	where     *location
	inventory []*item
}

type location struct {
	name     string
	to       []*location
	items    []*item
	actions  []*active
	obstacle []*active
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

// " можно пройти - ..."
func (l *location) getFullWaysToGo() string {
	return fmt.Sprintf(" можно пройти - %s",
		strings.Join(l.canGo(), ", "))
}

func (l *location) canGo() []string {
	res := make([]string, 0, len(l.to))
	for _, to := range l.to {
		res = append(res, to.name)
	}

	return res
}

func (l *location) look(inventory []*item) string {
	switch l.name {
	case kitchen:
		base := "ты находишься на кухне, %s%s.%s"

		itemsPosition := l.getFullItemsPosition()
		if itemsPosition == "" {
			itemsPosition = "на кухне пусто"
		}

		// ", надо собрать рюкзак и идти в универ."

		aims := make([]string, 0)
		if inventory == nil {
			aims = append(aims, "собрать рюкзак")
		}
		aims = append(aims, "идти в универ")

		var aimString string
		if len(aims) != 0 {
			aimString = fmt.Sprintf(", надо %s", strings.Join(aims, " и "))
		}

		return fmt.Sprintf(base,
			itemsPosition,
			aimString,
			l.getFullWaysToGo(),
		)
	case street:
		return "на улице весна. можно пройти - домой"
	case room:
		base := "%s.%s"

		itemsPosition := l.getFullItemsPosition()
		if itemsPosition == "" {
			itemsPosition = "пустая комната"
		}

		return fmt.Sprintf(base,
			itemsPosition,
			l.getFullWaysToGo(),
		)

	case hallway:
	}

	return ""
}

func (l *location) canGoTo(to string) bool {
	ways := l.canGo()

	return slices.Contains(ways, to)
}

func (l *location) goTo(to string) string {
	if !l.canGoTo(to) {
		return ""
	}

	destination := l.getTo(to)

	l = destination

	switch to {
	case hallway:
		base := "ничего интересного.%s"

		return fmt.Sprintf(base, l.getFullWaysToGo())

	case room:
		base := "ты в своей комнате.%s"

		return fmt.Sprintf(base, l.getFullWaysToGo())

	case street:
		return "на улице весна. можно пройти - домой"

	case kitchen:
		base := "кухня, ничего интересного.%s"

		return fmt.Sprintf(base, l.getFullWaysToGo())
	}

	_ = to
	return ""
}

func (l *location) getTo(to string) *location {
	if !l.canGoTo(to) {
		return nil
	}

	for _, way := range l.to {
		if way.name == to {
			return way
		}
	}

	return nil
}

func (l *location) takeItem(thing string) (string, *item) {
	if !l.canTakeItem(thing) {
		return "", nil
	}

	var pos int
	for i, item := range l.items {
		if item.name == thing {
			pos = i
		}
	}

	i := l.items[pos]

	l.items = append(l.items[:pos], l.items[pos+1:]...)

	return fmt.Sprintf("предмет добавлен в инвентарь: %s", thing), i

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

func (l *location) canTake() []string {
	res := make([]string, 0, len(l.items))

	for _, item := range l.items {
		res = append(res, item.name)
	}

	return res
}

func (l *location) hasObstacles() bool {
	for _, i := range l.obstacle {
		if i.isActive == false {
			return true
		}
	}

	return false
}

func (l *location) canTakeItem(thing string) bool {
	items := l.canTake()

	return slices.Contains(items, thing)
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

func (u *you) look() string {
	return u.where.look(u.inventory)
}

func (u *you) goTo(to string) string {
	res := u.where.goTo(to)
	if res == "" {
		return fmt.Sprintf("нет пути в %s", to)
	}

	destination := u.where.getTo(to)

	if destination.hasObstacles() {
		return "дверь закрыта"
	}

	if destination != nil {
		u.where = destination
	}

	return res
}

func (u *you) take(thing string) string {
	if u.inventory == nil {
		return "некуда класть"
	}

	if !u.where.canTakeItem(thing) {
		return "нет такого"
	}

	res, item := u.where.takeItem(thing)
	if res == "" {
		return "некуда класть"
	}

	u.inventory = append(u.inventory, item)

	return res
}

func (u *you) putOn(thing string) string {
	if !u.where.canTakeItem(thing) {
		return "нет такого"
	}

	res := u.where.punOnItem(thing)
	if res == "" {
		return "нет такого"
	}

	switch thing {
	case backpack:
		u.inventory = make([]*item, 0)
	}

	return res
}

func (u *you) contains(thing string) bool {
	for _, i := range u.inventory {
		if i.name == thing {
			return true
		}
	}

	return false
}

func (l *location) canAct() []*active {
	return l.actions
}

func (l *location) canActOn(actOn string) bool {
	actions := l.canAct()

	for _, active := range actions {
		if active.name == actOn {
			return true
		}
	}

	return false
}

// func (l *location) act(thing, on string) string {
// 	return ""
// }

// применить <что> над <чем>
func (u *you) act(thing, on string) string {
	if !u.contains(thing) {
		return fmt.Sprintf("нет предмета в инвентаре - %s", thing)
	}

	// if !u.where.canActOn(on) {
	// 	return "не к чему применить"
	// }

	// u.where.act()

	switch on {
	case "дверь":
		var action *active

		for _, i := range u.where.actions {
			if i.name == on {
				action = i
				if action.trigger == thing {
					action.isActive = true
				} else {
					return "не к чему применить"
				}
			}
		}

		if action == nil {
			return "не к чему применить"
		}

		return "дверь открыта"
	default:
		return "не к чему применить"
	}
}

func main() {

	/*
		в этой функции можно ничего не писать,
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/

	//initGame()

	//_ = handleCommand("some command")
}

func initGame() {
	keys := &item{
		name: "ключи",
		position: position{
			name:   table,
			onWhat: onTable,
		},
		pick: "взять",
	}
	paper := &item{
		name: "конспекты",
		position: position{
			name:   table,
			onWhat: onTable,
		},
		pick: "взять",
	}
	pack := &item{
		name: "рюкзак",
		position: position{
			name:   chair,
			onWhat: onChair,
		},
		pick: "надеть",
	}
	tea := &item{
		name: "чай",
		position: position{
			name:   table,
			onWhat: onTable,
		},
	}

	door := &active{
		name:     door,
		trigger:  "ключи",
		isActive: false,
	}

	kitchen := &location{
		name:  "кухня",
		items: []*item{tea},
	}
	street := &location{
		name:     "улица",
		obstacle: []*active{door},
		// actions:  []*active{door},
		// obstacle: []*active{door},
	}
	room := &location{
		name:  "комната",
		items: []*item{keys, paper, pack},
	}
	hallway := &location{
		name:    "коридор",
		actions: []*active{door},
	}

	kitchen.to = append(kitchen.to, hallway)

	room.to = append(room.to, hallway)

	street.to = append(street.to, hallway)

	hallway.to = append(hallway.to, kitchen, room, street)

	u = &you{
		where:     kitchen,
		inventory: nil,
	}

}

func handleCommand(command string) string {
	commands := strings.Fields(command)

	switch commands[0] {
	case "осмотреться":
		if len(commands) != 1 {
			return "неизвестная команда"
		}
		return u.look()
	case "идти":
		if len(commands) != 2 {
			return "неизвестная команда"
		}
		return u.goTo(commands[1])
	case "взять":
		if len(commands) != 2 {
			return "неизвестная команда"
		}
		return u.take(commands[1])
	case "надеть":
		if len(commands) != 2 {
			return "неизвестная команда"
		}
		return u.putOn(commands[1])
	case "применить":
		if len(commands) != 3 {
			return "неизвестная команда"
		}
		return u.act(commands[1], commands[2])
	}

	return "неизвестная команда"
}
