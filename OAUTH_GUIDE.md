# OAuth Authorization Guide

## Why You Need This

By default, Strava only gives `read` scope (athlete profile).  
To read activities you need `activity:read` scope.

## How to Get Full Access

### 1. Build OAuth URL

```
https://www.strava.com/oauth/authorize?client_id=YOUR_CLIENT_ID&response_type=code&redirect_uri=http://localhost/callback&approval_prompt=force&scope=read,activity:read
```

**Parameters:**
- `client_id` - your Client ID (from strava.com/settings/api)
- `redirect_uri` - must match your app settings
- `scope` - requested permissions (comma-separated)

### 2. Open in Browser

Browser redirects to:
```
http://localhost/callback?state=&code=AUTHORIZATION_CODE&scope=read,activity:read
```

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
  "refresh_token": "new_refresh_token",
  "access_token": "new_access_token",
  "athlete": {...}
}
```

### 4. Update credentials.json

```bash
nano ~/.config/strava-cli/credentials.json
```

Replace `access_token`, `refresh_token`, `token_expires_at`, `scope`.

## Available Scopes

- `read` - public profile
- `read_all` - full profile (including private data)
- `profile:read_all` - detailed profile
- `profile:write` - update profile
- `activity:read` - read activities ✅
- `activity:read_all` - read all activities (including private)
- `activity:write` - create/update activities

## Quick Fix for `activity:read`

If you don't want to run the full OAuth flow:

1. Go to https://www.strava.com/settings/api
2. In "Your Access Token" section, click "Regenerate"
3. **Important:** Strava will generate a new token, but scope will remain `read`
4. For full access you still need OAuth flow

## Automatic OAuth (TODO)

Planned to be built into CLI:
```bash
strava-cli auth --scope read,activity:read
# Opens browser, gets code, exchanges for tokens, saves
```

For now, follow the manual instructions above.
