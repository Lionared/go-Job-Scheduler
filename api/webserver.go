package api

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type muxEntry struct {
	f       func(http.ResponseWriter, *http.Request) // entry func
	pattern string                                   // entry pattern
}

type myServeMux struct {
	mu    sync.RWMutex        // my serve mux locker
	m     map[string]muxEntry // my serve mux router rules
	hosts bool                // host info in any rule
}

func (mux *myServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("request from", r.RemoteAddr, ", path:", r.RequestURI)
	pathInfo, _ := url.Parse(r.RequestURI)

	if found, exist := mux.m[pathInfo.Path]; exist {
		found.f(w, r)
	} else {
		http.Error(w, "404 Not Found", 404)
	}
}

func (mux *myServeMux) HandleFunc(pattern string, f func(http.ResponseWriter, *http.Request)) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if f == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	mux.m[pattern] = muxEntry{f: f, pattern: pattern}

	if pattern[0] != '/' {
		mux.hosts = true
	}
}

func (mux *myServeMux) Handle(pattern string, handler http.Handler) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

type Middleware func(handler http.Handler) http.Handler

type WebServer struct {
	Host string
	Port int
	http.Server
}

func NewWebServer(host string, port int, readTimeout int64, writeTimeout int64) *WebServer {
	server := &WebServer{
		Host: host,
		Port: port,
	}
	myServeMux := &myServeMux{}
	myServeMux.RegisterHandlers()
	server.Addr = server.Host + ":" + strconv.Itoa(server.Port)
	server.Handler = myServeMux
	server.ReadTimeout = time.Second * time.Duration(readTimeout)
	server.WriteTimeout = time.Second * time.Duration(writeTimeout)
	return server
}

func (server *WebServer) Start() {
	if l, err := net.Listen("tcp", server.Addr); err != nil {
		log.Fatal(err)
	} else {
		defer func() {
			_ = l.Close()
		}()
		log.Println("Server start at:", server.Addr)
		_ = server.Serve(l)
	}
}

// methodMiddleware 请求方式中间件
func methodMiddleware(m string) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}

// authMiddleware 认证中间件
func authMiddleware() Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := getHeaderAuthorization(r)
			if strings.TrimSpace(authToken) == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if !VerifyJwtToken(authToken) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}

// chain 中间件嵌套
func chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

func getHeaderAuthorization(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func (mux *myServeMux) RegisterHandlers() {
	mux.HandleFunc("/", handleIndex)
	mux.Handle("/api/jobs", chain(http.HandlerFunc(handleJobsList), methodMiddleware("GET")))
	mux.Handle("/api/job/add", chain(http.HandlerFunc(handleJobAdd), methodMiddleware("POST")))
}
