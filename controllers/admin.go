package controllers

import (
    "furniture-ecommerce/config"
    "furniture-ecommerce/models"
    "github.com/gin-gonic/gin"
    "net/http"
)

//List all users
func ListUsers(c *gin.Context) {
    var users []models.User
    if result := config.DB.Find(&users); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
        return
    }
    c.JSON(http.StatusOK, users)
}

//  Add a new product
func AddProduct(c *gin.Context) {
    var product models.Product
    if err := c.BindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    if result := config.DB.Create(&product); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Product added"})
}

//  Edit an existing product
func EditProduct(c *gin.Context) {
    var product models.Product
    id := c.Param("id")

    if err := config.DB.First(&product, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    if err := c.BindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    if result := config.DB.Save(&product); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update product"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

//  Delete a product
func DeleteProduct(c *gin.Context) {
    id := c.Param("id")
    if result := config.DB.Delete(&models.Product{}, id); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete product"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

//  Delete a user
func DeleteUser(c *gin.Context) {
    id := c.Param("id")
    if result := config.DB.Delete(&models.User{}, id); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
