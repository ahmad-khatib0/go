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

// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
func GetTotalPopulation(cities []string) int {
	result := 0
	for _, city := range cities {
		result += GetSingletonDatabase().GetPopulation(city) // DIP situation
	}
	return result
}

// the problem here is that you are on a real copy of the DB , and for example
// the population of Seoul and mexico city might ( or it'l change definitely  )
// so the test will break
func main() {
	cities := []string{"Seoul", "Mexico City"}

	tp := GetTotalPopulation(cities)
	ok := tp == (17500000 + 17400000)
	fmt.Println(ok)

}
