package controllers

import (
	"fmt"
	"go_jwt/initializers"
	"go_jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Loadsignup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

//Signup func

func SignUp(c *gin.Context) {
	// //get the email/pass of req body
	// var body struct {
	// 	Email    string
	// 	Password string
	// }

	// if c.Bind(&body) != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "failed to read body",
	// 	})
	// 	return
	// }
	// Get form data
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	//validate the form data
	if password != confirmPassword {
		// Render the signup template with an error message
		c.HTML(http.StatusOK, "signup.html", gin.H{
			"Error": "Passwords do not match",
		})
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	//create the user
	user := models.User{Username: username, Email: email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
		return
	}

	// Redirect to success page
	c.Redirect(http.StatusSeeOther, "/success")
}

func LoadSuccess(c *gin.Context) {
	c.HTML(http.StatusOK, "success.html", gin.H{})
}

func LoadLogin(c *gin.Context) {
	//session handle
	userID, _ := c.Get("user")
	var id models.ID
	initializers.DB.First(&id, userID)
	fmt.Println(id.Id)
	if id.Id == 0 {
		c.HTML(http.StatusOK, "login.html", gin.H{})

	} else {
		c.Redirect(http.StatusSeeOther, "/home")
	}

}
func Loadhome(c *gin.Context) {

	c.HTML(http.StatusOK, "home.html", gin.H{})
}

//login func

func Login(c *gin.Context) {
	//Get the email and password from the body
	username := c.PostForm("username")
	password := c.PostForm("password")

	//Look up requested user
	var user models.User
	initializers.DB.First(&user, "username = ?", username)

	if user.ID == 0 {
		// Render the signup template with an error message
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error": "Invalid password or username",
		})
		return
	}

	// if user.ID == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "invalid credential",
	// 	})
	// 	return
	// }

	//Compare sent in pass with saved user hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error": "Invalid password or username",
		})
		return
	}

	//Generate a jwt token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	//send it back
	//set cookie

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", " ", false, true)

	//store user

	ID := models.ID{Id: user.ID}
	initializers.DB.Create(&ID)
	//redirect to home page
	c.Redirect(http.StatusSeeOther, "/home")

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func AdminPageHandler(c *gin.Context) {
	var signups []models.User
	initializers.DB.Find(&signups)
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"signups": signups,
	})
}

func DeleteUser(c *gin.Context) {
	//get the post
	//id := c.Param("id")
	// var record models.User
	// if err := initializers.DB.First(&record, id).Error; err != nil {
	// 	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found"})
	// 	return
	// }

	// //Delete the posts
	// if err := initializers.DB.Delete(&record).Error; err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
	// 	return
	// }

	//respond with them
	//
	//get the post
	var user models.User
	id := c.Param("id")
	result := initializers.DB.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	initializers.DB.Delete(&user)
	// c.HTML(http.StatusOK, "admin.html", gin.H{
	// 	"message": "deleted successfully",
	// })

	//redirect to admin page
	c.Redirect(http.StatusSeeOther, "/admin")
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id") // extract the ID of the record to update
	var user models.User
	result := initializers.DB.First(&user, id) // find the record by ID
	if result.Error != nil {
		// handle error
		return
	}
	// update the record with the new data
	user.Username = c.PostForm("username")
	user.Email = c.PostForm("email")

	hash, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	user.Password = string(hash)

	initializers.DB.Save(&user)
	// redirect the user to the updated admin page
	c.Redirect(http.StatusSeeOther, "/admin")
}

func Signout(c *gin.Context) {
	var id models.ID
	userID, _ := c.Get("user")
	fmt.Println(userID)
	result := initializers.DB.First(&id, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	initializers.DB.Delete(&id)
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/login")

}
