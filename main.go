package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
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

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	if _, err := os.Stat(dir); err != nil {
		opts.Logger.Debug("Using '%s' (database already exist)\n", dir)

		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s'...\n", dir)

	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("missing collection")
	}

	if resource == "" {
		return fmt.Errorf("missing Resource")
	}

	mutex := d.getOrCreateMutex(collection)

	mutex.Lock()

	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)

	finalPath := filepath.Join(dir, resource+".json")

	tempPath := finalPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := os.WriteFile(tempPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tempPath, finalPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("missing collection")
	}

	if resource == "" {
		return fmt.Errorf("missing Resource")
	}

	record := filepath.Join(d.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}

	b, err := os.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("missing collection")
	}

	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := os.ReadDir(dir)

	var records []string

	for _, file := range files {
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		records = append(records, string(b))
	}

	return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("file not found")

	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}

	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	m, ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}

	return
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

	for _, f := range records {
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
