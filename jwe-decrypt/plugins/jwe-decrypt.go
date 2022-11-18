/*
	A custom plugin in Go to decrypt JWE
*/

package main

import (
	"crypto/x509"
	"encoding/base64"
	"strings"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"

	"gopkg.in/square/go-jose.v2"
)

func main() {
	server.StartServer(New, Version, Priority)
}

var Version = "0.1"
var Priority = 1

// Property of the Plugin
type Config struct {
	PrivateKey string
}

func New() interface{} {
	return &Config{}
}

//------------------------------------
// Decrypt JWE by using a Private Key
//------------------------------------
func decrypt(kong *pdk.PDK, myPrivateKey string, JWE string) string {

	kong.Log.Notice("*** Decrypt - Begin ***")

	kong.Log.Notice("decrypt - myPrivateKey=", myPrivateKey)

	kong.Log.Notice("decrypt - DecodeString")
	privateKeyBytes, err := base64.StdEncoding.DecodeString(myPrivateKey)
	if err != nil {
		kong.Log.Err("err: ", string(err.Error()))
		return ""
	}
	kong.Log.Notice("decrypt - ParsePKCS1PrivateKey")
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		kong.Log.Err("err: ", string(err.Error()))
		return ""
	}
	kong.Log.Notice("decrypt - ParseEncrypted")
	encryptedJwe, err := jose.ParseEncrypted(JWE)
	if err != nil {
		kong.Log.Err("err: ", string(err.Error()))
		return ""
	}

	kong.Log.Notice("decrypt - Decrypt")
	JWT, err := encryptedJwe.Decrypt(privateKey)

	if err != nil {
		kong.Log.Err("err: ", string(err.Error()))
		return ""
	}
	kong.Log.Notice("*** Decrypt - End ***")
	return string(JWT)
}

//-----------------------------------------------------------------------------
// Access => Executed for every request from a Client and before it is being
// proxied to the upstream service (i.e. producer)
//-----------------------------------------------------------------------------
func (conf Config) Access(kong *pdk.PDK) {

	var bearer = -1
	var JWE = ""

	kong.Log.Notice("*** jwe-decrypt - Begin Access() ***")
	header := map[string][]string{"content-type": {"application/json"}}

	// Get Header from the request Consumer
	authorization, _ := kong.Request.GetHeader("Authorization")
	if authorization != "" {
		// Find Bearer
		bearer = strings.Index(authorization, "Bearer ")
		if bearer != -1 {
			JWE = authorization[bearer+len("Bearer "):]
		}
	}
	if bearer == -1 {
		kong.Log.Err("Unable to get 'Authorization' header: ")
		kong.Response.Exit(500, `{"Error": "Unable to get 'Authorization' header"}`, header)
	}

	kong.Log.Notice("JWE found in 'Authorization' header=", JWE)

	// Decrypt JWE payload
	var JWT = decrypt(kong, conf.PrivateKey, JWE)
	if JWT == "" {
		kong.Log.Err("Unable to decrypt JWE payload")
		kong.Response.Exit(500, `{"Error": "Unable to decrypt JWE payload"}`, header)
	}

	// replace the JWE by the JWT in the Authorization header
	kong.ServiceRequest.SetHeader("Authorization", "Bearer "+JWT)

	kong.Log.Notice("*** jwe-decrypt - End Access() ***")
}

//-----------------------------------------------------------------------------
// Response => Executed after the whole response has been received from the
// upstream service (i.e. Producer), but before sending any part of it to
// the Client
//-----------------------------------------------------------------------------
func (conf Config) Response(kong *pdk.PDK) {
	kong.Log.Notice("*** jwe-decrypt - Begin Response() ***")

	// ...

	kong.Log.Notice("*** jwe-decrypt- End Response() ***")
}
