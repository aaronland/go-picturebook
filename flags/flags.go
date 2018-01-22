package flags

import (
	"fmt"
	"regexp"
	"strings"
)

type RegexpFlag []*regexp.Regexp

func (i *RegexpFlag) String() string {

	patterns := make([]string, 0)

	for _, re := range *i {
		patterns = append(patterns, fmt.Sprintf("%v", re))
	}

	return strings.Join(patterns, "\n")
}

func (i *RegexpFlag) Set(value string) error {

	re, err := regexp.Compile(value)

	if err != nil {
		return err
	}

	*i = append(*i, re)
	return nil
}

type PreProcessFlag []string

func (p *PreProcessFlag) String() string {
	return strings.Join(*p, "\n")
}

func (p *PreProcessFlag) Set(value string) error {
	*p = append(*p, value)
	return nil
}
