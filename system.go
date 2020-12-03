package ecs

type System interface {
	Add(entity *Entity)
	Update(elapsed float32)
	Remove(entity *Entity)
	RequiredTypes() []interface{}
}

type SystemInitializer interface {
	New(*World)
}
