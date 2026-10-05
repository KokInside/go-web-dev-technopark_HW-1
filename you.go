package main

import "slices"

func (u *you) contains(thing string) bool {
	i := slices.IndexFunc(u.inventory, func(i *item) bool { return i.name == thing })

	return i >= 0
}

func (u *you) goTo(label string) string {
	i := slices.IndexFunc(u.where.exits, func(e exit) bool { return e.label == label })
	if i < 0 {
		return "нет пути в " + label
	}

	e := u.where.exits[i]
	if e.lock != nil && !e.lock.opened {
		return e.lock.closedMsg // "дверь закрыта"
	}

	u.where = e.to
	return e.to.welcome()
}

func (u *you) act(thing, on string) string {
	if !u.contains(thing) {
		return "нет предмета в инвентаре - " + thing
	}

	for _, e := range u.where.exits {
		if e.lock != nil && e.lock.name == on && e.lock.trigger == thing {
			e.lock.opened = true
			return e.lock.message // "дверь открыта"
		}
	}

	return "не к чему применить"
}

func (u *you) look() string {
	return u.where.look()
}

func (u *you) take(thing string) string {
	if !u.hasBag {
		return "некуда класть"
	}

	if !u.where.canTakeItem(thing) {
		return "нет такого"
	}

	res, item := u.where.takeItem(thing)
	if res == "" {
		return "нет такого"
	}

	u.inventory = append(u.inventory, item)

	return res
}

func (u *you) putOn(thing string) string {
	if !u.where.canPutOnItem(thing) {
		return "нет такого"
	}

	res := u.where.punOnItem(thing)
	if res == "" {
		return "нет такого"
	}

	if thing == Backpack {
		u.hasBag = true
	}

	return res
}
