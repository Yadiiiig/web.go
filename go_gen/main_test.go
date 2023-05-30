package main

import "testing"

func BenchmarkMapPersonGen(b *testing.B) {
	p := Person{
		Name: "John Doe",
		Age:  30,
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
	}

	for n := 0; n < b.N; n++ {
		_ = MapPerson(p)
	}
}

func BenchmarkMapPersonReflection(b *testing.B) {
	p := Person{
		Name: "John Doe",
		Age:  30,
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
	}

	for n := 0; n < b.N; n++ {
		_ = MapPersonReflection(p, "")
	}
}
