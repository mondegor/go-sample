package view_shared

import (
    "regexp"
)

var (
    regexpArticle = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.+-]*[a-zA-Z0-9]$`)
)

func ValidateArticle(value string) bool {
    return regexpArticle.MatchString(value)
}
