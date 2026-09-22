# Museo Pixel

Rapidou includes a small video game museum as its complete example application.

- `GET /api/health` is public.
- A seeded user can log in with email and password.
- Authentication works through an HttpOnly cookie for the browser or a Bearer JWT for API clients.
- An authenticated user can list, create, edit, and delete users.
- A user cannot delete their own active account.
- Names, valid unique emails, and passwords of at least eight characters are required when creating users.
- Anyone can browse the museum collection.
- An authenticated curator can add, edit, and remove museum pieces.
- Adding and editing use one responsive modal; catalog cards remain focused on viewing.
- Each piece has a title, platform, release year, and optional description.
- Each piece may have multiple JPEG, PNG, GIF, or WebP images.
- Curators can upload images up to 5 MB; uploaded files are stored in SQLite.
- Curators can also attach HTTPS image URLs. The editor identifies URL images and uploaded files separately.
- The development sample contains six favorites from the referenced positive-games collection and two remote image URLs per game.
- Release years before 1950 or beyond next year are rejected.
- The frontend is embedded in the executable and uses plain HTML, CSS, and JavaScript.

Startup requires `APP_JWT_SECRET`. Set both `APP_ADMIN_EMAIL` and `APP_ADMIN_PASSWORD` to create the first user in an empty database. They are ignored once a user exists.
