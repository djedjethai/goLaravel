package data

// go get golang.org/x/crypto/bcrypt

import (
	"errors"
	// "fmt"
	"github.com/djedjethai/goframework"
	up "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	// `` if for upper.io(the ORM) to work with it
	// the field ID will match id or nothing(as omitempty)
	ID        int       `db:"id,omitempty"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Active    string    `db:"user_active"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Token     Token     `db:"-"` // means to the db to ignore that field
}

// that func is to over write the table name in case of legacy code
// could be very usefull
func (u *User) Table() string {
	return "users"
}

// add the validator to the model, got it from the imported module goframework
func (u *User) Validate(validator *goframework.Validation) {
	validator.Check(u.LastName != "", "last_name", "Last name must be provided")
	validator.Check(u.FirstName != "", "first_name", "First name must be provided")
	validator.Check(u.Email != "", "email", "Email must be provided")
	validator.IsMail("email", u.Email)
}

// up.Condition is where we are starting to use upper.io
func (u *User) GetAll() ([]*User, error) {
	// upper refers to things stored in db
	collection := upper.Collection(u.Table())

	var all []*User

	res := collection.Find().OrderBy("last_name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	var theUser User

	collection := upper.Collection(u.Table())

	res := collection.Find(up.Cond{"email =": email})
	err := res.One(&theUser) // as we expect only one res
	if err != nil {
		return nil, err
	}

	// add the user's token to its record to return
	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id =": theUser.ID, "expiry >": time.Now()}).OrderBy("created_at desc")
	err = res.One(&token)
	if err != nil {
		// in case user have no token
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}

	}

	// set the token to the user
	// if there is no token the token field will be empty
	theUser.Token = token

	return &theUser, nil
}

func (u *User) Get(id int) (*User, error) {
	var theUser User

	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&theUser)
	if err != nil {
		return nil, err
	}

	// add the user's token to its record to return
	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id =": theUser.ID, "expiry >": time.Now()}).OrderBy("created_at desc")
	err = res.One(&token)
	if err != nil {
		// in case user have no token
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}

	}

	// set the token to the user
	// if there is no token the token field will be empty
	theUser.Token = token

	return &theUser, nil
}

func (u *User) Update(theUser User) error {
	theUser.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	res := collection.Find(theUser.ID)
	err := res.Update(&theUser)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete(id int) error {
	collection := upper.Collection(u.Table())
	res := collection.Find(id)

	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Insert(theUser User) (int, error) {
	// []byte(theUser.Password) cast to a slice of bytes
	newHash, err := bcrypt.GenerateFromPassword([]byte(theUser.Password), 12)
	if err != nil {
		return 0, err
	}

	theUser.CreatedAt = time.Now()
	theUser.UpdatedAt = time.Now()
	// get back from bcrypt a slice of byte, so we cast it
	theUser.Password = string(newHash)

	collection := upper.Collection(u.Table())
	res, err := collection.Insert(theUser)
	if err != nil {
		return 0, err
	}

	// the res.ID can be in various type, depending of the db we are dealing with
	// so need to make sur the type is int
	id := GetInsertID(res.ID())
	return id, nil
}

func (u *User) ResetPassword(id int, password string) error {

	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	theUser, err := u.Get(id)
	if err != nil {
		return err
	}

	theUser.Password = string(newHash)

	err = u.Update(*theUser)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		// wrong password
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		// any other err
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *User) CheckForRememberToken(id int, token string) bool {
	var rememberToken RememberToken
	rt := RememberToken{}

	collection := upper.Collection(rt.Table())
	res := collection.Find(up.Cond{"user_id": id, "remember_token": token})
	err := res.One(&rememberToken)
	return err == nil // return bool
}
