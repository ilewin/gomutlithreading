package main

import "time"

func MoveTrain(tr *Train, distance int, crossings []*Crossing) {
	for tr.Front < distance {
		tr.Front += 1
		for _, cr := range crossings {
			if tr.Front == cr.Position {
				cr.Intersection.Mutex.Lock()
				cr.Intersection.LockedBy = tr.Id
			}

			back := tr.Front - tr.Trainlength
			if back == cr.Position {
				cr.Intersection.LockedBy = -1
				cr.Intersection.Mutex.Unlock()
			}
		}
	}

	time.Sleep(30 * time.Millisecond)
}
