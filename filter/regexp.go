package filter

import (
	"context"
	"fmt"
	"gocloud.dev/blob"
	"net/url"
	"regexp"
)

func init() {

	ctx := context.Background()
	err := RegisterFilter(ctx, "regexp", NewRegexpFilter)

	if err != nil {
		panic(err)
	}
}

type RegexpFilter struct {
	Filter
	mode string
	re   *regexp.Regexp
}

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

func (f *RegexpFilter) Continue(ctx context.Context, bucket *blob.Bucket, path string) (bool, error) {

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
