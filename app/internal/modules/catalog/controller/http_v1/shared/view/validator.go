package view_shared

import (
	"regexp"
)

var (
	regexpArticle = regexp.MustCompile(`^\S+$`)
)

func ValidateArticle(value string) bool {
	return regexpArticle.MatchString(value)
}
