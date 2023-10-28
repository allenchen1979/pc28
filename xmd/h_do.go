package xmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func hDo(userBase UserBase, method string, url RequestURL, s interface{}, t interface{}) error {
	buf := new(bytes.Buffer)

	// JSON Encode
	if err := json.NewEncoder(buf).Encode(s); err != nil {
		return err
	}

	// Sync
	req, err := http.NewRequest(method, string(url), buf)
	if err != nil {
		return err
	}

	// Sync Header
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("User-Agent", userBase.agent)

	req.Header.Set("Origin", "http://manorapp.pceggs.com")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", buildCookie(userBase.cookie))

	// Response
	http.DefaultClient.Timeout = 3 * time.Second
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// JSON Decode
	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return err
	}

	return nil
}

func buildCookie(q string) string {
	qs := strings.Split(q, "; ")

	ns := make([]string, 0)
	for _, s := range qs {
		s0 := strings.Split(s, "=")
		k, v := s0[0], strings.Join(s0[1:], "=")

		nv := v
		if strings.Contains(k, "lpvt") {
			nv = strconv.Itoa(int(time.Now().Unix() - 1))
		}

		ns = append(ns, strings.Join([]string{k, nv}, "="))
	}

	return strings.Join(ns, "; ")
}
