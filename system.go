package ecs

type System interface {
	Add(entity *Entity)
	Update(world *World, player *Entity)
	Remove(entity *Entity)
	RequiredTypes() []interface{}
}
