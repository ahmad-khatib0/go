package handlers

import (
	"crypto/sha256"
	"net/http"
	"time"

	"github.com/ahmad-khatib0/go/distributed-services/gin/api/models"
	"github.com/dgrijalva/jwt-go"
	gs "github.com/gin-contrib/sessions"
	g "github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// swagger:operation POST /signin auth signIn
// Login with username and password
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'401':
//	    description: Invalid credentials
func (handler *AuthHandler) SignInHandler(c *g.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, g.H{"error": err.Error()})
		return
	}

	h := sha256.New()

	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
		"password": string(h.Sum([]byte(user.Password))),
	})
	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, g.H{"error": "Invalid username or password"})
		return
	}

	sessionToken := xid.New().String()
	session := gs.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, g.H{"message": "User signed in"})
}

// swagger:operation POST /refresh auth refresh
// Refresh token
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'401':
//	    description: Invalid credentials
func (handler *AuthHandler) RefreshHandler(c *g.Context) {
	session := gs.Default(c)
	sessionToken := session.Get("token")
	sessionUser := session.Get("username")
	if sessionToken == nil {
		c.JSON(http.StatusUnauthorized, g.H{"error": "Invalid session cookie"})
		return
	}

	sessionToken = xid.New().String()
	session.Set("username", sessionUser.(string))
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, g.H{"message": "New session issued"})
}

// swagger:operation POST /signout auth signOut
// Signing out
// ---
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *AuthHandler) SignOutHandler(c *g.Context) {
	session := gs.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, g.H{"message": "Signed out..."})
}

func (handler *AuthHandler) AuthMiddleware() g.HandlerFunc {
	return func(c *g.Context) {
		session := gs.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, g.H{
				"message": "Not logged",
			})
			c.Abort()
		}
		c.Next()
	}
}
