package http

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)


func Post(inputs interface{}, Uri string) ([]byte, error) {

	strPost := []byte("POST")

	var strRequestURI = []byte(Uri)
	creatorJSON, _ := json.Marshal(inputs)
	req := fasthttp.AcquireRequest()

	req.SetBody(creatorJSON)
	req.Header.SetContentType("application/json")
	req.Header.SetMethodBytes(strPost)
	req.SetRequestURIBytes(strRequestURI)
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		fmt.Println("handle http req error,%v", err)
		return []byte{}, err
	}
	fasthttp.ReleaseRequest(req)

	body := res.Body()

	// Only when you are done with body!
	fasthttp.ReleaseResponse(res)

	return body, nil
}
