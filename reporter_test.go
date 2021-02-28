package metrics

import (
	"fmt"
	"github.com/maxim-kuderko/metrics/drivers"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

func BenchmarkReporter_Send(b *testing.B) {
	b.ReportAllocs()
	r := NewReporter(WithDriver(drivers.NewNoop()))
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
	r := NewReporter(WithDriver(drivers.NewNoop()))
	name := `name`
	v := 0.1
	b.ResetTimer()
	concurrency := 8
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	arr := randArr()
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/concurrency; i++ {
				r.Send(name, v, arr...)
			}
		}()
	}
	wg.Wait()
}

func TestReporter_Send(t *testing.T) {
	stu := drivers.NewTestStub()
	r := NewReporter(WithDriver(stu))
	count := 100000
	tagsAr := map[string][]string{}
	for i := 0; i < count; i++ {
		tagsA := randArr()
		k := ``
		for _, v := range tagsA {
			k += v
		}
		tagsAr[k] = tagsA
	}
	for _, v := range tagsAr {
		r.Send(`name`, 1.0, v...)
	}

	r.Close()
	c := int64(0)
	i := 0
	if len(tagsAr) != len(stu.Metrics()) {
		t.Fatalf(`bad aggregation expecting %v, got %v`, len(tagsAr), len(stu.Metrics()))
	}
	for _, m := range stu.Metrics() {
		c += m.Values.Count
		if m.Values.Count == 0 {
			fmt.Print(0)
		}
		k := ``
		for _, v := range m.Tags {
			k += v
		}
		if !reflect.DeepEqual(tagsAr[k], m.Tags) {
			t.Fatalf(`expecting %v, got %v`, tagsAr[k], m.Tags)
		}
		i++
	}

	if int(c) != count {
		t.Fatalf(`expecting %v, got %v`, count, c)
	}
}

func TestReporter_Send_Small(t *testing.T) {
	stu := drivers.NewTestStub()
	r := NewReporter(WithDriver(stu))
	count := 200
	tagsAr := map[string][]string{}
	for i := 0; i < count; i++ {
		tagsA := randArr()
		r.Send(`name`, 1.0)
		k := ``
		for _, v := range tagsA {
			k += v
		}
		tagsAr[k] = tagsA
	}
	r.Close()
	c := int64(0)
	i := 0
	for _, m := range stu.Metrics() {
		c += m.Values.Count
		k := ``
		for _, v := range m.Tags {
			k += v
		}
		if !reflect.DeepEqual(tagsAr[k], m.Tags) {
			t.Fatalf(`expecting %v, got %v`, tagsAr[k], m.Tags)
		}
		i++
	}

	if int(c) != count {
		t.Fatalf(`expecting %v, got %v`, count, c)
	}
}

func TestReporter_SendC(t *testing.T) {
	stu := drivers.NewTestStub()
	r := NewReporter(WithDriver(stu))
	concurrency := 8
	count := 10000000 * concurrency
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	arr := randArr()
	name := `name`
	v := 1.0
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < count/concurrency; i++ {
				r.Send(name, v, arr...)
			}
		}()
	}
	wg.Wait()
	r.Close()
	c := int64(0)
	for _, m := range stu.Metrics() {
		c += m.Values.Count
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
