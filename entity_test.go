package ecs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUIDIsSetOnCreation(t *testing.T) {
	assert.NotEqual(t, NewEntity().ID(), uuid.UUID{})
}

type TestComponent struct {
	X int
}

var IsTestable *Testable

type Testable interface {
	TestComponent() *TestComponent
}

func (c *TestComponent) TestComponent() *TestComponent {
	return c
}

func init() {
	RegisterComponent(&TestComponent{})
}

func TestAddingComponentToEntity(t *testing.T) {
	entity := NewEntity()
	component := TestComponent{
		X: 12345,
	}
	somethingElse := "blah"

	entity.Add(&component)
	entity.Add(&somethingElse)

	var testable *Testable
	matched := entity.Component(testable)

	require.NotNil(t, matched)

	assert.Equal(t, &component, matched.(Testable).TestComponent())
}

func TestRemovingComponentFromEntity(t *testing.T) {
	entity := NewEntity()
	component := TestComponent{
		X: 12345,
	}
	somethingElse := "blah"

	entity.Add(&component)
	entity.Add(&somethingElse)
	entity.Remove(&component)

	var testable *Testable
	matched := entity.Component(testable)

	require.Nil(t, matched)

	assert.Len(t, entity.Store.components, 1)
}

func TestRemovingComponentFromEntityLeavesNonMatchingComponentInPlace(t *testing.T) {
	entity := NewEntity()
	component := TestComponent{
		X: 12345,
	}
	somethingElse := "blah"

	entity.Add(&component)
	entity.Add(&somethingElse)
	entity.Remove(&somethingElse)

	var testable *Testable
	matched := entity.Component(testable)

	require.NotNil(t, matched)

	assert.Equal(t, &component, matched.(Testable).TestComponent())
	assert.Len(t, entity.Store.components, 1)
}

func TestEntityComponentSerialisation(t *testing.T) {
	entity := NewEntity()
	entity.Add(&TestComponent{
		X: 1,
	})
	data, err := json.Marshal(entity)
	if err != nil {
		t.Fatal(err)
	}
	var loaded Entity
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, entity, &loaded)
}

func BenchmarkAddingComponentToEntity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		entity := NewEntity()
		component := TestComponent{
			X: 12345,
		}
		b.StartTimer()
		entity.Add(component)
	}
}

func BenchmarkRemovingComponentFromEntity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		entity := NewEntity()
		component := TestComponent{
			X: 12345,
		}
		entity.Add(&component)
		b.StartTimer()
		entity.Remove(&component)
	}
}

func makeComponent() Component {
	return &struct {
		X int
	}{X: 123}
}

func BenchmarkRetrievingComponentFromEntityWith10OtherComponents(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var testable *Testable
		entity := NewEntity()
		component := TestComponent{
			X: 12345,
		}
		for i := 0; i < 10; i++ {
			entity.Add(makeComponent())
		}
		entity.Add(component)
		b.StartTimer()
		_ = entity.Component(testable)
	}
}
