package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"recipe_api/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthHandler(ctx context.Context,
	collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json: "expires"`
}

// swagger:operation POST /signin auth signIn
// Login with username and password
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '401':
//         description: Invalid credentials
func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	h := sha256.New()
	newPassword := hex.EncodeToString(h.Sum([]byte(user.Password)))

	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
		"password": newPassword,
	})

	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()
	// expirationTime := time.Now().Add(10 * time.Minute)
	// claims := &Claims{
	// 	Username: user.Username,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: expirationTime.Unix(),
	// 	},
	// }

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256,
	// 	claims)

	// tokenString, err := token.SignedString([]byte(os.
	// 	Getenv("JWT_SECRET")))

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError,
	// 		gin.H{
	// 			"error": err.Error(),
	// 		})
	// 	return
	// }
	// jwtOutput := JWTOutput{
	// 	Token:   tokenString,
	// 	Expires: expirationTime,
	// }

	c.JSON(http.StatusOK, gin.H{"message": "User signed in"})
}

// swagger:operation POST /refresh auth refresh
// Get new token in exchange for an old one
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Token is new and doesn't need
//                      a refresh
//     '401':
//         description: Invalid credentials
func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	session := sessions.Default(c)
	sessionToken := session.Get("token")
	sessionUser := session.Get("username")
	if sessionToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session cookie"})
		return
	}

	sessionToken = xid.New().String()
	session.Set("username", sessionUser.(string))
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "New session issued"})
	// tokenValue := c.GetHeader("Authorization")
	// claims := &Claims{}
	// tkn, err := jwt.ParseWithClaims(tokenValue, claims,
	// 	func(token *jwt.Token) (interface{}, error) {
	// 		return []byte(os.Getenv("JWT_SECRET")), nil
	// 	})

	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
	// if tkn == nil || !tkn.Valid {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"error": "Invalid Token",
	// 	})
	// 	return
	// }
	// if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) >
	// 	30*time.Second {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Token is not expired yet",
	// 	})
	// 	return
	// }

	// expirationTime := time.Now().Add(5 * time.Minute)
	// claims.ExpiresAt = expirationTime.Unix()
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenString, err := token.SignedString(os.Getenv(
	// 	"JWT_SECRET",
	// ))

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError,
	// 		gin.H{
	// 			"error": err.Error(),
	// 		})
	// 	return
	// }

	// jwtOutput := JWTOutput{
	// 	Token:   tokenString,
	// 	Expires: expirationTime,
	// }

	// c.JSON(http.StatusOK, jwtOutput)
}

func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "Signed out",
	})
}
