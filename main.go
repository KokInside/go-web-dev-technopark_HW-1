package main

import (
	"fmt"
	"strings"
)

var u *you

const (
	Table   = "стол"
	Chair   = "стул"
	OnTable = "на столе"
	OnChair = "на стуле"
)

const (
	Room    = "комната"
	Kitchen = "кухня"
	Street  = "улица"
	Hallway = "коридор"
)

const (
	Tea      = "чай"
	Backpack = "рюкзак"
	Keys     = "ключи"
	Paper    = "конспекты"
)

const (
	Door = "дверь"
)

const (
	PutOn = "надеть"
	Take  = "взять"
)

type exit struct {
	label string // как выход назван в тексте: "улица", "домой"
	to    *location
	lock  *active // nil — проход свободен
}

type location struct {
	name    string
	exits   []exit
	items   []*item
	welcome func() string // текст при входе
	look    func() string // текст на «осмотреться»
}

type position struct {
	name   string
	onWhat string
}

type active struct {
	name      string
	trigger   string
	closedMsg string
	message   string
	opened    bool
}

type item struct {
	name     string
	position position
	pick     string
}

type you struct {
	where     *location
	inventory []*item
	hasBag    bool
}

func main() {
	/*
		в этой функции можно ничего не писать,
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/

	// initGame()

	// _ = handleCommand("some command")
}

func initGame() {
	keys := &item{
		name: Keys,
		position: position{
			name:   Table,
			onWhat: OnTable,
		},
		pick: Take,
	}
	paper := &item{
		name: Paper,
		position: position{
			name:   Table,
			onWhat: OnTable,
		},
		pick: Take,
	}
	pack := &item{
		name: Backpack,
		position: position{
			name:   Chair,
			onWhat: OnChair,
		},
		pick: PutOn,
	}
	tea := &item{
		name: Tea,
		position: position{
			name:   Table,
			onWhat: OnTable,
		},
	}

	door := &active{
		name:      Door,
		trigger:   Keys,
		opened:    false,
		closedMsg: "дверь закрыта",
		message:   "дверь открыта",
	}

	// locations
	kitchen := &location{
		name:  Kitchen,
		items: []*item{tea},
	}
	street := &location{
		name: Street,
	}
	room := &location{
		name:  Room,
		items: []*item{keys, paper, pack},
	}
	hallway := &location{
		name: Hallway,
	}

	// путь с кухни
	kitchen.exits = []exit{
		{label: Hallway, to: hallway},
	}

	// путь с улицы
	street.exits = []exit{
		{label: "домой", to: hallway},
	}

	// путь из комнаты
	room.exits = []exit{
		{label: Hallway, to: hallway},
	}

	// путь из коридора
	hallway.exits = []exit{
		{label: kitchen.name, to: kitchen},
		{label: room.name, to: room},
		{label: street.name, to: street, lock: door},
	}

	// коридор
	hallway.welcome = func() string {
		return "ничего интересного." + hallway.ways()
	}
	hallway.look = hallway.welcome

	// кухня
	kitchen.welcome = func() string {
		return "кухня, ничего интересного." + kitchen.ways()
	}

	kitchen.look = func() string {
		aim := "надо идти в универ"
		if !u.hasBag {
			aim = "надо собрать рюкзак и идти в универ"
		}
		return fmt.Sprintf("ты находишься на кухне, %s, %s.%s", kitchen.itemsText(), aim, kitchen.ways())
	}

	// комната
	room.welcome = func() string {
		return "ты в своей комнате." + room.ways()
	}

	room.look = func() string {
		items := room.itemsText()
		if items == "" {
			items = "пустая комната"
		}
		return fmt.Sprintf("%s.%s", items, room.ways())
	}

	// улица
	street.welcome = func() string {
		return "на улице весна." + street.ways()
	}
	street.look = street.welcome

	u = &you{
		where:     kitchen,
		inventory: nil,
	}
}

func handleCommand(command string) string {
	commands := strings.Fields(command)

	if len(commands) < 1 {
		return ""
	}

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
