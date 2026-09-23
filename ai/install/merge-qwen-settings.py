#!/usr/bin/env python3
"""Merge Rapidou's local Ollama provider into Qwen Code user settings."""

import json
import os
import pathlib
import shutil
import sys
import tempfile


MARKER = "Rapidou local Ollama model"


def fail(message: str) -> None:
    print(f"error: {message}", file=sys.stderr)
    raise SystemExit(2)


def main() -> None:
    if len(sys.argv) != 4:
        fail("usage: merge-qwen-settings.py <settings.json> <model> <context-window>")

    settings_path = pathlib.Path(sys.argv[1])
    model = sys.argv[2]
    try:
        context_window = int(sys.argv[3])
    except ValueError:
        fail("context window must be an integer")

    settings = {}
    if settings_path.exists():
        try:
            settings = json.loads(settings_path.read_text())
        except (OSError, json.JSONDecodeError) as error:
            fail(f"cannot read valid JSON from {settings_path}: {error}")
        if not isinstance(settings, dict):
            fail(f"{settings_path} must contain a JSON object")

        backup = settings_path.with_name(settings_path.name + ".rapidou-backup")
        if not backup.exists():
            shutil.copy2(settings_path, backup)

    env = settings.setdefault("env", {})
    providers = settings.setdefault("modelProviders", {})
    security = settings.setdefault("security", {})
    auth = security.setdefault("auth", {})
    model_settings = settings.setdefault("model", {})
    if not all(isinstance(value, dict) for value in (env, providers, security, auth, model_settings)):
        fail("env, modelProviders, security.auth, and model settings must be JSON objects")

    openai_models = providers.setdefault("openai", [])
    if not isinstance(openai_models, list):
        fail("modelProviders.openai must be a JSON array")

    openai_models[:] = [
        entry
        for entry in openai_models
        if not (isinstance(entry, dict) and entry.get("description") == MARKER)
    ]
    openai_models.append(
        {
            "id": model,
            "name": f"{model} (local Ollama)",
            "description": MARKER,
            "envKey": "RAPIDOU_OLLAMA_API_KEY",
            "baseUrl": "http://127.0.0.1:11434/v1",
            "generationConfig": {
                "timeout": 300000,
                "streamIdleTimeoutMs": 600000,
                "maxRetries": 1,
                "contextWindowSize": context_window,
            },
        }
    )

    env["RAPIDOU_OLLAMA_API_KEY"] = "ollama"
    auth["selectedType"] = "openai"
    model_settings["name"] = model

    settings_path.parent.mkdir(parents=True, exist_ok=True)
    descriptor, temporary_name = tempfile.mkstemp(
        prefix="settings.json.", dir=settings_path.parent
    )
    try:
        with os.fdopen(descriptor, "w") as temporary:
            json.dump(settings, temporary, indent=2, ensure_ascii=False)
            temporary.write("\n")
        os.chmod(temporary_name, 0o600)
        os.replace(temporary_name, settings_path)
    finally:
        if os.path.exists(temporary_name):
            os.unlink(temporary_name)


if __name__ == "__main__":
    main()
