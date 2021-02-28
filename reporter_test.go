package metrics

import (
	"github.com/maxim-kuderko/metrics/drivers"
	"runtime"
	"sync"
	"testing"
)

func BenchmarkReporter_Send(b *testing.B) {
	b.ReportAllocs()
	r := NewReporter(WithDriver(drivers.NewNoop()), WithBuffer(100))
	name := `name`
	v := 0.1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Send(name, v)
	}
}

func BenchmarkReporter_Send_Concurrent(b *testing.B) {
	b.ReportAllocs()
	r := NewReporter(WithDriver(drivers.NewNoop()), WithBuffer(100))
	name := `name`
	v := 0.1
	b.ResetTimer()
	concurrency := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/concurrency; i++ {
				r.Send(name, v, `a`, `b`, `c`, `d`)
			}
		}()
	}
	wg.Wait()
}
