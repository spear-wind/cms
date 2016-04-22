package user

import (
	"errors"

	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/cfmgo/params"
	"gopkg.in/mgo.v2/bson"
)

type mongoUserRepository struct {
	Collection cfmgo.Collection
}

type userRecord struct {
	RecordID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserID           int64         `bson:"user_id",json:"match_id"`
	FacebookID       string        `bson:"fb_id",json:"fb_id"`
	Email            string        `bson:"email",json:"email"`
	FirstName        string        `bson:"first_name",json:"first_name"`
	LastName         string        `bson:"last_name",json:"last_name"`
	Hash             string        `bson:"hash",json:"hash"`
	Verified         bool          `bson:"verified",json:"verified"`
	VerificationCode string        `bson:"verification_code",json:"verification_code"`
}

func NewMongoUserRepository(col cfmgo.Collection) *mongoUserRepository {
	return &mongoUserRepository{
		Collection: col,
	}
}

func (repo *mongoUserRepository) Add(user *User) (err error) {
	repo.Collection.Wake()
	ur := toUserRecord(user)
	_, err = repo.Collection.UpsertID(ur.RecordID, ur)
	return
}

func (repo *mongoUserRepository) Update(user *User) (err error) {
	repo.Collection.Wake()
	foundUser, err := repo.getMongoUser(user.ID)
	if err == nil {
		ur := toUserRecord(user)
		ur.RecordID = foundUser.RecordID
		_, err = repo.Collection.UpsertID(ur.RecordID, ur)

	}

	return
}

func (repo *mongoUserRepository) listUsers() (users []*User) {
	repo.Collection.Wake()
	var ur []userRecord
	_, err := repo.Collection.Find(cfmgo.ParamsUnfiltered, &ur)
	if err == nil {
		users = make([]*User, len(ur))
		for k, v := range ur {
			users[k] = toUser(&v)
		}
	}

	return
}

func (repo *mongoUserRepository) getUser(id int64) (user *User, err error) {
	var ur *userRecord
	ur, err = repo.getMongoUser(id)
	if ur != nil && err == nil {
		user = toUser(ur)
	}

	return
}

func (repo *mongoUserRepository) getMongoUser(id int64) (user *userRecord, err error) {
	var users []userRecord
	query := bson.M{"user_id": id}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := repo.Collection.Find(params, &users)
	if count == 0 {
		err = errors.New("User not found")
	}
	if err == nil {
		user = &users[0]
	}

	return
}

func (repo *mongoUserRepository) FindByVerificationCode(verificationCode string) (user *User) {
	var users []userRecord
	query := bson.M{"verification_code": verificationCode}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := repo.Collection.Find(params, &users)
	if count == 0 {
		err = errors.New("User not found")
	}
	if err == nil {
		user = toUser(&users[0])
	}

	return
}

func (repo *mongoUserRepository) FindByEmail(email string) (user *User) {
	var users []userRecord
	query := bson.M{"email": email}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := repo.Collection.Find(params, &users)
	if count == 0 {
		err = errors.New("User not found")
	}
	if err == nil {
		user = toUser(&users[0])
	}

	return
}

func (repo *mongoUserRepository) FindByFacebookID(facebookID string) (user *User) {
	var users []userRecord
	query := bson.M{"facebook_id": facebookID}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := repo.Collection.Find(params, &users)
	if count == 0 {
		err = errors.New("User not found")
	}
	if err == nil {
		user = toUser(&users[0])
	}

	return
}

func (repo *mongoUserRepository) Exists(user *User) bool {
	ur, err := repo.getMongoUser(user.ID)
	return err == nil && ur != nil
}

func toUserRecord(u *User) (ur *userRecord) {
	ur = &userRecord{
		RecordID:         bson.NewObjectId(),
		UserID:           u.ID,
		FacebookID:       u.FacebookID,
		Email:            u.Email,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Hash:             u.hash,
		Verified:         u.Verified,
		VerificationCode: u.VerificationCode,
	}
	return
}

func toUser(ur *userRecord) (u *User) {
	u = &User{
		ID:               ur.UserID,
		FacebookID:       ur.FacebookID,
		Email:            ur.Email,
		FirstName:        ur.FirstName,
		LastName:         ur.LastName,
		hash:             ur.Hash,
		Verified:         ur.Verified,
		VerificationCode: ur.VerificationCode,
	}
	return
}
