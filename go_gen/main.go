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

	structType := reflect.TypeOf(Person{})

	fmt.Printf("func Map%s(dst %s) map[string]interface{} {\n", "Person", "Person")
	fmt.Printf("\toutput := map[string]interface{}{}\n\n")

	generateMappingCode(structType, "output", "")

	fmt.Printf("\treturn output\n")
	fmt.Printf("}\n")
}

func generateMappingCode(structType reflect.Type, outputVar, parent string) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldType := field.Type
		fieldName := field.Name

		if field.Type.Kind() == reflect.Struct {
			generateMappingCode(fieldType, outputVar, parent+fieldName+".")
		} else {
			fieldPath := fmt.Sprintf("%s%s", parent, fieldName)
			fmt.Printf("\tif !isZeroValue(dst.%s) {\n", fieldPath)
			fmt.Printf("\t\t%s[\"%s\"] = dst.%s\n", outputVar, fieldPath, fieldPath)
			fmt.Printf("\t}\n")
		}
	}
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0.0
	case reflect.String:
		return v.String() == ""
	default:
		return v.IsNil()
	}
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

// Testing

func MapPersonTest(dst Person) map[string]interface{} {
	output := map[string]interface{}{}

	if !isZeroValue(dst.Name) {
		output["Name"] = dst.Name
	}
	if !isZeroValue(dst.Age) {
		output["Age"] = dst.Age
	}
	if !isZeroValue(dst.Address.Street) {
		output["Address.Street"] = dst.Address.Street
	}
	if !isZeroValue(dst.Address.City) {
		output["Address.City"] = dst.Address.City
	}
	if !isZeroValue(dst.Address.Country) {
		output["Address.Country"] = dst.Address.Country
	}
	return output
}
