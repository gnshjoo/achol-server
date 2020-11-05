package main

import (
	"net/http"
	"strings"
)

type router struct {
	// handlers는 또 다른 맵을 값으로 사용한 2차원 맵이다.
	handlers map[string]map[string]HandlerFunc
}

func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {
	// http 메서드로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method]

	if !ok {
		// 등록된 맵이 없으면 새 맵 생성
		m = make(map[string]HandlerFunc)
		r.handlers[method] = m
	}
	// http 메서드로 등록된 맵에 URL 패턴과 핸들러 함수 등록
	m[pattern] = h
}

func match(pattern, path string) (bool, map[string]string) {
	// 패턴과 패스가 정확히 일치하면 바로 true 반환
	if pattern == path {
		return true, nil
	}

	// / 단위로 구문
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	if len(patterns) != len(paths) {
		return false, nil
	}

	params := make(map[string]string)

	// "/"로 구분된 패턴/패스의 각 문자열을 하나씩 비교
	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			params[patterns[i][1:]] = paths[i]
		default:
			return false, nil
		}
	}
	return true, params
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for pattern, handler := range r.handlers[req.Method] {
		if ok, params := match(pattern, req.URL.Path); ok {
			//Context 생성
			c := Context{
				Params:         make(map[string]interface{}),
				ResponseWriter: w,
				Request:        req,
			}
			for k, v := range params {
				c.Params[k] = v
			}
			//요청 url에 해당하는 handler 수행
			handler(&c)
			return
		}
	}
	http.NotFound(w, req)
}
