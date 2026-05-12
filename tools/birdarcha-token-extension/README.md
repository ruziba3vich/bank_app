# Birdarcha Token Refresher

Chrome/Edge extension that periodically reads the birdarcha JWT from
`office.birdarcha.uz` localStorage (`access-token` key) and PUTs it to the
bank-app `/api/v1/entrepreneurs/birdarcha-token` endpoint, authenticated with
an `X-Api-Key` header.

## Install (dev / unpacked)

1. Open `chrome://extensions`
2. Enable **Developer mode** (top right)
3. Click **Load unpacked** → select this directory
4. Click the extension's icon → **Options**
5. Fill in:
   - **Bank-app base URL** — e.g. `https://bank-back.shoha-coder.uz`
   - **X-Api-Key** — value of `BIRDARCHA_API_KEY` from the server's `.env`
   - **Refresh interval (minutes)** — 180 (= 3h) is a sensible default since
     the JWT TTL is ~4h
6. Click **Save**, then **Refresh now** to test

## How it works

- A `chrome.alarms` timer fires every N minutes.
- The background worker finds (or opens) a tab on `office.birdarcha.uz`.
- It reads `localStorage["access-token"]` via `chrome.scripting.executeScript`.
- If the token is expired (decoded `exp` claim), it reloads the tab so birdarcha
  silently re-issues a fresh JWT, then reads it again.
- It PUTs `{ "token": "..." }` with the `X-Api-Key` header to the bank-app.

## Requirements

- You stay logged into birdarcha in this browser profile. If the underlying
  egov SSO session expires, you must re-login manually (E-IMZO).
- The browser must be running.

## Status

The popup and options page show the last run's status (success or error).
