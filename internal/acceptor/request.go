package acceptor

import (
	"fmt"
)

const metadataID = "_id"
const metadataURL = "_url"
const metadataAttempts = "_attempts"
const metadataResponse = "_response"

type Request struct {
	Data map[string]interface{} `json:"data"`
}

func (r Request) MetadataID() string {
	return fmt.Sprintf("%v", r.Data[metadataID])
}

func (r Request) MetadataURL() string {
	return fmt.Sprintf("%v", r.Data[metadataURL])
}

func (r Request) MetadataAttempts() string {
	return fmt.Sprintf("%v", r.Data[metadataAttempts])
}

func (r Request) MetadataResponse() string {
	return fmt.Sprintf("%v", r.Data[metadataResponse])
}

func (r Request) AddMetadataAttempts(attempts int) {
	r.Data[metadataAttempts] = attempts
}

func (r Request) AddMetadataResponse(res interface{}) {
	r.Data[metadataResponse] = res
}
