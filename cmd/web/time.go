package main

import (
	"fmt"
	"math"
	"time"
)

func humanDate(t time.Time) string {
	return t.Format("January 02, 2006")
}

func timeAgo(pastTime time.Time) string {
	difference := time.Since(pastTime)

	seconds := int(math.Round(difference.Seconds()))
	if seconds < 60 {
		return fmt.Sprintf("%d %s ago", seconds, pluralize("second", seconds))
	}

	minutes := int(math.Round(difference.Minutes()))
	if minutes < 60 {
		return fmt.Sprintf("%d %s ago", minutes, pluralize("minute", minutes))
	}

	hours := int(math.Round(difference.Hours()))
	if hours < 24 {
		return fmt.Sprintf("%d %s ago", hours, pluralize("hour", hours))
	}

	days := hours / 24
	if days < 7 {
		return fmt.Sprintf("%d %s ago", int(days), pluralize("day", days))
	}

	weeks := days / 7
	if weeks < 4 {
		return fmt.Sprintf("%d %s ago", weeks, pluralize("week", weeks))
	}

	months := days / 30
	if days < 365 {
		return fmt.Sprintf("%d %s ago", months, pluralize("month", months))
	}

	years := days / 365
	return fmt.Sprintf("%d %s ago", years, pluralize("year", years))
}
