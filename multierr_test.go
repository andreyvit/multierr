package multierr_test

import (
	"errors"
	"fmt"

	"github.com/andreyvit/multierr"
)

var (
	oops     = errors.New("oops")
	whoops   = errors.New("whoops")
	whoopsie = errors.New("whoopsie")
)

func Example_nilError() {
	var err error
	fmt.Println(len(multierr.All(err)))
	// Output: 0
}

func Example_singleError() {
	var err error
	err = multierr.Append(err, oops)

	fmt.Println(err)
	fmt.Println(err == oops)
	fmt.Println(len(multierr.All(err)))
	// Output: oops
	// true
	// 1
}

func Example_multipleErrors() {
	var err error
	err = multierr.Append(err, oops)
	err = multierr.Append(err, whoops)
	err = multierr.Append(err, whoopsie)

	fmt.Println(err)
	// Output: 3 errors occurred:
	// (1) oops
	// (2) whoops
	// (3) whoopsie
}

func Example_is() {
	var err error
	err = multierr.Append(err, oops)
	err = multierr.Append(err, whoopsie)

	fmt.Println(errors.Is(err, oops))
	fmt.Println(errors.Is(err, whoops))
	fmt.Println(errors.Is(err, whoopsie))

	// Output: true
	// false
	// true
}
