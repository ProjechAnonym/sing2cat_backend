package middlewarefunc

import (
	"fmt"
	"net/http"
	ginconfig "sing2cat_web/GinConfig"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Token string `json:"token"`
	jwt.StandardClaims
}

func Generate_token() (string,error){
	key,err := ginconfig.Get_value("config","jwt","key")
	if err != nil{
		ginconfig.Logger_caller("Generate JWT failed!",err)
		return "",err
	}
	expire_time := time.Now().Add(30 * 24 * time.Hour)
	claims := Claims{Token: key.(string),StandardClaims: jwt.StandardClaims{Issuer: "SIFULIN",ExpiresAt: expire_time.Unix()}}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("linsifu"))
	if err != nil{
		ginconfig.Logger_caller("Generate JWT failed!",err)
		return "",err
	}
	return token,nil
}	

func Verify_token(token string) (*Claims,error){
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("linsifu"), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		fmt.Println(tokenClaims)
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func Jwt_auth() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == ""{
			ctx.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized"})
			ctx.Abort()
			return
		}
		claims,err := Verify_token(header)
		if err != nil{
			ctx.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Set("token",claims.Token)
	}
}