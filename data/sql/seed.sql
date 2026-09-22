-- The initial user is created from APP_ADMIN_EMAIL and APP_ADMIN_PASSWORD.
-- This museum sample is based on the positive game records in:
-- https://github.com/el-acantilado/juegos-peliculas-series/tree/main/data/juegos/positivos

INSERT INTO games (title, platform, release_year, description, created_at)
SELECT 'Commander Keen 4', 'MS-DOS', 1991, 'Mi referencia personal más alta entre los plataformeros.', CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM games WHERE title = 'Commander Keen 4' AND platform = 'MS-DOS');

INSERT INTO games (title, platform, release_year, description, created_at)
SELECT 'Doom', 'MS-DOS', 1993, 'Una experiencia excelente y favorita personal.', CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM games WHERE title = 'Doom' AND platform = 'MS-DOS');

INSERT INTO games (title, platform, release_year, description, created_at)
SELECT 'Quake', 'PC', 1996, 'Un clásico de acción que sigue siendo una experiencia favorita.', CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM games WHERE title = 'Quake' AND platform = 'PC');

INSERT INTO games (title, platform, release_year, description, created_at)
SELECT 'Half-Life', 'PC', 1998, 'Una experiencia narrativa y de acción valorada con 10/10.', CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM games WHERE title = 'Half-Life' AND platform = 'PC');

INSERT INTO games (title, platform, release_year, description, created_at)
SELECT 'Portal', 'PC', 2007, 'Breve y ligero, incluso donde todavía no alcanza el pulido de su secuela.', CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM games WHERE title = 'Portal' AND platform = 'PC');

INSERT INTO games (title, platform, release_year, description, created_at)
SELECT 'Machinarium', 'PC', 2009, 'Un conjunto exquisito de música, sonidos y narrativa sin diálogos.', CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM games WHERE title = 'Machinarium' AND platform = 'PC');

INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/9180/header.jpg', 'Commander Keen 4 — imagen 1', CURRENT_TIMESTAMP FROM games WHERE title = 'Commander Keen 4' AND platform = 'MS-DOS';
INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/9180/library_hero.jpg', 'Commander Keen 4 — imagen 2', CURRENT_TIMESTAMP FROM games WHERE title = 'Commander Keen 4' AND platform = 'MS-DOS';

INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/2280/header.jpg', 'Doom — imagen 1', CURRENT_TIMESTAMP FROM games WHERE title = 'Doom' AND platform = 'MS-DOS';
INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/2280/library_hero.jpg', 'Doom — imagen 2', CURRENT_TIMESTAMP FROM games WHERE title = 'Doom' AND platform = 'MS-DOS';

INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/2310/header.jpg', 'Quake — imagen 1', CURRENT_TIMESTAMP FROM games WHERE title = 'Quake' AND platform = 'PC';
INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/2310/library_hero.jpg', 'Quake — imagen 2', CURRENT_TIMESTAMP FROM games WHERE title = 'Quake' AND platform = 'PC';

INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/70/header.jpg', 'Half-Life — imagen 1', CURRENT_TIMESTAMP FROM games WHERE title = 'Half-Life' AND platform = 'PC';
INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/70/library_hero.jpg', 'Half-Life — imagen 2', CURRENT_TIMESTAMP FROM games WHERE title = 'Half-Life' AND platform = 'PC';

INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/400/header.jpg', 'Portal — imagen 1', CURRENT_TIMESTAMP FROM games WHERE title = 'Portal' AND platform = 'PC';
INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/400/library_hero.jpg', 'Portal — imagen 2', CURRENT_TIMESTAMP FROM games WHERE title = 'Portal' AND platform = 'PC';

INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/40700/header.jpg', 'Machinarium — imagen 1', CURRENT_TIMESTAMP FROM games WHERE title = 'Machinarium' AND platform = 'PC';
INSERT OR IGNORE INTO game_images (game_id, url, alt_text, created_at)
SELECT id, 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/40700/library_hero.jpg', 'Machinarium — imagen 2', CURRENT_TIMESTAMP FROM games WHERE title = 'Machinarium' AND platform = 'PC';
