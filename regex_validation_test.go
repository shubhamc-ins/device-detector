package devicedetector

import (
	"testing"
)

// TestRegexesLoad verifies that all embedded regex files load without errors
func TestRegexesLoad(t *testing.T) {
	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Fatalf("Failed to load embedded cache: %v", err)
	}

	if cache == nil {
		t.Fatal("Cache is nil after loading")
	}

	if cache.Bot == nil {
		t.Error("Bot cache is nil")
	}
	if cache.Client == nil {
		t.Error("Client cache is nil")
	}
	if cache.Device == nil {
		t.Error("Device cache is nil")
	}
	if cache.Hint == nil {
		t.Error("Hint cache is nil")
	}
	if cache.OS == nil {
		t.Error("OS cache is nil")
	}

	t.Log("All regex caches loaded successfully")
}

// TestCriticalUserAgents tests detection of critical user agents that must work
func TestCriticalUserAgents(t *testing.T) {
	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	testCases := []struct {
		name          string
		ua            string
		expectBot     bool
		expectBotName string   // for bots
		expectOS      string   // for devices (partial match ok)
		expectDevice  string   // device type
	}{
		// === BOTS ===
		{
			name:          "Googlebot",
			ua:            "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			expectBot:     true,
			expectBotName: "Googlebot",
		},
		{
			name:          "Bingbot",
			ua:            "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
			expectBot:     true,
			expectBotName: "BingBot",
		},
		{
			name:          "GPTBot",
			ua:            "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)",
			expectBot:     true,
			expectBotName: "GPTBot",
		},
		{
			name:          "ClaudeBot",
			ua:            "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; ClaudeBot/1.0; claudebot@anthropic.com)",
			expectBot:     true,
			expectBotName: "ClaudeBot",
		},

		// === MOBILE DEVICES ===
		{
			name:         "iPhone Safari",
			ua:           "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
			expectBot:    false,
			expectOS:     "iOS",
			expectDevice: "smartphone",
		},
		{
			name:         "Android Chrome",
			ua:           "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
			expectBot:    false,
			expectOS:     "Android",
			expectDevice: "smartphone",
		},
		{
			name:         "iPad Safari",
			ua:           "Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
			expectBot:    false,
			expectOS:     "iPadOS",
			expectDevice: "tablet",
		},

		// === DESKTOP ===
		{
			name:         "Chrome Windows",
			ua:           "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expectBot:    false,
			expectOS:     "Windows",
			expectDevice: "desktop",
		},
		{
			name:         "Safari macOS",
			ua:           "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
			expectBot:    false,
			expectOS:     "Mac",
			expectDevice: "desktop",
		},
		{
			name:         "Firefox Linux",
			ua:           "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
			expectBot:    false,
			expectOS:     "GNU/Linux",
			expectDevice: "desktop",
		},

		// === TV / STREAMING ===
		{
			name:         "Fire TV",
			ua:           "Mozilla/5.0 (Linux; Android 9; AFTSSS Build/PS7255) AppleWebKit/537.36 (KHTML, like Gecko) Silk/100.1.1 like Chrome/100.0.4896.127 Safari/537.36",
			expectBot:    false,
			expectOS:     "Fire OS",
			expectDevice: "tv",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test for panics
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic during detection: %v", r)
				}
			}()

			detector := New(cache, tc.ua)
			if detector == nil {
				t.Fatal("Detector is nil")
			}

			isBot := detector.IsBot()

			if tc.expectBot {
				if !isBot {
					t.Errorf("Expected bot, got device")
					return
				}
				botName := detector.BotName()
				if botName != tc.expectBotName {
					t.Errorf("Expected bot name '%s', got '%s'", tc.expectBotName, botName)
				}
			} else {
				if isBot {
					t.Errorf("Expected device, got bot: %s", detector.BotName())
					return
				}

				osName := detector.OSName()
				deviceType := detector.DeviceType()

				if tc.expectOS != "" && osName != tc.expectOS {
					t.Errorf("Expected OS '%s', got '%s'", tc.expectOS, osName)
				}
				if tc.expectDevice != "" && deviceType != tc.expectDevice {
					// Allow phablet as smartphone variant
					if !(tc.expectDevice == "smartphone" && deviceType == "phablet") {
						t.Errorf("Expected device type '%s', got '%s'", tc.expectDevice, deviceType)
					}
				}
			}
		})
	}
}

