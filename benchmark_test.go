package devicedetector

import (
	"testing"
)

// BenchmarkDeviceType measures DeviceType() performance with precompiled regexes
func BenchmarkDeviceType(b *testing.B) {
	cache, err := NewEmbeddedCache()
	if err != nil {
		b.Fatalf("Failed to load cache: %v", err)
	}

	testCases := []string{
		// Android Chrome - triggers Chrome/Safari regex checks
		"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
		// Android tablet
		"Mozilla/5.0 (Linux; Android 14; SM-X910) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		// Desktop Windows
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		// Android TV
		"Mozilla/5.0 (Linux; Android 9; BRAVIA 4K UR3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36",
		// Touch-enabled Windows
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Touch) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ua := testCases[i%len(testCases)]
		detector := New(cache, ua)
		_ = detector.DeviceType()
	}
}

// BenchmarkFullDetection measures complete detection including all fields
func BenchmarkFullDetection(b *testing.B) {
	cache, err := NewEmbeddedCache()
	if err != nil {
		b.Fatalf("Failed to load cache: %v", err)
	}

	ua := "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector := New(cache, ua)
		_ = detector.IsBot()
		_ = detector.OSName()
		_ = detector.DeviceType()
		_ = detector.DeviceBrand()
		_ = detector.DeviceName()
		_ = detector.Name()
		_ = detector.FullVersion()
	}
}
