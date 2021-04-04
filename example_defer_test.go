package multierr_test

import (
	"errors"
	"fmt"

	"github.com/andreyvit/multierr"
)

func Example() {
	fmt.Println(sprinkleMagicDust())

	// Output: 2 errors occurred:
	// (1) magic not available
	// (2) close: whoopsie
}

func sprinkleMagicDust() (err error) {
	d := &Dummy{}
	defer func() {
		err = multierr.Append(err, d.Close())
	}()

	return d.DoSomeMagic()
}

type Dummy struct{}

func (d Dummy) DoSomeMagic() error {
	return errors.New("magic not available")
}
func (d Dummy) Close() error {
	return errors.New("close: whoopsie")
}
