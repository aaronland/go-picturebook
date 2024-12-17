package filter

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/aaronland/go-picturebook/bucket"
)

func init() {

	ctx := context.Background()
	err := RegisterFilter(ctx, "regexp", NewRegexpFilter)

	if err != nil {
		panic(err)
	}
}

// type AnyFilter implements the `Filter` interface that determines whether an image should be included in a picturebook using a regular expression.
type RegexpFilter struct {
	Filter
	mode string
	re   *regexp.Regexp
}

// NewRegexpFilter returns a new instance of `RegExpFilter` for 'uri' which must be parsable as a valid `net/url` URL instance.
// That URI must contain a host value of either 'include' or 'exclude' and query parameter 'pattern' containing a valid regular
// expression used to test file paths for inclusion in a picturebook.
func NewRegexpFilter(ctx context.Context, uri string) (Filter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	var mode string

	switch u.Host {
	case "include":
		mode = "include"
	case "exclude":
		mode = "exclude"
	default:
		return nil, fmt.Errorf("Invalid mode '%s'", u.Host)
	}

	q := u.Query()

	pat := q.Get("pattern")

	if pat == "" {
		return nil, fmt.Errorf("Missing ?pattern= parameter")
	}

	re, err := regexp.Compile(pat)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile pattern, %w", err)
	}

	f := &RegexpFilter{
		mode: mode,
		re:   re,
	}

	return f, nil
}

// Continues returns a boolean value signaling whether or not 'path' should be included in a picturebook.
func (f *RegexpFilter) Continue(ctx context.Context, source_bucket bucket.Bucket, path string) (bool, error) {

	match := f.re.MatchString(path)

	switch f.mode {
	case "include":

		if !match {
			return false, nil
		}

	case "exclude":

		if match {
			return false, nil
		}

	default:
		// pass
	}

	return true, nil
}
