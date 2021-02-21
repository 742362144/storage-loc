package sche

import (
	"fmt"
	"runtime"
	"storage-loc/db"
	"storage-loc/util"
)

const (

	INITIALIZED = 1

	RUNNING = 2

	YIELDED = 3

	COMPLETED = 4

	STOPPED = 5

	WAITING = 6
)

type TaskState uint8

type TaskPriority uint8

type Task struct {

	state TaskState

	time  uint64

	db_time  uint64

	db db.DB

	ext util.Extension

	fname string

	id uint64
}

func (task *Task) run(ws <-chan uint8)  {
	state := Paused // Begin in the paused state.
	for {
		select {
		case state = <-ws:
			switch state {
			case Stopped:
				fmt.Printf("Worker %d: Stopped\n", task.id)
				return
			case Running:
				fmt.Printf("Worker %d: Running\n", task.id)
			case Paused:
				fmt.Printf("Worker %d: Paused\n", task.id)
			}

		default:
			// We use runtime.Gosched() to prevent a deadlock in this case.
			// It will not be needed of work is performed here which yields
			// to the scheduler.
			runtime.Gosched()

			if state == Paused {
				break
			}

			// Do actual work here.

		}
	}
}
