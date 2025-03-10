// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	User
	AccessToken string `json:"accessToken"`
	jwt.StandardClaims
}

func ParseJwtToken(token string) (*User, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(authConfig.JwtSecret), nil
	})

	if tokenClaims == nil {
		return nil, err
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok || !tokenClaims.Valid {
		return nil, TokenInvalidError()
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, TokenExpiredError(claims.ExpiresAt)
	}

	return &claims.User, nil
}
