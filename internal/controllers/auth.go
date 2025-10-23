package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "signup ok"})
}

func Login(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "login ok"})
}
