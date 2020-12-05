package gox

import "context"

// A RequestChan is a channel to send Requests to another goroutine.
type RequestChan chan<- ResponseChan

// NewRequestChan returns a RequestChan and a channel to receive the requests
// sent to it.
// nolint unnamedResult
func NewRequestChan() (RequestChan, <-chan ResponseChan) {
	ch := make(chan ResponseChan)
	return ch, ch
}

// Send is a convenience method to send a request through the RequestChan. It
// returns the channel to receive the response.
func (req RequestChan) Send(ctx context.Context) (<-chan interface{}, error) {
	sendResp, recvResp := NewResponseChan()
	select {
	case req <- sendResp:
		return recvResp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// WaitResponse is a convenience method to send a request through the
// RequestChan, then wait for and return the response.
func (req RequestChan) WaitResponse(ctx context.Context) (interface{}, error) {
	resp, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}
	select {
	case v := <-resp:
		return v, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// A ResponseChan is sent as a request to a RequestChan to get a response
// through it from the other side.
type ResponseChan chan<- interface{}

// NewResponseChan returns a ResponseChan to be sent as a request through a
// RequestChan, and a channel to receive the response sent to it.
// nolint unnamedResult
func NewResponseChan() (ResponseChan, <-chan interface{}) {
	ch := make(chan interface{})
	return ch, ch
}
