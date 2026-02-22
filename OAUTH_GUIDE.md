# OAuth Authorization Guide

## Why You Need This

By default, Strava only gives `read` scope (athlete profile).
To read activities, HR zones, best efforts, gear — you need `activity:read_all` scope.

## How to Get Full Access

### 1. Build OAuth URL

```
https://www.strava.com/oauth/authorize?client_id=YOUR_CLIENT_ID&response_type=code&redirect_uri=http://localhost&approval_prompt=force&scope=read,activity:read_all
```

**Parameters:**
- `client_id` - your Client ID (from strava.com/settings/api)
- `redirect_uri` - must match your app settings (use `http://localhost`)
- `scope` - requested permissions: `read,activity:read_all`
- `approval_prompt=force` - always show consent screen

### 2. Open in Browser

Authorize and browser redirects to:
```
http://localhost/?state=&code=AUTHORIZATION_CODE&scope=read,activity:read_all
```

Copy the `code` value from URL.

### 3. Exchange Code for Tokens

```bash
curl -X POST https://www.strava.com/oauth/token \
  -d client_id=YOUR_CLIENT_ID \
  -d client_secret=YOUR_CLIENT_SECRET \
  -d code=AUTHORIZATION_CODE \
  -d grant_type=authorization_code
```

**Response:**
```json
{
  "token_type": "Bearer",
  "expires_at": 1629387600,
  "expires_in": 21600,
  "refresh_token": "your_refresh_token",
  "access_token": "your_access_token",
  "athlete": {...}
}
```

### 4. Update credentials.json

```bash
mkdir -p ~/.config/strava-cli
nano ~/.config/strava-cli/credentials.json
```

```json
{
  "client_id": "YOUR_CLIENT_ID",
  "client_secret": "YOUR_CLIENT_SECRET",
  "access_token": "your_access_token",
  "refresh_token": "your_refresh_token",
  "token_expires_at": "2026-02-22T20:00:00Z",
  "scope": "read,activity:read_all"
}
```

Note: `token_expires_at` — convert `expires_at` unix timestamp to ISO format.

## Available Scopes

| Scope | Description |
|-------|-------------|
| `read` | Public profile |
| `read_all` | Full profile (including private) |
| `activity:read` | Read public activities |
| `activity:read_all` | Read all activities (including private) ✅ |
| `activity:write` | Create/update activities |

**Recommended:** `read,activity:read_all`

## Token Refresh

Access tokens expire every 6 hours. The CLI automatically refreshes using `refresh_token`.

No manual action needed after initial setup.

## Troubleshooting

**Error: "Authorization Error - activity:read_permission"**
- Your token has wrong scope
- Re-run OAuth flow with `scope=read,activity:read_all`

**Error: "401 Unauthorized"**
- Token expired and refresh failed
- Check `client_secret` is correct
- Re-run OAuth flow
