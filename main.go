package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mrrahbarnia/GameApp/infrastructure/bcrypt"
	"github.com/mrrahbarnia/GameApp/infrastructure/postgresql"
	authservice "github.com/mrrahbarnia/GameApp/service/auth"
	userservice "github.com/mrrahbarnia/GameApp/service/users"
)

const (
	JwtSignKey                 = "jwt_secret"
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", healthCheck)
	mux.HandleFunc("/users/register", registerUser)
	mux.HandleFunc("/users/login", loginUser)
	mux.HandleFunc("/users/profile", profile)

	server := http.Server{Addr: ":8090", Handler: mux}
	log.Println("Server is listening on port :8090")
	log.Fatal(server.ListenAndServe())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	// CURL http://localhost:8090/health-check
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, `{"error":"Invalid method"}`)

		return
	}
	fmt.Fprint(w, `{"message":"Everything works fine"}`)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	// curl -X POST "http://localhost:8090/users/register" \
	// -H "Content-Type: application/json" \
	// -d '{"name": "testUser", "phone_number": "09131234567", "password": "12345678"}'
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, `{"error":"Invalid method"}`)

		return
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
	bcrypt := bcrypt.New()
	authSvc := authservice.New(JwtSignKey, AccessTokenSubject,
		RefreshTokenSubject, AccessTokenExpireDuration, RefreshTokenExpireDuration)
	uSvc := userservice.New(pgRepo, bcrypt, authSvc)

	_, err = uSvc.Register(uReq)
	if err != nil {
		w.Write(
			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
		)

		return
	}

	w.Write([]byte(`{"message":"user created"}`))
}

func loginUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		fmt.Fprintf(w, `{"error":"Invalid method"}`)

		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(
			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
		)

		return
	}

	var lReq userservice.LoginRequest
	if err := json.Unmarshal(data, &lReq); err != nil {
		w.Write(
			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
		)

		return
	}

	pgRepo := postgresql.New()
	bcrypt := bcrypt.New()
	authSvc := authservice.New(JwtSignKey, AccessTokenSubject,
		RefreshTokenSubject, AccessTokenExpireDuration, RefreshTokenExpireDuration)
	uSvc := userservice.New(pgRepo, bcrypt, authSvc)

	resp, err := uSvc.Login(lReq)
	if err != nil {
		w.Write(
			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
		)

		return
	}

	data, err = json.Marshal(resp)
	if err != nil {
		w.Write(
			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
		)

		return
	}

	w.Write(data)
}

func profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, `{"error":"Invalid method"}`)

		return
	}

	authSvc := authservice.New(
		JwtSignKey,
		AccessTokenSubject,
		RefreshTokenSubject,
		AccessTokenExpireDuration,
		RefreshTokenExpireDuration,
	)

	authToken := r.Header.Get("Authorization")
	claims, err := authSvc.ParseToken(authToken)
	if err != nil {
		w.Write(
			[]byte(fmt.Sprintf("Token is not valid")),
		)

		return
	}

	pgRepo := postgresql.New()
	bcrypt := bcrypt.New()

	userSvc := userservice.New(pgRepo, bcrypt, authSvc)
	resp, err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	w.Write(data)
}
