package repository

import (
	"context"
	"users-task/models"
)

type UserRepository interface {
	List(ctx context.Context) (*[]models.User, error)
	Get(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, user *models.User) (int, error)
	Update(ctx context.Context, id int, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id int) error
	AddFile(ctx context.Context, id int) error
}

/*
func List(dest interface{}, query string) error {
	err := dbList(dest, query)
	if err != nil {
		return err
	}
	return nil
}


func Get(ID string, dest interface{}, query string) error {
	err := dbGet(ID, dest, query)
	if err != nil {
		return err
	}
	return nil
}

func Create(data interface{}, query string) error {
	err := dbCreate(data, query)
	if err != nil {
		return err
	}
	return nil
}

func Update(ID string, data interface{}, query string) error {
	err := dbUpdate(ID, data, query)
	if err != nil {
		return err
	}
	return nil
}

func Delete(ID string, query string) error {
	err := dbDelete(ID, query)
	if err != nil {
		return err
	}
	return nil
}
*/
