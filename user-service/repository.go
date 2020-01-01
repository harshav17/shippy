package main

import (
	pb "github.com/harshav17/shippy/user-service/proto/user"
	"github.com/jinzhu/gorm"
)

//Repository -
type Repository interface {
	GetAll() ([]*pb.User, error)
	Get(id string) (*pb.User, error)
	Create(user *pb.User) error
	GetByEmail(email string) (*pb.User, error)
}

//UserRepository -
type UserRepository struct {
	db *gorm.DB
}

//GetAll -
func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

//Get -
func (repo *UserRepository) Get(id string) (*pb.User, error) {
	var user *pb.User
	user.Id = id
	if err := repo.db.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

//GetByEmail -
func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	user := &pb.User{}
	if err := repo.db.Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

//Create -
func (repo *UserRepository) Create(user *pb.User) error {
	if err := repo.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
