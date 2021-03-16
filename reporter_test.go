package metrics

import (
	"encoding/binary"
	"fmt"
	marshaler "github.com/golang/protobuf/proto"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/maxim-kuderko/metrics/drivers"
	"go.uber.org/atomic"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func BenchmarkReporter_Send(b *testing.B) {
	b.ReportAllocs()
	r := NewReporter(WithDriver(func() Driver {
		return drivers.NewNoop()
	}, 1, 1))
	name := `name`
	v := 0.1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Send(name, v)
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
	r := NewReporter(WithDriver(func() Driver {
		return drivers.NewNoop()
	}, 1, 8))
	name := `name`
	v := 0.1
	concurrency := 32
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	tagsAr := make([][]string, 0, 1000)
	for i := 0; i < 100; i++ {
		tagsAr = append(tagsAr, randArr())
	}
	b.ResetTimer()
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/concurrency; i++ {
				r.Send(name, v, tagsAr[i%len(tagsAr)]...)
			}
		}()
	}
	wg.Wait()
}

func TestReporter_Send(t *testing.T) {
	stu := drivers.NewTestStub()
	r := NewReporter(WithDriver(func() Driver {
		return stu
	}, 1, 1))
	count := 1000
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
	r := NewReporter(WithDriver(func() Driver {
		return stu
	}, 1, 1))
	count := 2000
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
	concurrency := 8
	r := NewReporter(WithDriver(func() Driver {
		return stu
	}, 1, 1))
	count := 100000 * concurrency
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	name := `name`
	v := 1.0
	tagsAr := make([][]string, 0, 100)
	for i := 0; i < 100; i++ {
		tagsA := randArr()
		k := ``
		for _, v := range tagsA {
			k += v
		}
		tagsAr = append(tagsAr, tagsA)
	}
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < count/concurrency; i++ {
				r.Send(name, v, tagsAr[i%len(tagsAr)]...)
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

func TestReporter_Send_UDP(t *testing.T) {
	count := 1000000
	addr := `127.0.0.1:9999`
	addrS := net.UDPAddr{
		Port: 9999,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ln, _ := net.ListenUDP(`udp`, &addrS)
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := atomic.NewInt64(0)

	go func() {
		defer wg.Done()
		for {
			buff := make([]byte, drivers.UDPBufferSize)
			ln.SetReadDeadline(time.Now().Add(time.Second * 1))
			n, err := ln.Read(buff)
			if err != nil {
				if 1-(float64(c.Load())/float64(count)) > 0.001 {
					t.Fatalf("got %d expexted %d, loss is %0.2f", c.Load(), count, 1-(float64(c.Load())/float64(count)))
				}
				break
			}
			go func(buff []byte) {
				scanned := 0
				for scanned+4 < n {
					size := int(binary.BigEndian.Uint32(buff[scanned : scanned+4]))
					scanned += 4
					tmp := proto.Metric{}
					if err = marshaler.Unmarshal(buff[scanned:scanned+size], &tmp); err != nil {
						t.Fatal(err)
					}
					c.Add(tmp.Values.Count)
					scanned += size
				}
			}(buff[:n])
		}
	}()
	go func() {
		defer wg.Done()
		cardinality := 1
		r := NewReporter(WithDriver(func() Driver {
			return drivers.NewUDP(addr)
		}, 1, 1))
		tagsAr := make([][]string, 0, cardinality)
		for i := 0; i < cardinality; i++ {
			tagsAr = append(tagsAr, randArr())
		}
		for i := 0; i < count; i++ {
			r.Send(`name`, 1, tagsAr[i%cardinality]...)
		}
		r.Close()
	}()
	wg.Wait()
}

func BenchmarkReporter_Send_UDP(b *testing.B) {
	b.ReportAllocs()
	addr := `127.0.0.1:9999`
	r := NewReporter(WithDriver(func() Driver {
		return drivers.NewUDP(addr)
	}, 1, 1))
	b.ResetTimer()
	name := `name`
	v := 1.0

	for i := 0; i < b.N; i++ {
		r.Send(name, v)
	}
	r.Close()
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
