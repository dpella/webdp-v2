package config

import (
	"encoding/json"
	"fmt"
	"os"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/entity"
)

type Envs struct {
	Port_ext    string
	Port_int    string
	Db_host     string
	Dp_port     string
	Db_user     string
	Db_pw       string
	Db_name     string
	Root_pw     string
	Auth_key    string
	Config_path string
}

/*
Collects environment variables from container.
*/
func Getenvars() (*Envs, error) {
	// api ports
	apiPort := os.Getenv("PORT")
	if apiPort == "" {
		return nil, fmt.Errorf("%w: %s", errors.ErrMissingEnv, "api port")
	}

	internalApiPort := os.Getenv("INTERNAL_PORT")
	if internalApiPort == "" {
		return nil, fmt.Errorf("%w: %s", errors.ErrMissingEnv, "internal api port")
	}

	// root user
	pass := os.Getenv("ROOT_PASSWORD")
	if pass == "" {
		return nil, fmt.Errorf("%w: %s", errors.ErrMissingEnv, "root user")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return nil, fmt.Errorf("%w: %s", errors.ErrMissingEnv, "db setup")
	}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("%w: %s", errors.ErrMissingEnv, "config path")
	}

	skey := os.Getenv("AUTH_SIGN_KEY")
	if skey == "" {
		return nil, fmt.Errorf("%w: %s", errors.ErrMissingEnv, "auth key")
	}

	return &Envs{
		Port_ext:    apiPort,
		Port_int:    internalApiPort,
		Db_host:     dbHost,
		Dp_port:     dbPort,
		Db_user:     dbUser,
		Db_pw:       dbPassword,
		Db_name:     dbName,
		Root_pw:     pass,
		Auth_key:    skey,
		Config_path: path,
	}, nil
}

/*
Collects default engine and list of DP engines from config file.
*/
func GetDpEngines(path string) (*entity.EnginesConfig, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out entity.EnginesConfig
	err = json.Unmarshal(bs, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
