package ecs

import (
	"reflect"
	"testing"
)

func TestComponentDeserialisation(t *testing.T) {
	emptyComponent, err := ComponentFromName("TestComponent")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.TypeOf(emptyComponent).Implements(reflect.TypeOf(IsTestable).Elem()) {
		t.Fatalf("Component does not implement interface")
	}
}
