package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type singletonDatabase struct {
	capitals map[string]int
}

func (db *singletonDatabase) GetPopulation(name string) int {
	return db.capitals[name]
}

// both init and sync.Once are thread-safe but only sync.Once is lazy
// laziness basically means that you only construct the database. You only read it from
// a disk to memory whenever somebody asks for it. So laziness is not going to be
// guaranteed in the Init function, unfortunately, but it can be guaranteed
// using things that once inside our own function. var once sync.Once

var instance *singletonDatabase
var once sync.Once

func GetSingletonDatabase() *singletonDatabase {
	once.Do(func() {
		db := singletonDatabase{}
		caps, err := readData(".\\capitals.txt")
		if err == nil {
			db.capitals = caps
		}
		instance = &db
	})
	return instance
}

func readData(path string) (map[string]int, error) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	file, err := os.Open(exPath + path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	result := map[string]int{}

	for scanner.Scan() {
		k := scanner.Text()
		scanner.Scan()
		v, _ := strconv.Atoi(scanner.Text())
		result[k] = v
	}

	return result, nil

}

func main() {

	db := GetSingletonDatabase()
	pop := db.GetPopulation("Seoul")
	fmt.Println("Pop of Seoul = ", pop)
}
