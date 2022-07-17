package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) moveOne() {
	accel := b.calcAcceleration()
	rWlock.Lock()
	b.velocity = b.velocity.Add(accel).Limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	rWlock.Unlock()
}

func (b *Boid) calcAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)

	avgVelocity, avgPosition, separation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}

	count := 0.0
	rWlock.RLock()
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
					avgPosition = avgPosition.Add(boids[otherBoidId].position)
					separation = separation.Add(b.position.Subtract(boids[otherBoidId].position).DivisionV(dist))
				}
			}
		}
	}
	rWlock.RUnlock()

	accel := Vector2D{
		b.borderBounce(b.position.x, screenWidth),
		b.borderBounce(b.position.y, screenHeight),
	}

	if count > 0 {
		avgVelocity, avgPosition = avgVelocity.DivisionV(count), avgPosition.DivisionV(count)
		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyV(adjRate)
		accelSeparation := separation.MultiplyV(adjRate)
		accel = accel.Add(accelAlignment).Add(accelCohesion).Add(accelSeparation)
	}

	return accel
}

func (b *Boid) borderBounce(p float64, max float64) float64 {
	if p < viewRadius {
		return 1 / p
	} else if p > max-viewRadius {
		return 1 / (p - max)
	}
	return 0
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(1 * time.Millisecond)
	}
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{
			rand.Float64() * screenWidth,
			rand.Float64() * screenHeight,
		},
		velocity: Vector2D{
			(rand.Float64() * 2) - 1.0,
			(rand.Float64() * 2) - 1.0,
		},

		id: bid,
	}

	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	go b.start()

}
