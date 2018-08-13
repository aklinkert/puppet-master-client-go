package puppet_master

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func unexpectedResponse(res *http.Response) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read body of failed response (%v): %v", res.Status, err)
	}
	return fmt.Errorf("unexpected response %v: %v", res.Status, string(b))
}

func unprocessableEntity(res *http.Response, errs map[string][]string) error {
	var errStrs []string
	for field, e := range errs {
		errStrs = append(errStrs, fmt.Sprintf("%s (%v)", field, strings.Join(e, ", ")))
	}

	return fmt.Errorf("failed to save job. The following fields are invalid: %v", strings.Join(errStrs, ", "))
}

func dumpRequest(req *http.Request) {
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Printf("failed to stringify request body for logging: %v", err)
	} else {
		log.Print("request begin ##############################################")
		log.Print(string(b))
		log.Print("request end   ##############################################")
	}
}

func dumpResponse(res *http.Response) {
	b, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Printf("failed to stringify response body for logging: %v", err)
	} else {
		log.Print("response begin ##############################################")
		log.Print(string(b))
		log.Print("response end   ##############################################")
	}
}
