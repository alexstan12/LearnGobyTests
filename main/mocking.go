package main

import (
	"github.com/alexstan12/LearnGobyTests/mocking"
	"os"
	"time"
)

func main(){
	//sleeper := DefaultSleeper{}
	//mocking.Countdown(os.Stdout, &sleeper)
	sleeper := &mocking.ConfigurableSleeper{}
	sleeper.SetDuration(1*time.Second)
	sleeper.SetSleepFunc(time.Sleep)
	mocking.Countdown(os.Stdout, sleeper)
}

