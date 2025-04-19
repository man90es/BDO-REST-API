package utils

import (
	"testing"
	"time"
)

func TestFormatDateForHeaders(t *testing.T) {
	// Use a fixed date for deterministic output
	dt := time.Date(2025, 4, 19, 8, 6, 53, 0, time.FixedZone("CEST", 2*60*60))
	expected := dt.Format(time.RFC1123Z)
	got := FormatDateForHeaders(dt)
	if got != expected {
		t.Errorf("FormatDateForHeaders(%v) = %q; want %q", dt, got, expected)
	}

	// Check UTC
	dtUTC := time.Date(2025, 4, 19, 6, 6, 53, 0, time.UTC)
	expectedUTC := dtUTC.Format(time.RFC1123Z)
	gotUTC := FormatDateForHeaders(dtUTC)
	if gotUTC != expectedUTC {
		t.Errorf("FormatDateForHeaders(%v) = %q; want %q", dtUTC, gotUTC, expectedUTC)
	}
}
