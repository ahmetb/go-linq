package linq

import (
	"reflect"
	"testing"
)

func TestChannelToChannel(t *testing.T) {
	input := []int{30, 40, 50}

	inpCh := make(chan interface{})
	resCh := make(chan interface{})

	go func() {
		for _, i := range input {
			inpCh <- i
		}

		close(inpCh)
	}()

	go func() {
		FromChannel(inpCh).Where(func(i interface{}) bool {
			return i.(int) > 20
		}).ToChannel(resCh)
	}()

	result := []int{}
	for value := range resCh {
		result = append(result, value.(int))
	}

	if !reflect.DeepEqual(result, input) {
		t.Errorf("FromChannel().ToChannel()=%v expected %v", result, input)
	}
}
