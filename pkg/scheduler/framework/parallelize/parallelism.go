/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parallelize

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

// DefaultParallelism is the default parallelism used in scheduler.
const DefaultParallelism int = 16

type DoWorkPieceFunc func(piece int)

// Parallelizer holds the parallelism for scheduler.
type Parallelizer struct {
	parallelism int
	pool        *ants.Pool
	isStoped    bool
	mutex       sync.Mutex
}

func init() {
	// It releases the default pool from ants.
	ants.Release()
}

// NewParallelizer returns an object holding the parallelism.
func NewParallelizer(stopCh <-chan struct{}, parallelism int) Parallelizer {
	if parallelism <= 0 {
		parallelism = DefaultParallelism
	}

	antsPool, _ := ants.NewPool(parallelism,
		ants.WithDisablePurge(false),
		ants.WithPreAlloc(true),
		ants.WithExpiryDuration(10*time.Minute),
	)

	pa := Parallelizer{
		parallelism: parallelism,
		pool:        antsPool,
	}

	go func() {
		<-stopCh
		pa.Stop()
	}()

	return pa
}

// Until is a wrapper around workqueue.ParallelizeUntil to use in scheduling algorithms.
func (pa Parallelizer) Until(ctx context.Context, pieces int, doWorkPiece DoWorkPieceFunc) {
	if pieces == 0 {
		return
	}

	toProcess := make(chan int, pieces)
	for i := 0; i < pieces; i++ {
		toProcess <- i
	}
	close(toProcess)

	var stop <-chan struct{}
	if ctx != nil {
		stop = ctx.Done()
	}

	wg := sync.WaitGroup{}
	wg.Add(pa.parallelism)
	taskProcess := func() {
		defer utilruntime.HandleCrash()
		defer wg.Done()
		for i := range toProcess {
			select {
			case <-stop:
				return
			default:
				doWorkPiece(i)
			}
		}
	}
	for i := 0; i < pa.parallelism; i++ {
		_ = pa.pool.Submit(taskProcess)
	}
	wg.Wait()
}

func (pa Parallelizer) Stop() {
	pa.mutex.Lock()
	defer pa.mutex.Unlock()
	if pa.isStoped {
		return
	}
	if pa.pool != nil {
		pa.pool.Release()
	}
	pa.isStoped = true
}

func ceilDiv(a, b int) int {
	return (a + b - 1) / b
}

// chunkSizeFor returns a chunk size for the given number of items to use for
// parallel work. The size aims to produce good CPU utilization.
// returns max(1, min(sqrt(n), n/Parallelism))
func chunkSizeFor(n, parallelism int) int {
	s := int(math.Sqrt(float64(n)))

	if r := n/parallelism + 1; s > r {
		s = r
	} else if s < 1 {
		s = 1
	}
	return s
}
