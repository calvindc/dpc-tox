package main

import (
	"errors"
	"github.com/calvindc/dpc-tox/cmd/webtox/httpserve"
	"github.com/calvindc/dpc-tox/librarywrapper/libtox"
	"io/ioutil"
	"log"
	"os"
)

// getUserStatusAsString returns a string representing the given Tox user status
// status  the Tox user status to be converted
func getUserStatusAsString(status libtox.ToxUserStatus) string {
	switch status {
	case libtox.TOX_USERSTATUS_NONE:
		return "NONE"
	case libtox.TOX_USERSTATUS_AWAY:
		return "AWAY"
	case libtox.TOX_USERSTATUS_BUSY:
		return "BUSY"
	default:
		return "INVALID"
	}
}

// getUserStatusFromString returns the Tox user status represented by the given
// user status string
// status  the user status as a string to be converted
func getUserStatusFromString(status string) libtox.ToxUserStatus {
	switch status {
	case "NONE":
		return libtox.TOX_USERSTATUS_NONE
	case "AWAY":
		return libtox.TOX_USERSTATUS_AWAY
	case "BUSY":
		return libtox.TOX_USERSTATUS_BUSY
	default:
		return libtox.TOX_USERSTATUS_NONE
	}
}

// saveData writes the current Tox saveData to a file
// t         the libtox instance whichs saveData will be stored
// filepath  the path to the file the saveData will be stored in
func saveData(t *libtox.Tox, filepath string) error {
	if len(filepath) == 0 {
		return errors.New("Empty path")
	}

	data, err := t.GetSavedata()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath, data, 0644)
	return err
}

// loadData reads a file and returns its contents as a byte array
// filepath  the path to the file
func loadData(filepath string) ([]byte, error) {
	if len(filepath) == 0 {
		return nil, errors.New("Empty path")
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return data, err
}

// fileExists returns true if the given file or directory exists, otherwise false
// path  the given file or directory
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// storeDefaultHTTPAuth generates a random password and stores it into the
// database (used for initialisation)
func storeDefaultHTTPAuth() (string, string, string) {
	salt, err := httpserve.RandomString(32)
	if err != nil {
		panic("could not generate salt")
	}

	plainPass, err := httpserve.RandomString(32)
	if err != nil {
		panic("could not generate salt")
	}

	user := CFG_DEFAULT_AUTH_USER
	pass := httpserve.Sha512Sum(plainPass + salt)

	log.Println("Info: Username reset to: ", user)
	log.Println("Info: Password reset to: ", plainPass)

	storage.StoreKeyValue("settings_auth_user", user)
	storage.StoreKeyValue("settings_auth_pass", pass)
	storage.StoreKeyValue("settings_auth_salt", salt)

	return user, pass, salt
}
