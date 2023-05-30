package main

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/maps"
)

type Person struct {
	Name    string
	Age     int
	Address Address
}

type Address struct {
	Street  string
	City    string
	Country string
}

func main() {
	p := Person{
		Name: "John Doe",
		Age:  30,
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
	}

	fetchFields(p, "")

	outputGen := MapPerson(p)
	outputRef := MapPersonReflection(p, "")

	fmt.Println(outputGen)
	fmt.Println(outputRef)
}

type Field struct {
	Actual string
}

func fetchFields(s interface{}, parent string) map[string]string {
	output := map[string]string{}
	name := strings.ToLower(reflect.TypeOf(s).Name())
	v := reflect.ValueOf(s)
	t := v.Type()

	if parent == "" {
		parent = name
	} else {
		parent += fmt.Sprintf(".%s", name)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		fmt.Printf("%s.%s\n", parent, strings.ToLower(fieldType.Name))

		if field.Kind() == reflect.Struct {
			tmp := fetchFields(field.Interface(), parent)
			maps.Copy(output, tmp)
		} else {

		}
	}

	return output
}

const (
	MapFn = `
		func Map%s(dst %s) map[string]interface{} {
			output := map[string]interface{}

			%s

			return output
		}
	`
)

// Example code of what the generated output should kind of look like

func MapRouter() {}

func MapPerson(dst Person) map[string]interface{} {
	output := map[string]interface{}{}

	output["person.name"] = dst.Name
	output["person.age"] = dst.Age
	output["person.address.street"] = dst.Address.Street
	output["person.address.city"] = dst.Address.City
	output["person.address.country"] = dst.Address.Country

	return output
}

func MapPersonReflection(s interface{}, parent string) map[string]interface{} {
	output := map[string]interface{}{}

	name := strings.ToLower(reflect.TypeOf(s).Name())
	v := reflect.ValueOf(s)
	t := v.Type()

	if parent == "" {
		parent = name
	} else {
		parent += fmt.Sprintf(".%s", name)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			tmp := MapPersonReflection(field.Interface(), parent)
			maps.Copy(output, tmp)
		} else {
			output[fmt.Sprintf("%s.%s", parent, strings.ToLower(fieldType.Name))] = field.Interface()
		}
	}

	return output
}
