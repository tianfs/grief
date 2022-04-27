package errors

import (
    "testing"
)

const ERROR_NOT_FOUND string = "ERROR NOT FOUND"

func TestErrors(t *testing.T) {

    errors := getError()
    t.Logf("%+v", errors)
}
func getError() error {
    errors := New(123, string(ERROR_NOT_FOUND), "用户咋的了")
    return errors
}
