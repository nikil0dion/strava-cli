# strava-cli

Minimalist CLI for Strava API v3 built with Go 1.26

## Installation

```bash
# Build from source
cd strava-cli
go build -o strava-cli ./cmd/strava-cli

# Or install to $GOPATH/bin
go install ./cmd/strava-cli

# Move to PATH (optional)
sudo mv strava-cli /usr/local/bin/
```

## Configuration

Credentials are stored in `~/.config/strava-cli/credentials.json`:

```json
{
  "client_id": "your_client_id",
  "client_secret": "your_client_secret",
  "access_token": "your_access_token",
  "refresh_token": "your_refresh_token",
  "token_expires_at": "2026-02-17T11:47:16Z",
  "scope": "read"
}
```

### Getting Tokens

1. Create an application at https://www.strava.com/settings/api
2. Authorize via OAuth (or use the Access Token from the app page)
3. Save credentials to `~/.config/strava-cli/credentials.json`

### Scopes

- `read` - athlete profile, basic info ✅
- `activity:read` - read activities (runs, rides) ❌ (not yet)
- `activity:write` - upload activities ❌ (not yet)

For `activities` command you need `activity:read` scope.  
To re-authorize with a new scope - run OAuth flow again (see `OAUTH_GUIDE.md`).

## Commands

### Profile
```bash
strava-cli profile
```
Shows athlete profile: name, location, weight.

### Activities
```bash
strava-cli activities --limit 10
```
List recent activities (requires `activity:read` scope).

### Stats
```bash
strava-cli stats
```
Athlete statistics for year-to-date and recent 4 weeks (run + ride).

## Automatic Token Refresh

CLI automatically refreshes `access_token` using `refresh_token` when it expires (every 6 hours).

## Project Structure

```
strava-cli/
├── cmd/strava-cli/        # Entry point
│   └── main.go
├── internal/
│   ├── api/               # API client
│   │   └── client.go
│   ├── auth/              # OAuth and refresh
│   │   └── refresh.go
│   └── config/            # Credentials handling
│       └── config.go
├── go.mod
└── README.md
```

## TODO

- [ ] Built-in OAuth flow (currently requires manual token)
- [ ] Support `activity:read` scope
- [ ] Detailed activity info
- [ ] Export to JSON/CSV
- [ ] Webhook subscription for new activities
- [ ] Integration with Garmin via unified framework
# strava-cli
