package sibling

import "github.com/ahmad-khatib0/go/idiomatic-approach-book/modules-packages-imports/foo/internal"

func AlsoUseDoubler(i int) int {
	return internal.Doubler(i)
}
