package main

import (
	"encoding/json"
	"fmt"
	"github.com/jcelliott/lumber"
	"os"
	"sync"
)

const VERSION = "1.0.1"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New() () {

}

func Write() error {

}

func Read() error {

}

func ReadAll() () {

}

func Delete() error {

}

func getOrCreateMutex() *sync.Mutex {
	
}

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
		fmt.Println("Error ", err)
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

	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error ", err)
	}

	fmt.Println(records)

	allUsers := []User{}

	for _, value := range records {
		employeeFound := User{}

		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
			fmt.Println("Error ", err)
		}

		allUsers = append(allUsers, employeeFound)
	}

	fmt.Println((allUsers))

	// if err := db.Delete("users", "John"); err != nil {
	// 	fmt.Println("Error ", err)
	// }

	// if err := db.Delete("user", ""); err != nil {
	// 	fmt.Println("Error ", err)
	// }

}
