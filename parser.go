package httpreqparser

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/chonla/goline"
)

func Parse(req string) (*http.Request, error) {
	lines := goline.FromMultilineString(req)
	if len(lines) == 0 {
		return nil, errors.New("unexpected EOF")
	}

	headers := map[string]string{}
	bodyLines := []string{}
	body := ""
	var bodyBuffer *bytes.Buffer

	firstLine := lines[0]
	if captured, ok := firstLine.CaptureAll("^(CONNECT|HEAD|OPTIONS|POST|GET|PUT|PATCH|DELETE|TRACE) (.+) (.+)/(.+)$"); ok {
		method := captured[1]
		targetUrl := captured[2]

		collectingHeader := true
		for _, line := range lines[1:] {
			if collectingHeader {
				if headerCaptured, ok := line.CaptureAll("^([^:]+):(.+)$"); ok {
					headers[goline.Line(headerCaptured[1]).Lower().Trim().Value()] = goline.Line(headerCaptured[2]).Trim().Value()
				} else {
					if line.Value() == "" {
						collectingHeader = false
						continue
					}
				}
			} else {
				bodyLines = append(bodyLines, line.Value())
			}
		}

		if len(bodyLines) > 0 {
			body = strings.Join(bodyLines, "\n")
			bodyBuffer = bytes.NewBuffer([]byte(body))
		}

		pathInfo, err := url.Parse(targetUrl)
		if err != nil {
			return nil, err
		}
		// targetUrl has only path, rebuild url from collected data
		if pathInfo.Host == "" {
			if host, ok := headers["host"]; ok {
				targetUrl = fmt.Sprintf("http://%s%s", host, pathInfo.Path)
			}
		}

		var r *http.Request
		if bodyBuffer == nil {
			r, err = http.NewRequest(method, targetUrl, nil)
		} else {
			r, err = http.NewRequest(method, targetUrl, bodyBuffer)
		}
		if err != nil {
			return nil, err
		}
		for headerKey, headerValue := range headers {
			r.Header.Set(headerKey, headerValue)
		}
		return r, nil
	}
	return nil, errors.New("unexpected http request")
}
