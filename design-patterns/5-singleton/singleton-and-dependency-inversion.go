package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// think of a module as a singleton
type Database interface {
	GetPopulation(name string) int
}

// you want to be able to do is you want to have your Singleton implement some interface
// and then you depend on the interface. And the reason why it's important is because then
// you can substitute the implementer of this interface. You can replace the Singleton with,
// let's say, some sort of test dummy that you can use in your tests

func GetTotalPopulationEx(db Database, cities []string) int {
	result := 0
	for _, city := range cities {
		result += db.GetPopulation(city)
	}
	return result
}

type DummyDatabase struct {
	dummyData map[string]int
}

func (d *DummyDatabase) GetPopulation(name string) int {
	if len(d.dummyData) == 0 {
		d.dummyData = map[string]int{"alpha": 1, "beta": 2, "gamma": 3}
	}
	return d.dummyData[name]
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

	names := []string{"alpha", "gamma"} // expect 4
	tp := GetTotalPopulationEx(&DummyDatabase{}, names)
	ok := tp == 4
	fmt.Println(ok)

}
