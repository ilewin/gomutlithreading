package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	matrixSize = 250
)

var (
	matrixA = [matrixSize][matrixSize]int{}
	matrixB = [matrixSize][matrixSize]int{}
	result  = [matrixSize][matrixSize]int{}
	rwLock  = sync.RWMutex{}
	cond    = sync.NewCond(rwLock.RLocker())
	wg      = sync.WaitGroup{}
)

func generateRandomMatrix(matrix *[matrixSize][matrixSize]int) {
	for r := 0; r < matrixSize; r++ {
		for c := 0; c < matrixSize; c++ {
			matrix[r][c] += rand.Intn(10) - 5
		}
	}
}

func computeRow(row int) {
	rwLock.RLock()
	for {
		wg.Done()
		cond.Wait()
		for col := 0; col < matrixSize; col++ {
			for i := 0; i < matrixSize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
	}

}

func main() {
	start := time.Now()
	wg.Add(matrixSize)
	for i := 0; i < matrixSize; i++ {
		go computeRow(i)
	}

	for i := 0; i < 100; i++ {
		wg.Wait()
		rwLock.Lock()
		generateRandomMatrix(&matrixA)
		generateRandomMatrix(&matrixB)
		wg.Add(matrixSize)
		rwLock.Unlock()
		cond.Broadcast()
	}

	elapsed := time.Since(start)
	fmt.Printf("it took %v\n", elapsed)
}