// TestNewDevicePatterns tests patterns added after July 2024
// These are critical for validating upstream sync
func TestNewDevicePatterns(t *testing.T) {
	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	// These patterns require the 2024-2025 Matomo updates
	newPatterns := []struct {
		name        string
		ua          string
		expectBot   bool
		description string
	}{
		// New AI Bots
		{"GPTBot", "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)", true, "OpenAI crawler"},
		{"ChatGPT-User", "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko); compatible; ChatGPT-User/1.0; +https://openai.com/bot", true, "OpenAI browsing agent"},
		{"ClaudeBot", "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; ClaudeBot/1.0; claudebot@anthropic.com)", true, "Anthropic crawler"},
		{"PerplexityBot", "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; PerplexityBot/1.0; +https://perplexity.ai/bot)", true, "Perplexity crawler"},
		{"Meta-ExternalAgent", "meta-externalagent/1.1 (+https://developers.facebook.com/docs/sharing/webmasters/crawler)", true, "Meta AI crawler"},
		{"Bytespider", "Mozilla/5.0 (Linux; Android 5.0) AppleWebKit/537.36 (KHTML, like Gecko) Mobile Safari/537.36 (compatible; Bytespider; spider-feedback@bytedance.com)", true, "ByteDance crawler"},

		// New Devices (2024-2025)
		{"iPhone 16 Pro", "Mozilla/5.0 (iPhone17,1; U; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/602.1", false, "Released Sep 2024"},
		{"Pixel 9 Pro", "Mozilla/5.0 (Linux; Android 14; Pixel 9 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Mobile Safari/537.36", false, "Released Aug 2024"},
		{"Galaxy S24 Ultra", "Mozilla/5.0 (Linux; Android 14; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Mobile Safari/537.36", false, "Released Jan 2024"},
	}

	for _, tc := range newPatterns {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic during detection of %s: %v", tc.name, r)
				}
			}()

			detector := New(cache, tc.ua)
			isBot := detector.IsBot()

			if tc.expectBot && !isBot {
				t.Errorf("%s: Expected bot detection for %s", tc.name, tc.description)
			}
			if !tc.expectBot && isBot {
				t.Errorf("%s: Expected device, got bot for %s", tc.name, tc.description)
			}

			// For devices, verify we get meaningful detection
			if !tc.expectBot && !isBot {
				osName := detector.OSName()
				deviceType := detector.DeviceType()
				if osName == "" || deviceType == "" {
					t.Errorf("%s: Missing OS ('%s') or device type ('%s') for %s", tc.name, osName, deviceType, tc.description)
				}
			}
		})
	}
}

// TestNoPanicsOnMalformedUA ensures the library doesn't panic on edge cases
func TestNoPanicsOnMalformedUA(t *testing.T) {
	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	edgeCases := []string{
		"",                              // empty
		" ",                             // whitespace
		"Mozilla/5.0",                   // minimal
		"Mozilla/5.0 ()",                // empty parens
		"some random string",            // garbage
		"!@#$%^&*()",                    // special chars
		string(make([]byte, 10000)),     // very long
		"Mozilla/5.0 \x00\x01\x02",      // binary
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X)", // truncated
	}

	for i, ua := range edgeCases {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic on edge case UA: %v", r)
				}
			}()

			detector := New(cache, ua)
			// Just ensure these don't panic
			_ = detector.IsBot()
			_ = detector.BotName()
			_ = detector.OSName()
			_ = detector.DeviceType()
			_ = detector.DeviceBrand()
			_ = detector.DeviceName()
			_ = detector.Name()
			_ = detector.FullVersion()
		})
	}
}
