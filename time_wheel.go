package timerWheel

import (
	"errors"
	"time"
)

const (
	default_tick_duration = 100 * time.Millisecond //默认的时间轮 间隔时间 1秒
	default_wheel_count   = 512                    //默认的卡槽数512个
)

//初始化时间轮对象
func NewTimerWheel() *TimerWheel {

	return &TimerWheel{
		tickDuration:  default_tick_duration,
		wheelCount:    default_wheel_count,
		wheel:         createWheel(),
		wheelCursor:   0,
		mask:          default_wheel_count - 1,
		roundDuration: default_tick_duration * default_wheel_count,
	}
}

//启动时间轮
func (t *TimerWheel) Start() {
	t.lock.Lock()
	t.tick = time.NewTicker(default_tick_duration)
	defer t.lock.Unlock()
	go func() {
		for {
			select {
			case <-t.tick.C:
				t.wheelCursor++
				if t.wheelCursor == default_wheel_count {
					t.wheelCursor = 0
				}
				//判断当前卡槽中是否有超时的task
				iterator := t.wheel[t.wheelCursor]
				tasks := t.fetchExpiredTimeouts(iterator)
				t.notifyExpiredTimeOut(tasks)
			}
		}
	}()
}

//停止时间轮
func (t *TimerWheel) Stop() {
	t.tick.Stop()
}

func createWheel() []*Iterator {
	arr := make([]*Iterator, default_wheel_count)

	for v := 0; v < default_wheel_count; v++ {
		arr[v] = &Iterator{items: make(map[string]*WheelTimeOut)}
	}
	return arr
}

//添加一个超时任务
func (t *TimerWheel) AddTask(task Task, delay time.Duration) (string, error) {

	if task == nil {
		return "", errors.New("task is empty")
	}
	if delay <= 0 {
		return "", errors.New("delay Must be greater than zero")
	}
	timeOut := &WheelTimeOut{
		delay: delay,
		task:  task,
	}

	tid, err := t.scheduleTimeOut(timeOut)

	return tid, err
}

func (t *TimerWheel) RemoveTask(taskId string) error {
	for _, it := range t.wheel {
		for k, _ := range it.items {
			if taskId == k {
				delete(it.items, k)
			}
		}
	}
	return nil
}

func (t *TimerWheel) scheduleTimeOut(timeOut *WheelTimeOut) (string, error) {
	if timeOut.delay < t.tickDuration {
		timeOut.delay = t.tickDuration
	}
	lastRoundDelay := timeOut.delay % t.roundDuration
	lastTickDelay := timeOut.delay % t.tickDuration

	//计算卡槽位置
	relativeIndex := lastRoundDelay / t.tickDuration
	if lastTickDelay != 0 {
		relativeIndex = relativeIndex + 1
	}
	//计算时间轮圈数
	remainingRounds := timeOut.delay / t.roundDuration
	if lastRoundDelay == 0 {
		remainingRounds = remainingRounds - 1
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	stopIndex := t.wheelCursor + int(relativeIndex)
	if stopIndex > default_wheel_count {
		timeOut.index = stopIndex - default_wheel_count
		timeOut.rounds = int(remainingRounds) + 1
	} else {
		timeOut.index = stopIndex
		timeOut.rounds = int(remainingRounds)
	}
	item := t.wheel[stopIndex]
	if item == nil {
		item = &Iterator{
			items: make(map[string]*WheelTimeOut),
		}
	}

	key, err := GetGuid()
	if err != nil {
		return "", err
	}
	item.items[key] = timeOut
	t.wheel[stopIndex] = item

	return key, nil
}

//判断当前卡槽中是否有超时任务,将超时task加入切片中
func (t *TimerWheel) fetchExpiredTimeouts(iterator *Iterator) []*WheelTimeOut {
	t.lock.Lock()
	defer t.lock.Unlock()

	task := []*WheelTimeOut{}

	for k, v := range iterator.items {
		if v.rounds <= 0 { //已经超时了
			task = append(task, v)
			delete(iterator.items, k)
		} else {
			v.rounds--
		}
	}

	return task
}

//执行超时任务
func (t *TimerWheel) notifyExpiredTimeOut(tasks []*WheelTimeOut) {

	for _, task := range tasks {
		go task.task.Expire()
	}
}
