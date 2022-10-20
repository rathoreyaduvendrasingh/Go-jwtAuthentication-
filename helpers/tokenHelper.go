package helpers

import(
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/yaduvendra/E-commerce/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct{
	Email 		string
	FirstName 	string
	LastName 	string
	Uid			string
	UserType 	string
	jwt.StandardClaims
}


var userCollection *mongo.Collection = database.OpenCollection(database.Client,"gouser")
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string,uid string) (signedToken string, signedRefreshToken string, err error){
	claims := &SignedDetails{
		Email:email,
		FirstName: firstName,
		LastName: lastName,
		Uid: uid,
		UserType:userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(25)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token,err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err
}

func UpdateAllTokens(signedToken,signedrefreshToken,userId string){
	ctx ,cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var updateObj primitive.D
	
	updateObj = append(updateObj, bson.E{"Token",signedToken})
	updateObj = append(updateObj, bson.E{"ResfreshToken",signedrefreshToken})

	Updated_at,_ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj,bson.E{"UpdatedAt",Updated_at})

	upsert := true
	filter := bson.M{"userId":userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set",updateObj},
		},
		&opt,
	)
	defer cancel()
	if err != nil{
		log.Panic(err)
		return
	}
	return
}

func ValidateToken(signedToken string) (claims *SignedDetails,msg string){
	token,err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token)(interface{}, error){
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil{
		msg = err.Error()
		return
	}
	claims,ok := token.Claims.(*SignedDetails)
	if !ok{
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}