package ecs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// Entity. See https://en.wikipedia.org/wiki/Entity_component_system
type Entity struct {
	UUID  uuid.UUID       `json:"uuid"`
	Store *ComponentStore `json:"components"`
}

// NewEntity creates an entity with a unique identifier
func NewEntity() *Entity {
	return &Entity{
		UUID:  uuid.New(),
		Store: &ComponentStore{},
	}
}

// ID returns the unique identifier for the entity
func (e *Entity) ID() uuid.UUID {
	return e.UUID
}

// Add a component to the entity. WARNING: This will not add the entity/component to the relevant systems. If you
// want to do this, use World.AddComponentToEntity() instead.
func (e *Entity) Add(component Component) {
	if reflect.TypeOf(component).Kind() != reflect.Ptr {
		panic(fmt.Sprintf("%#v is not a pointer - it is: %s", component, reflect.TypeOf(component).Kind()))
	}
	e.Store.Add(component)
}

// Component returns the first component matching the provided interface pointer, or nil.
func (e *Entity) Component(face interface{}) Component {
	interfaceType := reflect.TypeOf(face).Elem()
	for _, c := range e.Store.components {
		if reflect.TypeOf(c.Inner).Implements(interfaceType) {
			return c.Inner
		}
	}
	return nil
}

// Remove a component from the entity. WARNING: This will not remove the entity/component to the relevant systems.
// If you want to do this, use World.RemoveComponentFromEntity() instead.
func (e *Entity) Remove(component Component) {
	for i, c := range e.Store.components {
		if c.Inner == component {
			// copy whatever is at the end of the list to the position we're removing
			e.Store.components[i] = e.Store.components[len(e.Store.components)-1]
			// delete whatever is at the end of the list
			e.Store.components = e.Store.components[:len(e.Store.components)-1]
			return
		}
	}
}

func RemoveEntityFromSlice(slice []Entity, i int) []Entity {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

type ComponentStore struct {
	components []serialisableComponent
}

func (s *ComponentStore) Add(component interface{}) {
	s.components = append(s.components, serialisableComponent{
		Inner: component,
	})
}

func (s *ComponentStore) List() []interface{} {
	var list []interface{}
	for _, c := range s.components {
		list = append(list, c.Inner)
	}
	return list
}

type serialisableComponent struct {
	Inner interface{}
}

func (s *ComponentStore) MarshalJSON() ([]byte, error) {
	var parts []string
	for _, comp := range s.components {
		data, err := comp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		parts = append(parts, string(data))
	}
	return []byte("[" + strings.Join(parts, ",") + "]"), nil
}

func (s *ComponentStore) UnmarshalJSON(data []byte) error {
	var comps []savedComponent
	if err := json.Unmarshal(data, &comps); err != nil {
		return err
	}
	for _, c := range comps {
		empty, err := ComponentFromName(c.Type)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(c.Data, empty); err != nil {
			return err
		}
		s.Add(empty)
	}
	return nil
}

func (c *serialisableComponent) MarshalJSON() ([]byte, error) {
	componentData, err := json.Marshal(c.Inner)
	if err != nil {
		return nil, err
	}
	return json.Marshal(savedComponent{
		Type: reflect.TypeOf(c.Inner).Elem().Name(),
		Data: componentData,
	})
}

type savedComponent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
