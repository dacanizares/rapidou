package tst

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// RunUI executes one high-value browser journey. It skips when Chrome is absent.
func RunUI(t *testing.T, factory Factory, credentials Credentials) {
	t.Helper()
	if !browserAvailable() {
		t.Skip("Chrome or Chromium is not installed; API contract still ran")
	}

	app, closeApp := factory(t)
	defer closeApp()
	server := httptest.NewServer(app)
	defer server.Close()

	allocator, cancelAllocator := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", true))...)
	defer cancelAllocator()
	ctx, cancel := chromedp.NewContext(allocator)
	defer cancel()
	ctx, timeout := context.WithTimeout(ctx, 20*time.Second)
	defer timeout()
	uploadPath := filepath.Join(t.TempDir(), "museum.png")
	if err := os.WriteFile(uploadPath, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 0}, 0o600); err != nil {
		t.Fatal(err)
	}

	var museumText string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitVisible("#login-view", chromedp.ByQuery),
		chromedp.SendKeys("#login-email", credentials.Email, chromedp.ByQuery),
		chromedp.SendKeys("#login-password", "wrong-password", chromedp.ByQuery),
		chromedp.Click("#login", chromedp.ByQuery),
		chromedp.Poll(`document.querySelector("#login-error").textContent.includes("invalid email or password")`, nil),
		chromedp.SetValue("#login-password", credentials.Password, chromedp.ByQuery),
		chromedp.Click("#login", chromedp.ByQuery),
		chromedp.WaitVisible("#new-game", chromedp.ByQuery),
		chromedp.Click("#new-game", chromedp.ByQuery),
		chromedp.WaitVisible("#game-dialog", chromedp.ByQuery),
		chromedp.SendKeys("#title", "The Legend of Zelda", chromedp.ByQuery),
		chromedp.SendKeys("#platform", "NES", chromedp.ByQuery),
		chromedp.SendKeys("#release-year", "1986", chromedp.ByQuery),
		chromedp.SendKeys("#description", "An adventure preserved by the museum.", chromedp.ByQuery),
		chromedp.Click("#save", chromedp.ByQuery),
		chromedp.Poll(`document.querySelector("#games").textContent.includes("The Legend of Zelda")`, nil),
		chromedp.WaitVisible("#image-manager", chromedp.ByQuery),
		chromedp.SetUploadFiles("#image-upload-form input[type=file]", []string{uploadPath}, chromedp.ByQuery),
		chromedp.Click("#image-upload-form .btn", chromedp.ByQuery),
		chromedp.Poll(`document.querySelector("#managed-images .source-upload") !== null`, nil),
		chromedp.SendKeys("#image-url", "https://example.com/zelda.jpg", chromedp.ByQuery),
		chromedp.Click("#image-url-form .btn", chromedp.ByQuery),
		chromedp.Poll(`document.querySelector("#managed-images .source-url") !== null`, nil),
		chromedp.Poll(`document.querySelector("#games img[src*='/api/game-images/']") !== null`, nil),
		chromedp.SetValue("#title", "The Legend of Zelda — Museum Edition", chromedp.ByQuery),
		chromedp.Click("#save", chromedp.ByQuery),
		chromedp.Poll(`document.querySelector("#games").textContent.includes("Museum Edition")`, nil),
		chromedp.Click("#close-dialog", chromedp.ByQuery),
		chromedp.Text("#games", &museumText, chromedp.ByQuery),
	); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(museumText, "The Legend of Zelda — Museum Edition") {
		t.Fatal(fmt.Errorf("edited museum piece is not visible; collection contains %q", museumText))
	}
}

func browserAvailable() bool {
	for _, name := range []string{"google-chrome", "chromium", "chromium-browser", "chrome"} {
		if _, err := exec.LookPath(name); err == nil {
			return true
		}
	}
	return false
}
