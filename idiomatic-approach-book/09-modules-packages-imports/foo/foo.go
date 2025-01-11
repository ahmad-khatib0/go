package foo

import "github.com/ahmad-khatib0/go/idiomatic-approach-book/modules-packages-imports/foo/internal"

func UseDoubler(i int) int {
	return internal.Doubler(i)
}
