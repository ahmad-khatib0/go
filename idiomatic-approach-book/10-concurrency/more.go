package main

import (
	"errors"
	"sync"
	"time"
)

// Time Out Code

func timeLimit() (int, error) {

	var result int
	var err error
	done := make(chan struct{})
	go func() {
		result, err = doSomeWork()
		close(done)
	}()

	select {
	case <-done:
		return result, err
	case <-time.After(2 * time.Second):
		return 0, errors.New("Sorry, your request timed out")
	}
}

func doSomeWork() (int, error) {
	time.Sleep(time.Second * 1)
	return 2, nil
}

// *********************************   Using WaitGroups  *********************************
func waitGroups() {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		doSomeWork1()
	}()
	go func() {
		defer wg.Done()
		doSomeWork2()
	}()
	go func() {
		defer wg.Done()
		doSomeWork3()
	}()

	wg.Wait()
}

func doSomeWork1() {}
func doSomeWork2() {}
func doSomeWork3() {}

func processAndGather(in <-chan int, processor func(int) int, num int) []int {
	out := make(chan int, num)
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			for v := range in {
				//  The for-range channel loop exits when out is closed and the buffer is empty
				out <- processor(v)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	var result []int
	for v := range out {
		result = append(result, v)
	}
	return result
}

// ********************************* Running Code Exactly Once  *********************************
// sometimes you want to lazy load, or call some initialization code exactly once after program launch time. This
// is usually because the initialization is relatively slow and may not even be needed every time your program runs.
type SlowComplicatedParser interface {
	Parse(string) string
}

var parser SlowComplicatedParser
var once sync.Once

// Declaring a sync.Once instance inside a function is usually the wrong thing to do, as a new instance
// will be created on every function call and there will be no memory of previous invocations.

func Parse(dateToParse string) string {
	// If Parse is called more than once, once.Do will not execute the closure again.
	once.Do(func() {
		parser = initParser()
	})
	return parser.Parse("")
}

func initParser() SlowComplicatedParser {
	// do all sorts of setup and loading here
	var p SlowComplicatedParser
	p.Parse("")
	return p
}

// ********************************* When to use Mutexes instead of channels *********************************
func scoreboardManager(in <-chan func(map[string]int), done <-chan struct{}) {
	scoreboard := map[string]int{}
	for {
		select {
		case <-done:
			return
		case f := <-in:
			f(scoreboard)
		}
	}
}

type ChannelScoreboardManager chan func(map[string]int)

func NewChannelScoreboardManager() (ChannelScoreboardManager, func()) {
	ch := make(ChannelScoreboardManager)
	done := make(chan struct{})
	go scoreboardManager(ch, done)
	return ch, func() {
		close(done)
	}
}

func (csm ChannelScoreboardManager) Update(name string, val int) {
	csm <- func(m map[string]int) {
		m[name] = val
	}
}

func (csm ChannelScoreboardManager) Read(name string) (int, bool) {
	var out int
	var ok bool
	done := make(chan struct{})
	csm <- func(m map[string]int) {
		out, ok = m[name]
		close(done)
	}
	<-done
	return out, ok
}

// An implementations using mutexes:
type MutexScoreboardManager struct {
	l          sync.RWMutex
	scoreboard map[string]int
}

func NewMutexScoreboardManager() *MutexScoreboardManager {
	return &MutexScoreboardManager{
		scoreboard: map[string]int{},
	}
}
func (msm *MutexScoreboardManager) Update(name string, val int) {
	msm.l.Lock()
	defer msm.l.Unlock()
	msm.scoreboard[name] = val
}
func (msm *MutexScoreboardManager) Read(name string) (int, bool) {
	msm.l.RLock()
	defer msm.l.RUnlock()
	val, ok := msm.scoreboard[name]
	return val, ok
}
