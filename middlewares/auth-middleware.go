package middlewares

import (
    "context"
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "github.com/golang-jwt/jwt"
)

func AuthMiddleware(jwtSecret []byte) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            tokenString := r.Header.Get("Authorization")
            if tokenString == "" {
                http.Error(w, "Missing token", http.StatusUnauthorized)
                return
            }

            tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                }
                return jwtSecret, nil
            })

            if err != nil || !token.Valid {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                http.Error(w, "Invalid token claims", http.StatusUnauthorized)
                return
            }

            adminID, err := strconv.Atoi(fmt.Sprintf("%v", claims["adminID"]))
            if err != nil {
                http.Error(w, "Invalid admin ID", http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), "adminID", adminID)
            next.ServeHTTP(w, r.WithContext(ctx))
        }
    }
}