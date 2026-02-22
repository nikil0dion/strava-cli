# strava-cli

Minimalist CLI for Strava API v3 built with Go 1.25+

## Features

- Activity details with pace, HR, calories, suffer score
- HR zones breakdown with time percentages
- Best efforts (PRs) tracking
- Gear tracking with total distance
- Splits per kilometer
- Auto token refresh

## Installation

```bash
# Build from source
go build -o strava-cli ./cmd/strava-cli

# Or install to $GOPATH/bin
go install ./cmd/strava-cli

# Move to PATH (optional)
sudo mv strava-cli /usr/local/bin/
```

## Configuration

Credentials stored in `~/.config/strava-cli/credentials.json`:

```json
{
  "client_id": "your_client_id",
  "client_secret": "your_client_secret",
  "access_token": "your_access_token",
  "refresh_token": "your_refresh_token",
  "token_expires_at": "2026-02-22T20:09:00Z",
  "scope": "read,activity:read_all"
}
```

### Getting Tokens

1. Create app at https://www.strava.com/settings/api
2. Authorize with correct scope (see `OAUTH_GUIDE.md`)
3. Save credentials to config file

## Commands

### Activities List
```bash
strava-cli activities --limit 10
```

### Activity Details
```bash
strava-cli activity <activity_id>
```
Output:
```
🏃 Morning Run (Run)
📅 2026-01-15T08:30:00Z

📊 STATS
• Distance: 5.00 km
• Time: 30:00
• Pace: 6:00 /km
• Elevation: +50 m

❤️ HEART RATE
• Avg: 145 bpm
• Max: 170 bpm
• Cadence: 160 spm

🔥 Calories: 350 kcal
💪 Suffer Score: 45

👟 Gear: Running Shoes (200 km total)

🏆 BEST EFFORTS
• 1K: 5:30 (PR #1)
• 1 mile: 9:15
• 5K: 30:00

📏 SPLITS
• Km 1: 6:10 | 140 bpm
• Km 2: 6:00 | 145 bpm
• Km 3: 5:55 | 148 bpm
• Km 4: 6:00 | 146 bpm
• Km 5: 5:55 | 150 bpm
```

### HR Zones
```bash
strava-cli zones <activity_id>
```
Output:
```
❤️ HR ZONES

Zone 1 (0-120 bpm): 2:00 (7%)
Zone 2 (121-150 bpm): 20:00 (67%)
Zone 3 (151-165 bpm): 6:00 (20%)
Zone 4 (166-180 bpm): 2:00 (6%)
Zone 5 (181-max bpm): 0:00 (0%)
```

### Laps/Splits
```bash
strava-cli laps <activity_id>
```

### Stats
```bash
strava-cli stats
```
YTD and recent 4-week totals for running and cycling.

### Profile
```bash
strava-cli profile
```

## Required Scope

For full functionality, authorize with `activity:read_all`:
```
https://www.strava.com/oauth/authorize?client_id=YOUR_ID&response_type=code&redirect_uri=http://localhost&scope=read,activity:read_all&approval_prompt=force
```

## Project Structure

```
strava-cli/
├── cmd/strava-cli/main.go     # CLI entry point
├── internal/
│   ├── api/client.go          # Strava API client
│   ├── auth/refresh.go        # Token refresh
│   └── config/config.go       # Credentials
├── go.mod
├── README.md
└── OAUTH_GUIDE.md
```

## API Limits

- 200 requests / 15 min
- 2,000 requests / day

## License

MIT
