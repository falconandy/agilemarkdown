package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	timestampLayout = "2006-01-02 03:04 PM"
)

var (
	separators  = regexp.MustCompile(`[ \\/&_=+:]`)
	dashes      = regexp.MustCompile(`[\-]+`)
	illegalName = regexp.MustCompile(`[^[:alnum:]-]`)
	spacesRe    = regexp.MustCompile(`\s+`)
)

func PadIntLeft(value, width int) string {
	result := strconv.Itoa(value)
	if len(result) < width {
		result = strings.Repeat(" ", width-len(result)) + result
	}
	return result
}

func PadStringRight(value string, width int) string {
	result := value
	if len(result) < width {
		result += strings.Repeat(" ", width-len(result))
	}
	return result
}

func PadStringLeft(value string, width int) string {
	result := value
	if len(result) < width {
		result = strings.Repeat(" ", width-len(result)) + result
	}
	return result
}

func WeekStart(value time.Time) time.Time {
	weekday := value.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	weekStart := value.Add(time.Duration(-(weekday-1)*24) * time.Hour)
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
	return weekStart
}

func WeekDelta(baseValue, value time.Time) int {
	return int(WeekStart(value).Sub(WeekStart(baseValue)).Hours()) / 24 / 7
}

func AreEqualStrings(items1, items2 []string) bool {
	if len(items1) != len(items2) {
		return false
	}
	for i := range items1 {
		if items1[i] != items2[i] {
			return false
		}
	}
	return true
}

func WrapLinesToMarkdownCodeBlock(lines []string) []string {
	result := make([]string, 0, len(lines)+2)
	if len(lines) > 0 {
		result = append(result, "```")
		result = append(result, lines...)
		result = append(result, "```")
	}
	return result
}

func TitleFirstLetter(s string) string {
	first := true
	return strings.Map(
		func(r rune) rune {
			if first {
				first = false
				return unicode.ToTitle(r)
			}
			return r
		},
		s)
}

func ContainsStringIgnoreCase(items []string, item string) bool {
	item = strings.ToLower(item)
	for i := range items {
		if strings.ToLower(items[i]) == item {
			return true
		}
	}
	return false
}

func GetCurrentTimestamp() string {
	return GetTimestamp(time.Now())
}

func GetTimestamp(moment time.Time) string {
	return moment.Format(timestampLayout)
}

func ParseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(timestampLayout, timestamp)
}

func MakeMarkdownLink(linkTitle, linkPath, baseDir string) string {
	linkPath, _ = filepath.Abs(linkPath)
	baseDir, _ = filepath.Abs(baseDir)

	upCount := 0
	for !strings.HasPrefix(linkPath, baseDir+string(os.PathSeparator)) {
		upCount++
		baseDir = filepath.Dir(baseDir)
	}

	linkPath = strings.TrimPrefix(linkPath, baseDir)
	linkPath = strings.TrimPrefix(linkPath, string(os.PathSeparator))

	if upCount > 0 {
		linkPath = strings.Repeat(fmt.Sprintf("..%c", os.PathSeparator), upCount) + linkPath
	}

	return fmt.Sprintf("[%s](%s)", linkTitle, linkPath)
}

func JoinMarkdownLinks(links ...string) string {
	return strings.Join(links, " || ")
}

func GetValidFileName(name string) string {
	fileName := strings.TrimSpace(name)
	fileName = separators.ReplaceAllString(fileName, "-")
	fileName = illegalName.ReplaceAllString(fileName, "")
	fileName = dashes.ReplaceAllString(fileName, "-")
	return fileName
}

func CollapseWhiteSpaces(value string) string {
	return strings.TrimSpace(spacesRe.ReplaceAllString(value, " "))
}
