package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mrrahbarnia/GameApp/repository/postgresql"
	"github.com/mrrahbarnia/GameApp/service/userservice"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", healthCheck)
	mux.HandleFunc("/users/register", registerUser)

	server := http.Server{Addr: ":8090", Handler: mux}
	log.Fatal(server.ListenAndServe())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	// CURL http://localhost:8090/health-check
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, `{"error":"Invalid method"}`)
	}
	fmt.Fprint(w, `{"message":"Everything works fine"}`)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	// curl -X POST "http://localhost:8090/users/register" \
	// -H "Content-Type: application/json" \
	// -d '{"Name": "testUser", "PhoneNumber": "09131234567"}'
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, `{"error":"Invalid method"}`)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(
			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
		)

		return
	}

	var uReq userservice.RegisterRequest
	if err = json.Unmarshal(data, &uReq); err != nil {
		w.Write(
			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
		)

		return
	}

	pgRepo := postgresql.New()
	uSvc := userservice.New(pgRepo)
	_, err = uSvc.Register(uReq)
	if err != nil {
		w.Write(
			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
		)

		return
	}

	w.Write([]byte(`{"message":"user created"}`))
}
