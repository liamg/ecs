package ecs

import (
	"fmt"
	"reflect"
)

var registeredComponents []reflect.Type

func RegisterComponent(component interface{}) {

	t := reflect.TypeOf(component)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for _, comp := range registeredComponents {
		if t.Name() == comp.Name() {
			panic(fmt.Sprintf("%s is already registered", comp.Name()))
		}
	}

	registeredComponents = append(registeredComponents, t)
}

func ComponentFromName(name string) (interface{}, error) {
	for _, comp := range registeredComponents {
		if comp.Name() == name {
			return reflect.New(comp).Interface(), nil
		}
	}

	return nil, fmt.Errorf("component '%s' was not found in the registry", name)
}
