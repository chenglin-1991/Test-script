/*
 Copyright (c) 2016-2017 XTAO technology <www.xtaotech.com>
 All rights reserved.

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions
 are met:
  1. Redistributions of source code must retain the above copyright
     notice, this list of conditions and the following disclaimer.
  2. Redistributions in binary form must reproduce the above copyright
     notice, this list of conditions and the following disclaimer in the
     documentation and/or other materials provided with the distribution.

  THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
  ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
  ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
  FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
  DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
  OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
  HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
  LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
  OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
  SUCH DAMAGE.
*/
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "github.com/xtao/xtor/common"
)

var globalAdminSkey string = XTOR_INITIAL_ADMIN_PASSWORD

func GetGlobalAdminSkey() string {
	return globalAdminSkey
}

func SetGlobalAdminSkey(key string) {
	globalAdminSkey = key
}

func AuthHeaderFromClient(req *http.Request) (error, *UserAccountInfo) {
	signedHeaders := []string{}

	if _, exists := req.Header["Content-Type"]; exists {
		signedHeaders = append(signedHeaders, "Content-Type")
	}

	signedHeaders = append(signedHeaders, AccountHeader)

	options := Options{
		SignedHeaders: signedHeaders,
		SecretKey: KeyLocator(func(apiKey string) (string, string) {
			AesAkey := NewAesEncrypt(XTOR_API_KEY_SEED)
			AesSkey := NewAesEncrypt(XTOR_SECURITY_KEY_SEED)

			user, err := AesAkey.Decrypt(apiKey)
			if err != nil {
				return "", ""
			}
			skey, err := AesSkey.Encrypt(user)
			if err != nil {
				return "", user
			}
			return skey, user
		}),
	}

	err, user, skey := HMACAuth(options, req)
	if err != nil {
		Logger.Printf("authenticate failure: %s\n",
			err.Error())
		return err, nil
	}

	Logger.Printf("user:%s,Skey:%s \n", user, skey)

	dkey := skey
	if len(dkey) < 16 {
		dkey = fmt.Sprintf("%16s", skey)
	}

	aesEnc := NewAesEncrypt(dkey)
	enAccountHeader := req.Header.Get(AccountHeader)
	accountHeader, err := aesEnc.Decrypt(enAccountHeader)
	if err != nil {
		Logger.Printf("Decrypt header failure: %s\n",
			err.Error())
		return err, nil
	}

	account := UserAccountInfo{}
	err = json.Unmarshal([]byte(accountHeader), &account)

	if err != nil {
		Logger.Printf("Failed to parse JSON account:%s\n",
			accountHeader)
		return err, nil
	}

	return nil, &account
}
