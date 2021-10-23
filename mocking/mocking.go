package mocking

import (
	"fmt"
	"io"
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
