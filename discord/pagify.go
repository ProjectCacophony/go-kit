package discord

import (
	"strings"
)

// Pagify splits the given text into pages
func Pagify(text string) []string {
	var pages []string
	for _, page := range pagify(text, "\n") {
		if len(page) <= 1992 {
			if len(page) > 0 {
				pages = append(pages, page)
			}
		} else {
			for _, page := range pagify(page, ",") {
				if len(page) <= 1992 {
					if len(page) > 0 {
						pages = append(pages, page)
					}
				} else {
					for _, page := range pagify(page, "-") {
						if len(page) <= 1992 {
							if len(page) > 0 {
								pages = append(pages, page)
							}
						} else {
							for _, page := range pagify(page, " ") {
								if len(page) <= 1992 {
									if len(page) > 0 {
										pages = append(pages, page)
									}
								} else {
									pages = append(pages, page[0:1992])
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
	var currentOutputPart string
	var result []string
	textParts := strings.Split(text, delimiter)

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

	if len(result) == 0 || result[0] == "" {
		return []string{text}
	}

	return result
}
