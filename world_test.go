package ecs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

type TestSystem struct {
	updateCount     int
	addedEntities   []*Entity
	removedEntities []*Entity
	world           *World
}

func (s *TestSystem) Add(entity *Entity) {
	s.addedEntities = append(s.addedEntities, entity)
}

func (s *TestSystem) Update(elapsed float32) {
	s.updateCount++
}

func (s *TestSystem) Remove(entity *Entity) {
	s.removedEntities = append(s.removedEntities, entity)
}

func (s *TestSystem) RequiredTypes() []interface{} {
	var testable *Testable
	return []interface{}{
		testable,
	}
}

func (s *TestSystem) New(world *World) {
	s.world = world
}

func TestNewGetsCalledIfProvided(t *testing.T) {
	world := NewWorld()

	system := &TestSystem{}

	world.AddSystem(system)

	assert.Equal(t, world, system.world)
}

func TestSystemsAddedToWorldAreUpdated(t *testing.T) {
	world := NewWorld()

	system := &TestSystem{}

	world.AddSystem(system)
	world.Update(1)

	assert.Equal(t, 1, system.updateCount)
}

func TestEntitiesAreAddedToRelevantSystemsWhenAddedToWorld(t *testing.T) {

	world := NewWorld()

	system := &TestSystem{}
	world.AddSystem(system)

	e := NewEntity()
	e.Add(&TestComponent{})

	world.AddEntity(e)

	require.Len(t, system.addedEntities, 1)
	assert.Equal(t, e, system.addedEntities[0])

}

func TestEntitiesAreRemovedFromTheRelevantSystemsWhenRemovedFromWorld(t *testing.T) {

	world := NewWorld()

	system := &TestSystem{}
	world.AddSystem(system)

	e := NewEntity()
	e.Add(&TestComponent{})

	world.AddEntity(e)
	world.RemoveEntity(e)

	require.Len(t, system.removedEntities, 1)
	assert.Equal(t, e, system.removedEntities[0])

}

func TestComponentsAreAddedToEntitiesAndEntitiesToRelevantSystems(t *testing.T) {
	world := NewWorld()

	system := &TestSystem{}
	world.AddSystem(system)

	e := NewEntity()

	world.AddEntity(e)

	require.Len(t, system.addedEntities, 0)

	testComponent := &TestComponent{}
	world.AddComponentToEntity(testComponent, e)

	var testable *Testable
	matchedComponent := e.Component(testable)
	assert.Equal(t, testComponent, matchedComponent)

	require.Len(t, system.addedEntities, 1)
	assert.Equal(t, e, system.addedEntities[0])
}

func TestComponentsAreRemovedFromEntitiesAndEntitiesFromRelevantSystems(t *testing.T) {
	world := NewWorld()

	system := &TestSystem{}
	world.AddSystem(system)

	testComponent := &TestComponent{}

	e := NewEntity()
	e.Add(testComponent)

	world.AddEntity(e)

	world.RemoveComponentFromEntity(testComponent, e)

	var testable *Testable
	matchedComponent := e.Component(testable)
	assert.Nil(t, matchedComponent)

	require.Len(t, system.removedEntities, 1)
	assert.Equal(t, e, system.removedEntities[0])
}
