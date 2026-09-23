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
	browserPath := browserExecutable()
	if browserPath == "" {
		t.Fatal("Chrome or Chromium is required; run the mandatory containerized harness with ./run/test")
	}

	app, closeApp := factory(t)
	defer closeApp()
	server := httptest.NewServer(app)
	defer server.Close()

	allocatorOptions := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browserPath),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	allocator, cancelAllocator := chromedp.NewExecAllocator(context.Background(), allocatorOptions...)
	defer cancelAllocator()
	ctx, cancel := chromedp.NewContext(allocator)
	defer cancel()
	ctx, timeout := context.WithTimeout(ctx, 45*time.Second)
	defer timeout()
	uploadPath := filepath.Join(t.TempDir(), "museum.png")
	if err := os.WriteFile(uploadPath, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 0}, 0o600); err != nil {
		t.Fatal(err)
	}

	var museumText string
	var detailText string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitVisible("#open-login", chromedp.ByQuery),
		chromedp.Click("#open-login", chromedp.ByQuery),
		chromedp.WaitVisible("#login-dialog", chromedp.ByQuery),
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
		chromedp.SendKeys("#steam-url", "https://store.steampowered.com/app/123", chromedp.ByQuery),
		chromedp.SendKeys("#gog-url", "https://www.gog.com/en/game/example", chromedp.ByQuery),
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
		chromedp.WaitNotVisible("#game-dialog", chromedp.ByQuery),
		chromedp.Poll(`document.querySelectorAll("#games .store-link").length === 2`, nil),
		chromedp.Click("#games .artifact", chromedp.ByQuery),
		chromedp.WaitVisible("#detail-dialog", chromedp.ByQuery),
		chromedp.Poll(`document.querySelectorAll("#detail-gallery img").length === 2`, nil),
		chromedp.Poll(`document.querySelectorAll("#detail-stores .store-link").length === 2`, nil),
		chromedp.Text("#detail-dialog", &detailText, chromedp.ByQuery),
		chromedp.Click("#close-detail", chromedp.ByQuery),
		chromedp.Text("#games", &museumText, chromedp.ByQuery),
	); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(museumText, "The Legend of Zelda — Museum Edition") {
		t.Fatal(fmt.Errorf("edited museum piece is not visible; collection contains %q", museumText))
	}
	if !strings.Contains(detailText, "An adventure preserved by the museum.") || !strings.Contains(detailText, "NES · 1986") {
		t.Fatal(fmt.Errorf("museum detail is incomplete; popup contains %q", detailText))
	}
}

func browserExecutable() string {
	for _, name := range []string{"google-chrome", "chromium", "chromium-browser", "chrome"} {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}
