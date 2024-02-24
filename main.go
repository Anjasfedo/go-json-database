package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	employees := []User{
		{"John", "23", "12121212", "Apel", Address{"bengkulu", "bengkulu", "indonesia", "1212"}},
		{"Alice", "25", "13131313", "Jeruk", Address{"jakarta", "jakarta", "indonesia", "1313"}},
		{"Michael", "30", "14141414", "Anggur", Address{"bandung", "jawa barat", "indonesia", "1414"}},
		{"Samantha", "28", "15151515", "Pisang", Address{"surabaya", "jawa timur", "indonesia", "1515"}},
		{"Emily", "26", "16161616", "Mangga", Address{"medan", "sumatera utara", "indonesia", "1616"}},
		{"Daniel", "32", "17171717", "Strawberry", Address{"makassar", "sulawesi selatan", "indonesia", "1717"}},
		{"Olivia", "27", "18181818", "Nanas", Address{"pontianak", "kalimantan barat", "indonesia", "1818"}},
	}
	
}
