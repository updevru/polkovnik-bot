package domain

import (
	"fmt"
	"math"
	"strings"
)

type Time struct {
	seconds int
}

func NewTime(seconds int) Time {
	return Time{seconds: seconds}
}

func (t Time) Seconds() int {
	return t.seconds
}

func (t Time) ToHumanFormat() string {

	result := make([]string, 0)
	val := float64(t.seconds)
	weeks := math.Floor(val / (60 * 60 * 24 * 7))
	if weeks > 0 {
		val -= 60 * 60 * 24 * 7 * weeks
		result = append(result, fmt.Sprintf("%dw", int(weeks)))
	}

	days := math.Floor(val / (60 * 60 * 24))
	if days > 0 {
		val -= days * 60 * 60 * 24
		result = append(result, fmt.Sprintf("%dd", int(days)))
	}

	hours := math.Floor(val / (60 * 60))
	if hours > 0 {
		val -= hours * 60 * 60
		result = append(result, fmt.Sprintf("%dh", int(hours)))
	}

	minutes := math.Floor(val / 60)
	if minutes > 0 {
		result = append(result, fmt.Sprintf("%dm", int(minutes)))
	}

	return strings.Join(result, " ")
}
