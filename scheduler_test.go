package timerWheel

import (
	"testing"
	"fmt"
	"time"
)

func TestScheculer(t *testing.T){
	
	wheel := NewTimerWheel()
	
	wheel.AddTask(&MyTask{}, 4 * time.Second)
	
	if len(wheel.wheel[40].items) ==0 {
		t.Fail()
	}
}

func TestScheculer1(t *testing.T){
	
	wheel := NewTimerWheel()
	
	wheel.AddTask(&MyTask{}, 1 * time.Minute)
	
	if len(wheel.wheel[88].items) ==0 {
		t.Fail()
	}
	
}

func TestNotifyTask4s(t *testing.T){
	wheel := NewTimerWheel()
	wheel.Start()
	wheel.AddTask(&MyTask{}, 4 * time.Second)
	wheel.AddTask(&MyTask{},	5 * time.Second)
	wheel.AddTask(&MyTask{}, 7 * time.Second)
	wheel.AddTask(&MyTask{}, 512 * 2 * 100 * time.Millisecond + 4 * time.Second)
	wheel.AddTask(&MyTask{}, 5 * time.Minute)
	time.Sleep(512 * 2 * 100 *40 * time.Millisecond)
}

type MyTask struct {
	
}

func(tk *MyTask) Expire() {
	fmt.Printf("-----------------> 执行 \n")	
}