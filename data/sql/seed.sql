-- The initial user is created from APP_ADMIN_EMAIL and APP_ADMIN_PASSWORD.
-- Keeping password hashes out of this file makes the seed safe to inspect.

INSERT INTO games (title, platform, release_year, description, created_at)
SELECT 'Pong', 'Arcade', 1972, 'Una de las primeras piezas que llevó el videojuego al público general.', CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM games WHERE title = 'Pong' AND platform = 'Arcade');

INSERT INTO games (title, platform, release_year, description, created_at)
SELECT 'Super Mario Bros.', 'NES', 1985, 'Una referencia fundamental del diseño de plataformas.', CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM games WHERE title = 'Super Mario Bros.' AND platform = 'NES');

INSERT INTO games (title, platform, release_year, description, created_at)
SELECT 'Doom', 'PC', 1993, 'Una pieza decisiva en la historia de los juegos de acción en primera persona.', CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM games WHERE title = 'Doom' AND platform = 'PC');
