package metrics

import (
	"testing"
)

func BenchmarkReporter_Map(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range map[string]string{`test`: `test`} {

		}
	}
}
