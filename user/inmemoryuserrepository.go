package user

import "errors"

type inMemoryRepository struct {
	users map[int64]*User
}

func NewInMemoryRepository() *inMemoryRepository {
	repo := &inMemoryRepository{}
	repo.users = make(map[int64]*User)
	return repo
}

func (repo *inMemoryRepository) Add(user *User) (err error) {
	user.ID = int64(len(repo.users) + 1)
	repo.users[user.ID] = user
	return err
}

func (repo *inMemoryRepository) Update(user *User) (err error) {
	repo.users[user.ID] = user
	return err
}

func (repo *inMemoryRepository) listUsers() (users []*User) {
	for _, user := range repo.users {
		users = append(users, user)
	}

	return users
}

func (repo *inMemoryRepository) getUser(userID int64) (user *User, err error) {
	found := false

	for _, target := range repo.users {
		if userID == target.ID {
			user = target
			found = true
		}
	}
	if !found {
		err = errors.New("Could not find user in repository")
	}
	return user, err
}

func (repo *inMemoryRepository) FindByVerificationCode(verificationCode string) (user *User) {
	for _, target := range repo.users {
		if target.VerificationCode == verificationCode {
			user = target
			break
		}
	}

	return user
}

func (repo *inMemoryRepository) FindByEmail(email string) (user *User) {
	for _, target := range repo.users {
		if target.Email == email {
			user = target
			break
		}
	}

	return user
}

func (repo *inMemoryRepository) FindByFacebookID(facebookID string) (user *User) {
	for _, target := range repo.users {
		if target.FacebookID == facebookID {
			user = target
			break
		}
	}

	return user
}

func (repo *inMemoryRepository) Exists(user *User) bool {
	for _, target := range repo.users {
		if user.Email == target.Email {
			return true
		}
	}

	return false
}
