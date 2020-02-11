package utils

import (
"encoding/json"
"fmt"
"github.com/mandoju/BindAPI/models/config"
"io/ioutil"
"os"
)

//GetJwtKey gets the secret from configuration file
func GetJwtKey() ([]byte, error) {
	path := "./config/JWT.json"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
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
		fmt.Println(err)
		return nil, err
	}
	secret := []byte(configuration.Secret)
	return secret, nil
}

func GetDnsSecretKey() (string, error) {
	path := "./config/DNS.json"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return "", err
	}
	return configuration.Secret, nil
}
