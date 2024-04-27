package controllers

import "sync"

var mu sync.Mutex
var rwmu sync.RWMutex

type goRoutineManager struct {
	goRoutineCnt chan bool
}

// goRoutineManager to limit amount of concurrent running goRoutine

func (g *goRoutineManager) Run(f func()) {
	select {
	case g.goRoutineCnt <- true:
		wg.Add(1)
		go func() {
			f()
			<-g.goRoutineCnt
			wg.Done()
		}()
	default:
		f()
	}
}

// func (g *goRoutineManager) Run(f func()) {
//     g.goRoutineCnt <- true
//     wg.Add(1)
//     go func() {
        
//         f()
//         <-g.goRoutineCnt
// 		wg.Done()
//     }()
// }

func NewGoRoutineManager(goRoutineLimit int) *goRoutineManager {
	return &goRoutineManager{
		goRoutineCnt: make(chan bool, goRoutineLimit),
	}
}
