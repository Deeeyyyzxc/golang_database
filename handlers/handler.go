package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"sample/db"

	"github.com/gofiber/fiber/v2"
)

// RegisterUser handles user registration
func RegisterUser(c *fiber.Ctx) error {
	// Get form data
	name := c.FormValue("name")
	email := c.FormValue("email")
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image is required"})
	}

	// Ensure the images directory exists
	imageDir := filepath.Join("Documents", "user", "images", email)
	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating image directory"})
	}

	// Save the uploaded image
	imagePath := filepath.Join(imageDir, fmt.Sprintf("%s.jpg", email))
	if err := c.SaveFile(file, imagePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving user image"})
	}

	// Create user with image path
	user := db.User{
		Name:      name,
		Email:     email,
		ImageData: []byte{}, // Change this to save image data if necessary
	}

	// Save the user to the database
	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving user to database"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully!"})
}

// CallPythonFaceRecognition runs the Python script and returns recognized names
func CallPythonFaceRecognition(email string, imagePath string) ([]string, error) {
	scriptPath := os.Getenv("PYTHON_SCRIPT_PATH")
	if scriptPath == "" {
		scriptPath = "/Users/Andrey_Delos_Reyes/Documents/De Los Reyes, Andrey/Face-Recognition/main.py"
	}

	cmd := exec.Command("python3", scriptPath, email, imagePath) // Pass email and imagePath
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing python script:", err)
		fmt.Println("Output:", out.String())
		return nil, fmt.Errorf("error executing python script: %w", err)
	}

	// Parse the output from the Python script
	var output map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &output); err != nil {
		fmt.Println("Error parsing JSON output:", err)
		fmt.Println("Raw output:", out.String())
		return nil, fmt.Errorf("error parsing JSON output: %w", err)
	}

	// Extract recognized names
	recognizedNames, _ := output["recognized_names"].([]interface{})
	names := make([]string, len(recognizedNames))
	for i, name := range recognizedNames {
		names[i] = name.(string)
	}

	return names, nil
}

// CheckIn handles user check-in
func CheckIn(c *fiber.Ctx) error {
	email := c.FormValue("email")
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image is required"})
	}

	imageDir := filepath.Join("Documents", "user", "images", email)
	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating image directory"})
	}

	checkinImagePath := filepath.Join(imageDir, fmt.Sprintf("%s-checkin.jpg", email))
	if err := c.SaveFile(file, checkinImagePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving check-in image"})
	}

	// Call the Python face recognition
	recognizedNames, err := CallPythonFaceRecognition(email, checkinImagePath) // Pass email and check-in image path
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error during face recognition"})
	}

	// Process recognized names
	if len(recognizedNames) > 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Check-in successful!", "recognized_names": recognizedNames})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Face not recognized"})
	}
}
