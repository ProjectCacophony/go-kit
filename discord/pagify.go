package discord

import (
	"strings"
)

// Pagify splits the given text into pages
func Pagify(text string) (pages []string) {
	for _, page := range pagify(text, "\n") {
		if len(page) <= 1992 {
			pages = append(pages, page)
		} else {
			for _, page := range pagify(page, ",") {
				if len(page) <= 1992 {
					pages = append(pages, page)
				} else {
					for _, page := range pagify(page, "-") {
						if len(page) <= 1992 {
							pages = append(pages, page)
						} else {
							for _, page := range pagify(page, " ") {
								if len(page) <= 1992 {
									pages = append(pages, page)
								} else {
									panic("unable to pagify text")
								}
							}
						}
					}
				}
			}
		}
	}
	return pages
}

func pagify(text string, delimiter string) []string {
	result := make([]string, 0)
	textParts := strings.Split(text, delimiter)
	currentOutputPart := ""
	for _, textPart := range textParts {
		if len(currentOutputPart)+len(textPart)+len(delimiter) <= 1992 {
			if len(currentOutputPart) > 0 || len(result) > 0 {
				currentOutputPart += delimiter + textPart
			} else {
				currentOutputPart += textPart
			}
		} else {
			result = append(result, currentOutputPart)
			currentOutputPart = ""
			if len(textPart) <= 1992 {
				currentOutputPart = textPart
			}
		}
	}
	if currentOutputPart != "" {
		result = append(result, currentOutputPart)
	}
	return result
}
