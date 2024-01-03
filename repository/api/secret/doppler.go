package secret

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	cg "github.com/pharmaniaga/auth-user/domain/constants/general"
	mg "github.com/pharmaniaga/auth-user/domain/model/general"
	"github.com/pharmaniaga/auth-user/domain/utils"
)

func GetCredentials(url, key, password string) (*mg.AppService, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Basic "+utils.BasicAuth(key, password))

	client := http.Client{
		Timeout: cg.APITimeDuration5s,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var credentials mg.SectionService
		err = json.Unmarshal(body, &credentials)
		if err != nil {
			return nil, err
		}

		maxIdleConnsRead, err := utils.StrToInt(credentials.DatabaseReadMaxIdleConns)
		if err != nil {
			return nil, err
		}

		maxIdleConnsWrite, err := utils.StrToInt(credentials.DatabaseWriteMaxIdleConns)
		if err != nil {
			return nil, err
		}

		maxOpenConnsRead, err := utils.StrToInt(credentials.DatabaseReadMaxOpenConns)
		if err != nil {
			return nil, err
		}

		maxOpenConnsWrite, err := utils.StrToInt(credentials.DatabaseWriteMaxOpenConns)
		if err != nil {
			return nil, err
		}

		maxLifeTimeRead, err := utils.StrToInt(credentials.DatabaseReadMaxLifeTime)
		if err != nil {
			return nil, err
		}

		maxLifeTimeWrite, err := utils.StrToInt(credentials.DatabaseWriteMaxLifeTime)
		if err != nil {
			return nil, err
		}

		redisPort, err := utils.StrToInt(credentials.RedisPort)
		if err != nil {
			return nil, err
		}

		redisMinIdleConns, err := utils.StrToInt(credentials.RedisMinIdleConns)
		if err != nil {
			return nil, err
		}

		authJWTIsActive, err := utils.StrToBool(credentials.AuthorizationJWTIsActive)
		if err != nil {
			return nil, err
		}

		jwtAccessDuration, err := utils.StrToInt(credentials.AuthorizationJWTAccessTokenDuration)
		if err != nil {
			return nil, err
		}

		jwtRefreshDuration, err := utils.StrToInt(credentials.AuthorizationJWTRefreshTokenDuration)
		if err != nil {
			return nil, err
		}

		return &mg.AppService{
			App: mg.AppAccount{
				Name:         credentials.AppName,
				Environtment: credentials.AppEnvirontment,
				URL:          credentials.AppURL,
				Port:         credentials.AppPort,
				SecretKey:    credentials.AppSecretKey,
			},
			Route: mg.RouteAccount{
				Methods: strings.Split(credentials.RouteMethods, ","),
				Headers: strings.Split(credentials.RouteHeaders, ","),
				Origins: strings.Split(credentials.RouteOrigins, ","),
			},
			DatabaseUser: mg.DatabaseUser{
				Read: mg.DBDetailUser{
					Username:     credentials.DatabaseReadUsername,
					Password:     credentials.DatabaseReadPassword,
					URL:          credentials.DatabaseReadURL,
					Port:         credentials.DatabaseReadPort,
					DBName:       credentials.DatabaseReadDBName,
					MaxIdleConns: maxIdleConnsRead,
					MaxOpenConns: maxOpenConnsRead,
					MaxLifeTime:  maxLifeTimeRead,
					Timeout:      credentials.DatabaseReadTimeout,
					SSLMode:      credentials.DatabaseReadSSLMode,
				},
				Write: mg.DBDetailUser{
					Username:     credentials.DatabaseWriteUsername,
					Password:     credentials.DatabaseWritePassword,
					URL:          credentials.DatabaseWriteURL,
					Port:         credentials.DatabaseWritePort,
					DBName:       credentials.DatabaseWriteDBName,
					MaxIdleConns: maxIdleConnsWrite,
					MaxOpenConns: maxOpenConnsWrite,
					MaxLifeTime:  maxLifeTimeWrite,
					Timeout:      credentials.DatabaseWriteTimeout,
					SSLMode:      credentials.DatabaseWriteSSLMode,
				},
			},
			Redis: mg.RedisAccount{
				Username:     credentials.RedisUsername,
				Password:     credentials.RedisPassword,
				URL:          credentials.RedisURL,
				Port:         redisPort,
				MinIdleConns: redisMinIdleConns,
				Timeout:      credentials.RedisTimeout,
			},
			Authorization: mg.AuthAccount{
				JWT: mg.JWTCredential{
					IsActive:              authJWTIsActive,
					AccessTokenSecretKey:  credentials.AuthorizationJWTAccessTokenSecretKey,
					AccessTokenDuration:   jwtAccessDuration,
					RefreshTokenSecretKey: credentials.AuthorizationJWTRefreshTokenSecretKey,
					RefreshTokenDuration:  jwtRefreshDuration,
				},
				Public: mg.PublicCredential{
					SecretKey: credentials.AuthorizationPublicSecretKey,
				},
			},
			KeyData: mg.KeyAccount{
				User: credentials.KeyAccountUser,
			},
			Minio: mg.MinioSecret{
				BucketName: credentials.MinioBucketName,
				Endpoint:   credentials.MinioEndpoint,
				Key:        credentials.MinioKey,
				Secret:     credentials.MinioSecret,
				Region:     credentials.MinioRegion,
				TempFolder: credentials.MinioTempFolder,
				BaseURL:    credentials.MinioBaseURL,
			},
		}, nil
	}

	return nil, errors.New("response not ok")
}
