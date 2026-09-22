# Rapidou base behavior

Rapidou is currently a small authenticated user directory and a foundation for the actual product domain.

- `GET /api/health` is public.
- A seeded user can log in with email and password.
- Authentication works through an HttpOnly cookie for the browser or a Bearer JWT for API clients.
- An authenticated user can list, create, edit, and delete users.
- A user cannot delete their own active account.
- Names, valid unique emails, and passwords of at least eight characters are required when creating users.
- The frontend is embedded in the executable and uses plain HTML, CSS, and JavaScript.

Startup requires `APP_JWT_SECRET`. Set both `APP_ADMIN_EMAIL` and `APP_ADMIN_PASSWORD` to create the first user in an empty database. They are ignored once a user exists.

