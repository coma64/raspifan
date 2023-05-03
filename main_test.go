package main

import (
	"raspifan/client"
	"testing"
)

func TestCalcFanSpeed(t *testing.T) {
	if speed := calcFanSpeed(5 * client.OneDegree); speed != 0 {
		t.Fatalf("Fan should not be turned on at 5 degrees; Got speed: %d", speed)
	}

	if speed := calcFanSpeed(19 * client.OneDegree); speed != 0 {
		t.Fatalf("Fan should not be turned on at 19 degrees; Got speed: %d", speed)
	}

	if speed := calcFanSpeed(41 * client.OneDegree); speed != 100_000 {
		t.Fatalf("Fan should be at full power at 41 degrees; Got speed: %d", speed)
	}

	if speed := calcFanSpeed(30 * client.OneDegree); speed < 0 || speed > 100_000 {
		t.Fatalf("Fan should be in between 0%% and 100%% at 30 degree; Got spedd: %d", speed)
	}
}
