package linq

import (
	"reflect"
	"testing"
)

func TestChannelToChannel(t *testing.T) {
	input := []int{30, 40, 50}

	inpCh := make(chan int)
	resCh := make(chan int)

	go func() {
		for _, i := range input {
			inpCh <- i
		}

		close(inpCh)
	}()

	go func() {
		FromChannel(inpCh).Where(func(i int) bool {
			return i > 20
		}).ToChannel(resCh)
	}()

	result := []int{}
	for value := range resCh {
		result = append(result, value)
	}

	if !reflect.DeepEqual(result, input) {
		t.Errorf("FromChannel().ToChannel()=%v expected %v", result, input)
	}
}
