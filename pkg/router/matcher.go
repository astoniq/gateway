package router

import "regexp"

const (
	matchRule = `(\/\*(.+)?)`
)

type ListenPathMatcher struct {
	req *regexp.Regexp
}

func NewListenPathMatcher() *ListenPathMatcher {
	return &ListenPathMatcher{regexp.MustCompile(matchRule)}
}

func (l *ListenPathMatcher) Match(listenPath string) bool {
	return l.req.MatchString(listenPath)
}

func (l *ListenPathMatcher) Extract(listenPath string) string {
	return l.req.ReplaceAllString(listenPath, "")
}
