package ecs

import (
	"reflect"

	"github.com/google/uuid"
)

type World struct {
	registrations []systemRegistration
	entities      []*Entity
	player        *Entity
}

type systemRegistration struct {
	system System
	types  []reflect.Type
}

func NewWorld() *World {
	return &World{}
}

func (w *World) GetEntity(id uuid.UUID) *Entity {
	for _, entity := range w.entities {
		if entity.ID() == id {
			return entity
		}
	}
	return nil
}

// repeatable refer st osystems that can be run without incrementing game state, i.e. renderers etc.
func (w *World) AddSystem(system System) {

	var refTypes []reflect.Type
	for _, t := range system.RequiredTypes() {
		refTypes = append(refTypes, reflect.TypeOf(t).Elem())
	}

	reg := systemRegistration{
		system: system,
		types:  refTypes,
	}

	for _, e := range w.entities {
		match := true
		for _, t := range reg.types {
			var found bool
			for _, c := range e.Store.components {
				found = reflect.TypeOf(c.Inner).Implements(t)
				//fmt.Printf("%v implements %v = %v\n", reflect.TypeOf(c.Inner).Name(), t.Name(), found)
				if found {
					break
				}
			}
			match = match && found
			if !match {
				break
			}
		}
		if match {
			reg.system.Add(e)
		}
	}

	w.registrations = append(w.registrations, reg)

	s, ok := system.(SystemInitializer)
	if ok {
		s.New(w)
	}
}

func (w *World) Update(elapsed float32) {
	for _, reg := range w.registrations {
		reg.system.Update(elapsed)
	}
}

func (w *World) AddEntity(e *Entity) {
	w.entities = append(w.entities, e)
	for _, reg := range w.registrations {
		match := true
		for _, t := range reg.types {
			var found bool
			for _, c := range e.Store.components {
				found = reflect.TypeOf(c.Inner).Implements(t)
				if found {
					break
				}
			}
			match = match && found
			if !match {
				break
			}
		}
		if match {
			reg.system.Add(e)
		}
	}
}

func (w *World) RemoveEntity(entity *Entity) {

	for i, e := range w.entities {
		if e == entity {
			w.entities[i] = w.entities[len(w.entities)-1]
			w.entities = w.entities[:len(w.entities)-1]
			break
		}
	}

	for _, reg := range w.registrations {
		reg.system.Remove(entity)
	}
}

func (w *World) ClearEntities() {
	tmp := make([]*Entity, len(w.entities))
	copy(tmp, w.entities)
	for _, entity := range tmp {
		w.RemoveEntity(entity)
	}
}

// AddComponentToEntity adds a given component to an entity. The component (c) must always be a struct pointer.
// If this change makes the entity a match for any previously uninvolved systems, it is added to those systems.
func (w *World) AddComponentToEntity(c interface{}, e *Entity) {

	e.Add(c)

	for _, reg := range w.registrations {
		match := true
		newTypeFound := false

		for _, t := range reg.types {
			var found bool

			for _, c := range e.Store.components {
				found = reflect.TypeOf(c.Inner).Implements(t)
				if found {
					break
				}
			}
			match = match && found
			if !match {
				break
			}
			if reflect.TypeOf(c).Implements(t) {
				newTypeFound = true
			}
		}
		if match && newTypeFound {
			reg.system.Add(e)
		}
	}
}

func (w *World) GetEntities() []*Entity {
	return w.entities
}

// RemoveComponentFromEntity removes a given component from an entity. The component must always be a struct pointer.
// If this change makes the entity a non-match for any previously matched systems, it is removed from those systems.
func (w *World) RemoveComponentFromEntity(c interface{}, e *Entity) {

	for _, reg := range w.registrations {
		match := true
		newTypeFound := false
		for _, t := range reg.types {
			var found bool
			for _, c := range e.Store.components {
				found = reflect.TypeOf(c.Inner).Implements(t)
				if found {
					break
				}
			}
			match = match && found
			if !match {
				break
			}
			if reflect.TypeOf(c).Implements(t) {
				newTypeFound = true
			}
		}
		if match && newTypeFound {
			reg.system.Remove(e)
		}
	}

	// we must remove after the above checks so the implements() checks still pass
	e.Remove(c)
}
