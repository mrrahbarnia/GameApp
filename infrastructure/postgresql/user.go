package postgresql

import (
	"database/sql"

	"github.com/mrrahbarnia/GameApp/entity"
	"github.com/mrrahbarnia/GameApp/pkg/errmsg"
	"github.com/mrrahbarnia/GameApp/pkg/richerror"
)

func (d *PostgreSQLDB) IsPhoneNumberExist(phoneNumber string) (bool, error) {
	var userID uint

	err := d.db.QueryRow("SELECT id FROM users WHERE phone_number=$1", phoneNumber).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, richerror.New("postgresql.IsPhoneNumberExist").
				WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
				WithKind(richerror.KindUnexpected)
		}
	}

	return true, nil
}

func (d *PostgreSQLDB) Register(u entity.User) (entity.User, error) {
	var userID uint

	if err := d.db.QueryRow(
		"INSERT INTO users(name, phone_number, hashed_password) VALUES($1, $2, $3) RETURNING id",
		u.Name,
		u.PhoneNumber,
		u.HashedPassword,
	).Scan(&userID); err != nil {
		return entity.User{}, richerror.New("postgresql.Register").
			WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return entity.User{
		ID:             userID,
		Name:           u.Name,
		PhoneNumber:    u.PhoneNumber,
		HashedPassword: u.HashedPassword,
	}, nil
}

func (d *PostgreSQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error) {
	u := entity.User{}
	err := d.db.QueryRow(
		"SELECT id, phone_number, name, hashed_password FROM users WHERE phone_number=$1",
		phoneNumber,
	).Scan(&u.ID, &u.PhoneNumber, &u.Name, &u.HashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return u, false, nil
		}

		return u, false, richerror.New("postgresql.GetUserByPhoneNumber").
			WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return u, true, nil
}

func (d *PostgreSQLDB) GetUserById(userId uint) (entity.User, bool, error) {
	u := entity.User{}
	err := d.db.QueryRow(
		"SELECT id, phone_number, name, hashed_password FROM users WHERE id=$1",
		userId,
	).Scan(&u.ID, &u.PhoneNumber, &u.Name, &u.HashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return u, false, nil
		}

		return u, false,
			richerror.New("postgresql.GetUserById").
				WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
				WithKind(richerror.KindUnexpected)
	}

	return u, true, nil
}
