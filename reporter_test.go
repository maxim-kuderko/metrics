package metrics

import (
	"github.com/maxim-kuderko/metrics/drivers"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
)

func BenchmarkReporter_Send(b *testing.B) {
	b.ReportAllocs()
	r := NewReporter(WithDriver(drivers.NewNoop()), WithBuffer(100))
	name := `name`
	v := 1.0
	b.ResetTimer()
	tArr := randArr()
	for i := 0; i < b.N; i++ {
		r.Send(name, v, tArr...)
	}
}

func randArr() []string {
	c := int(rand.Int31n(10)) + 1
	output := make([]string, 0, c)
	for c > 0 {
		c--
		output = append(output, RandStringRunes(5))

	}
	return output
}

func BenchmarkReporter_Send_Concurrent(b *testing.B) {
	b.ReportAllocs()
	r := NewReporter(WithDriver(drivers.NewNoop()), WithBuffer(500))
	name := `name`
	v := 0.1
	b.ResetTimer()
	concurrency := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			arr := randArr()
			for i := 0; i < b.N/concurrency; i++ {
				r.Send(name, v, arr...)
			}
		}()
	}
	wg.Wait()
}

func TestReporter_Send(t *testing.T) {
	stu := drivers.NewTestStub()
	r := NewReporter(WithDriver(stu), WithBuffer(500))
	count := 1000
	tagsAr := make([][]string, 0, count)
	for i := 0; i < count; i++ {
		tagsA := randArr()
		r.Send(`name`, 1.0, tagsA...)
		tagsAr = append(tagsAr, tagsA)
	}
	r.Close()
	c := 0.0
	i := 0
	for _, m := range stu.Metrics() {
		c += m.Value
		if !reflect.DeepEqual(tagsAr[i], m.Tags) {
			t.Fatalf(`expecting %v, got %v`, tagsAr[i], m.Tags)
		}
		i++
	}

	if int(c) != count {
		t.Fatalf(`expecting %v, got %v`, count, c)
	}
}

func TestReporter_Send_Small(t *testing.T) {
	stu := drivers.NewTestStub()
	r := NewReporter(WithDriver(stu), WithBuffer(1))
	count := 2
	tagsAr := make([][]string, 0, count)
	for i := 0; i < 2; i++ {
		tagsA := randArr()
		r.Send(`name`, 1.0, tagsA...)
		tagsAr = append(tagsAr, tagsA)
	}
	r.Close()
	c := 0.0
	i := 0
	for _, m := range stu.Metrics() {
		c += m.Value
		if !reflect.DeepEqual(tagsAr[i], m.Tags) {
			t.Fatalf(`expecting %v, got %v`, tagsAr[i], m.Tags)
		}
		i++
	}

	if int(c) != count {
		t.Fatalf(`expecting %v, got %v`, count, c)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
