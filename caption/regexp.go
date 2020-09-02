package caption

import (
	"regexp"
)

var flickr_re *regexp.Regexp
var orthis_re *regexp.Regexp

func init() {
	flickr_re = regexp.MustCompile(`o_\.\.*$`)
	orthis_re = regexp.MustCompile(`^(\d+)_[a-zA-Z0-9]+_o\.jpg$`)
}
