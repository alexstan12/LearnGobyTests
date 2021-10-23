package mocking

import (
	"fmt"
	"io"
	"time"
)

func Countdown(w io.Writer, sleeper Sleeper){
	for i:=3 ; i>0; i-- {
		sleeper.Sleep()
		fmt.Fprintln(w,i)
	}
	sleeper.Sleep()
	fmt.Fprint(w,"Go!")
}

type Sleeper interface {
	Sleep()
}

type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep(){
	s.Calls++
}

type SpyCountdownOperations struct {
	Calls []string
}

func (s *SpyCountdownOperations) Sleep(){
	s.Calls = append(s.Calls, sleep)
}

func (s *SpyCountdownOperations) Write(p []byte) (n int, err error){
	s.Calls = append(s.Calls, write)
	return
}

type ConfigurableSleeper struct {
	duration time.Duration
	sleep func(duration time.Duration)
}

func (c *ConfigurableSleeper) SetDuration(d time.Duration){
	c.duration = d
}
func (c *ConfigurableSleeper) SetSleepFunc(f func(d time.Duration)){
	c.sleep = f
}

func (c *ConfigurableSleeper) Sleep(){
	c.sleep(c.duration)
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration){
	s.durationSlept = duration
}

const write = "write"
const sleep = "sleep"
