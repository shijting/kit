package kit

import (
	"fmt"
	"testing"
)

type Context struct {
}
type Filter func(ctx *Context)
type FilterBuilder func(Filter) Filter

func do(fbs ...FilterBuilder) {
	f := func(cxt *Context) {
		fmt.Println("业务核心")
	}
	for i := len(fbs) - 1; i >= 0; i-- {
		f = fbs[i](f)
	}
	f(&Context{})
}

func MetricsFilterBuilder(f Filter) Filter {
	return func(cxt *Context) {
		fmt.Println("调用MetricsFilterBuilder1")
		f(cxt)
		fmt.Println("调用MetricsFilterBuilder2")
	}
}

func LoggingFilterBuilder(f Filter) Filter {
	return func(cxt *Context) {
		fmt.Println("调用LoggingFilterBuilder1")
		f(cxt)
		fmt.Println("调用LoggingFilterBuilder2")
	}
}

func TestFilter(t *testing.T) {
	do(MetricsFilterBuilder, LoggingFilterBuilder)
}
