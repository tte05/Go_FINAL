package tests

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestRegistrationEndToEnd(t *testing.T) {
	caps := selenium.Capabilities{"browserName": "chrome"}

	service, err := selenium.NewChromeDriverService("app/Resourses/chromedriver-win64/chromedriver-win64/chromedriver.exe", 9515)
	if err != nil {
		t.Fatal("Error starting the WebDriver service:", err)
	}
	defer service.Stop()

	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		t.Fatal("Error connecting to WebDriver:", err)
	}
	defer wd.Quit()

	err = wd.Get("http://localhost:8080/register")
	if err != nil {
		t.Fatal("Error navigating to registration page:", err)
	}

	inputs := map[string]string{
		"username": "tamer",
		"email":    "tamertazhenov2005@gmail.com",
		"password": "tamer2005",
	}

	for id, value := range inputs {
		element, err := wd.FindElement(selenium.ByID, id)
		if err != nil {
			t.Fatalf("Error finding %s input field: %v", id, err)
		}
		err = element.SendKeys(value)
		if err != nil {
			t.Fatalf("Error entering %s: %v", id, err)
		}
	}

	registerButton, err := wd.FindElement(selenium.ByID, "register-button")
	if err != nil {
		t.Fatal("Error finding register button:", err)
	}
	err = registerButton.Click()
	if err != nil {
		t.Fatal("Error clicking register button:", err)
	}

	time.Sleep(2 * time.Second)

	currentURL, err := wd.CurrentURL()
	if err != nil {
		t.Fatal("Error getting current URL:", err)
	}
	if currentURL != "http://localhost:8080/login" {
		t.Error("Expected to be redirected to login page after registration, but got:", currentURL)
	}
}
