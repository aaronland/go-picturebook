package filter

import (
	"context"
	"errors"
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
		return nil, err
	}

	var mode string

	switch u.Host {
	case "include":
		mode = "include"
	case "exclude":
		mode = "exclude"
	default:
		return nil, errors.New("Invalid mode")
	}

	q := u.Query()

	pat := q.Get("pattern")

	if pat == "" {
		return nil, errors.New("Missing pattern")
	}

	re, err := regexp.Compile(pat)

	if err != nil {
		return nil, err
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
