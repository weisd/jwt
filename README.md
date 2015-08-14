# jwt


## echo example
```
	var jwtSigningKeys = map[string]string{"da": "weisd"}

	j := e.Group("/jwt")
	j.Use(jwt.EchoJWTAuther(func(c *echo.Context) (key string, err error) {
		// get the clientId from header
		clientId := c.Request().Header.Get("client-id")
		key, ok := jwtSigningKeys[clientId]
		if !ok {
			return "", echo.NewHTTPError(http.StatusUnauthorized)
		}
		return key, nil
	}))

	j.Get("", func(c *echo.Context) error {
		return c.String(200, "jwt Access ok with claims %v", jwt.Claims(c))
	})

	e.Get("/test/jwt/token", func(c *echo.Context) error {
		claims := map[string]interface{}{"token": "weisd"}
		token, err := jwt.NewToken("weisd", claims)
		if err != nil {
			return err
		}
		// show the token use for test
		return c.String(200, "token : %s", token)
	})

```
bat 工具测试 访问成功
```
bat GET http://localhost:1323/jwt Authorization:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0Mzk4ODkzNjAsInRva2VuIjoid2Vpc2QifQ.7pZHBX16DL3qxzsrJYHD2rq1jT_zIOeYVXO96pYJ2C0" client-id:da
GET /jwt HTTP/1.1
Host: localhost:1323
Accept: application/json
Accept-Encoding: gzip, deflate
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0Mzk4ODkzNjAsInRva2VuIjoid2Vpc2QifQ.7pZHBX16DL3qxzsrJYHD2rq1jT_zIOeYVXO96pYJ2C0
Client-Id: da
User-Agent: bat/0.0.3



HTTP/1.1 200 OK
Content-Type : text/plain; charset=utf-8
Date : Fri, 14 Aug 2015 09:53:31 GMT
Content-Length : 61


jwt Access ok with claims map[exp:1.43988936e+09 token:weisd]
```