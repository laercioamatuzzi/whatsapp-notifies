package utils

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/joho/godotenv"
)

const (
	DEBUG            = false
	WHATSAPP_DB_NAME = "./whatsapp.db"
	NOTIFIES_DB_NAME = "./notifies.db"
)

func LoadEnv() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func GetThumbnail(filePath string) []byte {

	f, err := imaging.Open(filePath)
	if err != nil {
		return nil
	}
	thumbnail := imaging.Thumbnail(f, 100, 100, imaging.Lanczos)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, thumbnail, nil)

	if err != nil {
		return nil
	}

	return buf.Bytes()
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func Logging(text string) {

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(os.Getenv("WHATSAPP_NOTIFIES_CONFIG_PATH")+"log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open log file.", err)
		return
	}
	defer file.Close()
	log.SetOutput(file)

	log.Println(text)
}

func LoggingInfo(text string) {

	Logging(fmt.Sprintf("[INFO] %s", text))
}

func LoggingDebug(text string) {

	if os.Getenv("DEBUG") == "true" {
		Logging(fmt.Sprintf("[DEBUG] %s", text))
	}
}

func LoggingError(text string) {

	Logging(fmt.Sprintf("[ERROR] %s", text))
}

func ValidLineFormat(number string) bool {

	// Brasil Celular
	if strings.HasPrefix(number, "55") && len(number) == 13 {
		return true
	}

	// Brasil Fixo
	if strings.HasPrefix(number, "55") && len(number) == 12 {

		if number[4:5] == "2" || number[4:5] == "3" || number[4:5] == "4" || number[4:5] == "5" {
			return true
		}
		return false
	}

	// Mexico
	if strings.HasPrefix(number, "52") && len(number) == 13 {
		return true
	}

	// Portugal
	if strings.HasPrefix(number, "351") && len(number) == 12 {
		return true
	}

	return false
}

func AddNinthdigit(number string) string {

	if strings.HasPrefix(number, "55") && len(number) == 12 {

		if number[4:5] != "2" && number[4:5] != "3" && number[4:5] != "4" && number[4:5] != "5" {
			number = fmt.Sprintf("%s9%s", number[:4], number[4:])
		}
	}

	return number
}
