package errors

import (
    "errors"
    "fmt"
)

type Error struct {
    Code     int32             `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
    Reason   string            `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
    Message  string            `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
    Metadata map[string]string `protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func New(code int, reason string, message string) *Error {
    return &Error{
        Code:    int32(code),
        Message: message,
        Reason:  reason,
    }
}
func (e *Error) Error() string {
    return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", e.Code, e.Reason, e.Message, e.Metadata)
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
    if se := new(Error); errors.As(err, &se) {
        return se.Code == e.Code && se.Reason == e.Reason
    }
    return false
}

// 添加
func (e *Error) WithMetadata(md map[string]string) *Error {
    e.Metadata = md
    return e
}
