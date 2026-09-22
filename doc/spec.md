# Museo Pixel

Rapidou includes a small video game museum as its complete example application.

- `GET /api/health` is public.
- A seeded user can log in with email and password.
- Authentication works through an HttpOnly cookie for the browser or a Bearer JWT for API clients.
- Browser login opens in a centered responsive popup instead of occupying space in the public catalog.
- Popup backdrops transition from transparent to dark over 500 ms. Clicking the backdrop or the close control dismisses the popup.
- An authenticated user can list, create, edit, and delete users.
- A user cannot delete their own active account.
- Names, valid unique emails, and passwords of at least eight characters are required when creating users.
- Anyone can browse the museum collection.
- An authenticated curator can add, edit, and remove museum pieces.
- Adding and editing use one compact responsive popup; catalog cards remain focused on viewing and expose only small curator actions.
- Each piece has a title, platform, release year, optional description, optional Steam URL, and optional GOG URL.
- Store URLs must use HTTPS and their corresponding official domain. Public cards show a Steam or GOG button only when that link exists.
- Each piece may have multiple JPEG, PNG, GIF, or WebP images.
- Curators can upload images up to 5 MB; uploaded files are stored in SQLite.
- Curators can also attach HTTPS image URLs. The editor identifies URL images and uploaded files separately.
- The development sample contains six favorites from the referenced positive-games collection, two remote image URLs per game, Steam links for all six, and GOG links where a matching catalog entry exists.
- Existing SQLite databases are migrated in place to add the store-link fields without discarding museum data.
- Release years before 1950 or beyond next year are rejected.
- The frontend is embedded in the executable and uses plain HTML, CSS, and JavaScript.

## Visual language

- The public page uses a warm neutral palette and spacious editorial layout inspired by Anthropic's visual language.
- The hero uses `https://el-acantilado.com/images/top.jpg` as its background and has a straight lower edge without rounded corners.
- Share Tech is the display face for titles, subtitles, section labels, years, counters, and other catalog metadata.
- Space Grotesk is the body and interface face for paragraphs, forms, buttons, and controls.
- Both font families are bundled with the application and must not depend on a third-party font service.
- The catalog and both popups remain usable on narrow mobile viewports without horizontal scrolling.

Startup requires `APP_JWT_SECRET`. Set both `APP_ADMIN_EMAIL` and `APP_ADMIN_PASSWORD` to create the first user in an empty database. They are ignored once a user exists.
