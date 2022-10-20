package helpers

import(
	"errors"
	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (e error){
	userType := c.GetString("userType")
	e = nil
	if userType != role{
		e = errors.New("Unauthorized to acces this resource")
		return e
	}
	return e
}

func MatchUserTypeToUid(c *gin.Context,userId string) (e error){
	userType := c.GetString("userType")
	uid := c.GetString("uid")
	e = nil
	if userType == "USER" && uid != userId {
		e = errors.New("Unauthorized to acces this resource")
		return e
	}

	e = CheckUserType(c, userType)
	return e;
}