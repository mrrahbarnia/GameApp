package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/mrrahbarnia/GameApp/entity"
)

func (d *PostgreSQLDB) IsPhoneNumberExist(phoneNumber string) (bool, error) {
	var userID uint

	err := d.db.QueryRow("SELECT id FROM users WHERE phone_number=$1", phoneNumber).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, fmt.Errorf("Cannot run the SQL query due to: %w", err)
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
		return entity.User{}, fmt.Errorf("Cannot run the SQL command due to: %w", err)
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

		return u, false, fmt.Errorf("Cannot run the SQL query due to: %w", err)
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

		return u, false, fmt.Errorf("Cannot run the SQL query due to: %w", err)
	}

	return u, true, nil
}
