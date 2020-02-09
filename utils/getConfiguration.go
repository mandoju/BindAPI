package utils

import (
	"encoding/json"
	"github.com/mandoju/BindAPI/models/config"
	"io/ioutil"
	"os"
)

//GetJwtKey gets the secret from configuration file
func GetJwtKey() (string,error) {
	path := "../config/JWT.json"
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(file)

	// we initialize our Users array
	var configuration config.ConfigurationJWT

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &configuration)
	if err != nil {
		return "",err
	}
	return configuration.Secret, nil
}
