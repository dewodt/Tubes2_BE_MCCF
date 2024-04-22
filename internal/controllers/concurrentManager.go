package controllers

import "sync"

var wg sync.WaitGroup

type goRoutineManager struct {
	goRoutineCnt chan bool
}

// goRoutineManager to limit amount of concurrent running goRoutine

func (g *goRoutineManager) Run(f func()) {
	wg.Add(1)
	select {
	case g.goRoutineCnt <- true:
		go func() {
			f()
			<-g.goRoutineCnt
		}()
	default:
		f()
	}
	wg.Done()
}

func NewGoRoutineManager(goRoutineLimit int) *goRoutineManager {
	return &goRoutineManager{
		goRoutineCnt: make(chan bool, goRoutineLimit),
	}
}
