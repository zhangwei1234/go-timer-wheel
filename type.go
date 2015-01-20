package timerWheel

import (
	"time"
	"sync"
)

//时钟轮
type TimerWheel struct {
	state			int 			//启动(1)or 停止(-1) 状态
	tickDuration	time.Duration   //卡槽每次跳动的时间间隔
	roundDuration	time.Duration	//一轮耗时
	wheelCount		int				//卡槽数
	wheel			[]*Iterator	    //卡槽
	tick			*time.Ticker    //时钟
	lock			sync.Mutex		//锁
	wheelCursor		int				//当前卡槽位置
	mask			int				//卡槽最大索引数
}

//到期执行的任务
type Task interface {
	//到期执行的函数
	Expire()
}

//时间轮卡槽迭代器
type Iterator struct {
	items 			map[string]*WheelTimeOut
}

//超时处理对象
type WheelTimeOut struct {
	delay			time.Duration //延迟时间
	index			int			  //卡槽索引位置
	rounds			int			  //需要转动的周期数
	task			Task		  //到期执行的任务
}