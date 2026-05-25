package utils

import (
	"testing"
)

func TestGenSec(t *testing.T) {
	tests := []struct {
		name             string
		useragent        string
		expectedUA       string
		expectedMobile   string
		expectedPlatform string
	}{
		{
			name:             "Windows Chrome",
			useragent:        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
			expectedUA:       "\"Chromium\";v=\"131\", \"Google Chrome\";v=\"131\", \"Not/A)Brand\";v=\"99\"",
			expectedMobile:   "?0",
			expectedPlatform: "\"Windows\"",
		},
		{
			name:             "Windows Edge",
			useragent:        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0",
			expectedUA:       "\"Chromium\";v=\"131\", \"Microsoft Edge\";v=\"131\", \"Not/A)Brand\";v=\"99\"",
			expectedMobile:   "?0",
			expectedPlatform: "\"Windows\"",
		},
		{
			name:             "macOS Chrome",
			useragent:        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			expectedUA:       "\"Chromium\";v=\"130\", \"Google Chrome\";v=\"130\", \"Not/A)Brand\";v=\"99\"",
			expectedMobile:   "?0",
			expectedPlatform: "\"macOS\"",
		},
		{
			name:             "Android Mobile",
			useragent:        "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Mobile Safari/537.36",
			expectedUA:       "\"Chromium\";v=\"131\", \"Google Chrome\";v=\"131\", \"Not/A)Brand\";v=\"99\"",
			expectedMobile:   "?1",
			expectedPlatform: "\"Android\"",
		},
		{
			name:             "iPhone iOS",
			useragent:        "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
			expectedUA:       "\"Chromium\";v=\"131\", \"Google Chrome\";v=\"131\", \"Not/A)Brand\";v=\"99\"",
			expectedMobile:   "?1",
			expectedPlatform: "\"iOS\"",
		},
		{
			name:             "Linux Chrome",
			useragent:        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
			expectedUA:       "\"Chromium\";v=\"129\", \"Google Chrome\";v=\"129\", \"Not/A)Brand\";v=\"99\"",
			expectedMobile:   "?0",
			expectedPlatform: "\"Linux\"",
		},
		{
			name:             "Default Fallback",
			useragent:        "CustomScraper/1.0",
			expectedUA:       "\"Chromium\";v=\"131\", \"Google Chrome\";v=\"131\", \"Not/A)Brand\";v=\"99\"",
			expectedMobile:   "?0",
			expectedPlatform: "\"Windows\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ua, mobile, platform := GenSec(tt.useragent)
			if ua != tt.expectedUA {
				t.Errorf("%s: ua mismatch\nGot:  %s\nWant: %s", tt.name, ua, tt.expectedUA)
			}
			if mobile != tt.expectedMobile {
				t.Errorf("%s: mobile mismatch\nGot:  %s\nWant: %s", tt.name, mobile, tt.expectedMobile)
			}
			if platform != tt.expectedPlatform {
				t.Errorf("%s: platform mismatch\nGot:  %s\nWant: %s", tt.name, platform, tt.expectedPlatform)
			}
		})
	}
}
