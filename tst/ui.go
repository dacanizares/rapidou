package tst

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os/exec"
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

	var usersText string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitVisible("#login-view", chromedp.ByQuery),
		chromedp.SendKeys("#login-email", credentials.Email, chromedp.ByQuery),
		chromedp.SendKeys("#login-password", "wrong-password", chromedp.ByQuery),
		chromedp.Click("#login", chromedp.ByQuery),
		chromedp.Poll(`document.querySelector("#login-error").textContent.includes("invalid email or password")`, nil),
		chromedp.SetValue("#login-password", credentials.Password, chromedp.ByQuery),
		chromedp.Click("#login", chromedp.ByQuery),
		chromedp.WaitVisible("#app-view", chromedp.ByQuery),
		chromedp.SendKeys("#name", "Browser User", chromedp.ByQuery),
		chromedp.SendKeys("#email", "browser@rapidou.test", chromedp.ByQuery),
		chromedp.SendKeys("#password", "browser-password", chromedp.ByQuery),
		chromedp.Click("#save", chromedp.ByQuery),
		chromedp.Poll(`document.querySelector("#users").textContent.includes("Browser User")`, nil),
		chromedp.Text("#users", &usersText, chromedp.ByQuery),
	); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(usersText, "Browser User") {
		t.Fatal(fmt.Errorf("created user is not visible; users contain %q", usersText))
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
