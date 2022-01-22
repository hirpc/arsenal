package validation

import "strings"

type Validation interface {
	SecureString(s string) string
}

type validation struct {
	opt Options
}

func New(opts ...Option) *validation {
	var opt = Options{
		Level: High,
	}
	for _, o := range opts {
		o(&opt)
	}
	return &validation{
		opt: opt,
	}
}

func (v validation) SecureString(s string) string {
	for _, rule := range getFilterRules(v.opt.Level) {
		s = strings.ReplaceAll(s, rule, "")
	}
	return s
}

// getFilterRules TODO: to fixup rules
func getFilterRules(l level) []string {
	rules := []string{"'"}
	if l == High {
		rules = append(rules, "`")
	}
	return rules
}
