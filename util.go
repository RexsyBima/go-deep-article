package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func createFolderIfNotExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create folder: %v", err)
		}
	}
	return nil
}

func getParseGetArg(args []string, expected string) (string, error) {
	for _, arg := range args {
		if strings.Contains(arg, expected) {
			return arg, nil
		}
	}
	return "", fmt.Errorf("Argument youtube not found, this should not be happening tho")
}

func extractYTIDFromURL(url string) (string, bool, error) {
	re := regexp.MustCompile(`(?:v=|youtu\.be/)([a-zA-Z0-9_-]{11})`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1], false, nil
	} else {
		return "", true, fmt.Errorf("fail to parse regex expression value")
	}
}
