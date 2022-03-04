package data

import (
	"crypto/sha256"
	"encoding/base32"
	"errors"
	up "github.com/upper/db/v4"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	// `` if for upper.io(the ORM) to work with it
	ID        int       `db:"id,omitempty" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	FirstName string    `db:"first_name" json:"first_name"`
	Email     string    `db:"email" json:"email"`
	PlainText string    `db:"token" json:"token"`
	Hash      []byte    `db:"token_hash" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Expires   time.Time `db:"expiry" json:"expiry"`
}

func (t *Token) Table() string {
	return "tokens"
}

func (t *Token) GetUserForToken(token string) (*User, error) {
	// get the token from token table
	var u User
	var theToken Token

	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"token =": token})
	err := res.One(&theToken)
	if err != nil {
		return nil, err
	}

	// get user using the token
	collection = upper.Collection("users")
	res = collection.Find(up.Cond{"id =": theToken.UserID})
	err = res.One(&u)
	if err != nil {
		return nil, err
	}

	u.Token = theToken

	return &u, nil
}

func (t *Token) GetTokensForUser(id int) ([]*Token, error) {
	var tokens []*Token

	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"user_id =": id})
	err := res.All(&tokens)
	// from the prof but even no token, no err are returned
	if err != nil {
		return nil, err
	}

	// so i add That
	if len(tokens) < 1 {
		return nil, errors.New("No token found")
	}

	return tokens, nil
}

// get the token value from the token.id
func (t *Token) Get(id int) (*Token, error) {
	var token Token

	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"id =": id})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *Token) GetByToken(plainText string) (*Token, error) {
	var token Token

	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token =": plainText})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *Token) DeleteById(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"id =": id})
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) DeleteByToken(plainText string) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token =": plainText})
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) Insert(token Token, u User) error {
	collection := upper.Collection(t.Table())

	// delete existing token first
	res := collection.Find(up.Cond{"user_id =": u.ID})
	err := res.Delete()
	if err != nil {
		return nil
	}

	// make sure the token is properly setted
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	token.FirstName = u.FirstName
	token.Email = u.Email

	_, err = collection.Insert(token)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) GenerateToken(userID int, ttl time.Duration) (*Token, error) {

	userIDString := strconv.Itoa(userID)

	// my fix: we generate the Expires date in .UTC()
	// as it look like my db return token in UTC format
	// but !!! after it may create problems(all validations must be done in UTC)...
	token := &Token{
		UserID:  userIDString,
		Expires: time.Now().Add(ttl).UTC(),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// got the token encoded on base32
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	// hash the token (we cast the token.PlainText to a slice of bytes)
	hash := sha256.Sum256([]byte(token.PlainText))

	// convert our hash into an array
	token.Hash = hash[:]

	return token, nil
}

func (t *Token) AuthenticateToken(r *http.Request) (*User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("No authorization header received")
	}

	// the token format will be
	// Bearer thetokenblblblblbbl
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("No authorization header received")
	}

	token := headerParts[1]

	// all our token are of this size
	if len(token) != 26 {
		return nil, errors.New("Token wrong size")
	}

	tkn, err := t.GetByToken(token)
	if err != nil {
		return nil, errors.New("No matching token found")
	}

	// if mytest.Before(time.Now()) {
	if tkn.Expires.Before(time.Now()) {
		return nil, errors.New("Expired token")
	}

	user, err := tkn.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("No matching user found")
	}

	return user, nil
}

func (t *Token) ValidToken(token string) (bool, error) {
	user, err := t.GetUserForToken(token)
	if err != nil {
		return false, errors.New("No matching user found")
	}

	if user.Token.PlainText == "" {
		return false, errors.New("No matching token found")
	}

	if user.Token.Expires.Before(time.Now()) {
		return false, errors.New("Expired token")
	}

	return true, nil
}
