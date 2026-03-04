package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	StaffName   string `json:"staff_name"`
	HospitalURL string `json:"hospital_url"`
	StaffID     string `json:"staff_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(appSecret string) gin.HandlerFunc {
	jwtSecret := []byte(appSecret)

	return func(c *gin.Context) {
		tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimSpace(tokenString[7:])
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok {
			c.Set("staff_name", claims.StaffName)
			c.Set("hospital_url", claims.HospitalURL)
			c.Set("staff_id", claims.StaffID)
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
	}
}
