package phonenumber

import "strconv"

func IsValid(phoneNumber string) bool {
	// TODO - use REGEX for better validation
	if len(phoneNumber) != 11 {
		return false
	}

	if phoneNumber[0:2] != "09" {
		return false
	}

	if _, err := strconv.Atoi(phoneNumber); err != nil {
		return false
	}

	return true
}
