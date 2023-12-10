package main

import (
	"regexp"
	"testing"
)

func TestPrintCurrentTimeFormat(t *testing.T) {
	// Вызываем тестируемую функцию
	nowInfo, err := printCurrentTime()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверка вывода времени в формате RFC3339
	expectedTimeRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	if !expectedTimeRegex.MatchString(nowInfo.timePackage) {
		t.Errorf("Expected time in RFC3339 format, got: %s", nowInfo.timePackage)
	}

	if !expectedTimeRegex.MatchString(nowInfo.ntpPackage) {
		t.Errorf("Expected time in RFC3339 format, got: %s", nowInfo.ntpPackage)
	}
}
