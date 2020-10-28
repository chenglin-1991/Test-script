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

package client

import (
    "fmt"
    "net/http"
    "encoding/json"
    . "github.com/xtao/xtor/common"
)

func SignXtorReqHeader(req *http.Request, account *UserAccountInfo) error {
    var user, uid string

    signedHeaders := []string{}

    if _, exists := req.Header["Content-Type"]; exists {
        signedHeaders = append(signedHeaders, "Content-Type")
    }
    accHeader := make(map[string]interface{})

    if account.Username != "" {
        user = account.Username
    }

    accHeader["user"] = user

    AesAkey := NewAesEncrypt(XTOR_API_KEY_SEED)
    AesSkey := NewAesEncrypt(XTOR_SECURITY_KEY_SEED)
    aKey, err := AesAkey.Encrypt(user)
    if err != nil {
        return err
    }

    var sKey string

    sKey, err = AesSkey.Encrypt(user)
    if err != nil {
        return err
    }

    if account.Uid != "" {
        uid = account.Uid
        accHeader["uid"] = uid
    }

    js, err := json.Marshal(accHeader)
    if err != nil {
        return err
    }

    // Encrypt the account information
    dKey := sKey
    if len(dKey) < 16 {
        dKey = fmt.Sprintf("%16s", sKey)
    }

    aesEnc := NewAesEncrypt(dKey)
    aHeader, err := aesEnc.Encrypt(string(js))
    if err != nil {
        return err
    }

    signedHeaders = append(signedHeaders, AccountHeader)
    req.Header.Add(AccountHeader, aHeader)

    options := Options {
        SignedHeaders: signedHeaders,
    }

    str, err := StringToSign(req, &options)
    if err != nil {
        return err
    }

    signature := SignString(str, sKey)

    authHeader := fmt.Sprintf("APIKey=%s,Signature=%s", aKey, signature)
    req.Header.Add("Authorization", authHeader)

    return nil
}
