package main

import (
	"time"

	"github.com/mrrahbarnia/GameApp/config"
	"github.com/mrrahbarnia/GameApp/delivery/httpserver"
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
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8090},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshSubject:        RefreshTokenSubject,
			AccessSubject:         AccessTokenSubject,
		},
		PostgreSQL: postgresql.Config{
			Username: "admin",
			Password: "123456",
			Host:     "localhost",
			Port:     5432,
			DBName:   "db",
		},
	}

	authSvc, userSvc := setupServices(cfg)

	server := httpserver.New(cfg, userSvc, authSvc)
	server.Serve()

}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	bcrypt := bcrypt.New()
	PGRepo := postgresql.New(cfg.PostgreSQL)
	userSvc := userservice.New(PGRepo, bcrypt, authSvc)

	return authSvc, userSvc
}

// func registerUser(w http.ResponseWriter, r *http.Request) {
// 	// curl -X POST "http://localhost:8090/users/register" \
// 	// -H "Content-Type: application/json" \
// 	// -d '{"name": "testUser", "phone_number": "09131234567", "password": "12345678"}'
// 	if r.Method != http.MethodPost {
// 		fmt.Fprintf(w, `{"error":"Invalid method"}`)

// 		return
// 	}

// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		w.Write(
// 			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
// 		)

// 		return
// 	}

// 	var uReq userservice.RegisterRequest
// 	if err = json.Unmarshal(data, &uReq); err != nil {
// 		w.Write(
// 			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
// 		)

// 		return
// 	}

// 	pgRepo := postgresql.New()
// 	bcrypt := bcrypt.New()
// 	authSvc := authservice.New(JwtSignKey, AccessTokenSubject,
// 		RefreshTokenSubject, AccessTokenExpireDuration, RefreshTokenExpireDuration)
// 	uSvc := userservice.New(pgRepo, bcrypt, authSvc)

// 	_, err = uSvc.Register(uReq)
// 	if err != nil {
// 		w.Write(
// 			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
// 		)

// 		return
// 	}

// 	w.Write([]byte(`{"message":"user created"}`))
// }

// func loginUser(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodPost {
// 		fmt.Fprintf(w, `{"error":"Invalid method"}`)

// 		return
// 	}

// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		w.Write(
// 			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
// 		)

// 		return
// 	}

// 	var lReq userservice.LoginRequest
// 	if err := json.Unmarshal(data, &lReq); err != nil {
// 		w.Write(
// 			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
// 		)

// 		return
// 	}

// 	pgRepo := postgresql.New()
// 	bcrypt := bcrypt.New()
// 	authSvc := authservice.New(JwtSignKey, AccessTokenSubject,
// 		RefreshTokenSubject, AccessTokenExpireDuration, RefreshTokenExpireDuration)
// 	uSvc := userservice.New(pgRepo, bcrypt, authSvc)

// 	resp, err := uSvc.Login(lReq)
// 	if err != nil {
// 		w.Write(
// 			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
// 		)

// 		return
// 	}

// 	data, err = json.Marshal(resp)
// 	if err != nil {
// 		w.Write(
// 			[]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())),
// 		)

// 		return
// 	}

// 	w.Write(data)
// }

// func profile(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		fmt.Fprintf(w, `{"error":"Invalid method"}`)

// 		return
// 	}

// 	authSvc := authservice.New(
// 		JwtSignKey,
// 		AccessTokenSubject,
// 		RefreshTokenSubject,
// 		AccessTokenExpireDuration,
// 		RefreshTokenExpireDuration,
// 	)

// 	authToken := r.Header.Get("Authorization")
// 	claims, err := authSvc.ParseToken(authToken)
// 	if err != nil {
// 		w.Write(
// 			[]byte(fmt.Sprintf("Token is not valid")),
// 		)

// 		return
// 	}

// 	pgRepo := postgresql.New()
// 	bcrypt := bcrypt.New()

// 	userSvc := userservice.New(pgRepo, bcrypt, authSvc)
// 	resp, err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
// 	if err != nil {
// 		w.Write([]byte(
// 			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
// 		))

// 		return
// 	}

// 	data, err := json.Marshal(resp)
// 	if err != nil {
// 		w.Write([]byte(
// 			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
// 		))

// 		return
// 	}

// 	w.Write(data)
// }
