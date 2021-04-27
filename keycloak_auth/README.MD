# KEYCLOAK Login endpoint For Fiber

## Endpoint
1. `POST: {BASEURL}/token` 
    #### 1.1 Get Token
    > **Header:**
    >
    > > `Content-Type: multipart/form-data;`
    > 
    > **FromData: ** 
    > 
    > - `username` - string username
    > - `password` - string password
    > - `scope` - string scope for information

    #### 1.2 Refresh Token
    > **Header**
    > > `Content-Type: multipart/form-data;`
    >
    > **FromData ** 
    > - `grant_type=refresh_token` -  set keyname `grant_type` =`refresh_token`
     > - `refresh_token` - refresh token from [Get Token](#get_token)

2. `GET: {BASEURL}/userinfo`

   - **Header**

   > `Authorization: {Access_Token}`

3. `GET: {BASEURL}/logout

   > **Header**
   >
   > >  `Authorization: {Access_Token}`



## Example:

```go
package main

import (
	"github.com/attapon-th/go-pkg/keycloak_auth"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	keycloak_auth.New(keycloak_auth.Config{
		FiberRouter:    app,
		BaseInstantUrl: "https://example.go.th/auth/",
		RealmName:      "realm_test",
		ClientID:       "client_test",
		ClientSecret:   "secret_key",
	})
	_ = app.Listen("127.0.0.1:8080")
}

```

