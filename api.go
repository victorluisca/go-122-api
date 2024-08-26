package main

import (
	"log"
	"net/http"
)

// representa o servidor da API e armazena o endereço em que o servidor sera executado
type APIServer struct {
	addr string
}

// cria e retorna uma instância do APIServer
func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

// método que inicia o servidor
func (s *APIServer) Run() error {
	// roteador HTTP
	router := http.NewServeMux()
	// define a rota users
	router.HandleFunc("GET /users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		w.Write([]byte("User ID: " + userID))
	})

	middlewareChain := MiddlewareChain(
		RequestLoggerMiddleware,
		RequireAuthMiddleware,
	)

	// configura o servidor
	server := http.Server{
		Addr:    s.addr,
		Handler: middlewareChain(router),
	}

	log.Printf("Server has started")

	// inicia o servidor
	return server.ListenAndServe()
}

// log das informações dos requests, método e url
func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method: %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

// verifica o token de autorização, token = 123
func RequireAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "123" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// alias para a assinatura de função dos middlewares
type Middleware func(http.Handler) http.HandlerFunc

// cria uma cadeia de middlewares e os aplica na ordem inversa
func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next.ServeHTTP
	}
}
