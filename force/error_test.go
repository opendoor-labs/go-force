package force

import (
	"fmt"
	"testing"
)

func TestWasNotFOund(t *testing.T) {
	apiErr := ApiErrors{
		&ApiError{ErrorCode: "NOT_FOUND"},
	}

	found, err := WasNotFound(apiErr)
	if err != nil {
		fmt.Println(err)
		t.Error("expected WasNotFound not to return an error")
	}
	if !found {
		t.Error("expected err to say it was not found")
	}

	apiErr = ApiErrors{
		&ApiError{ErrorCode: "SO_TOTALLY_FOUND"},
	}

	found, err = WasNotFound(apiErr)
	if err != nil {
		fmt.Println(err)
		t.Error("expected WasNotFound not to return an error")
	}
	if found {
		t.Error("expected err not to say it was not found")
	}
}
