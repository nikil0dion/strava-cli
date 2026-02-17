# Integration Ideas

## Strava + Garmin Framework

Two data sources, different strengths:

### Garmin (gccli)
- **Strengths:** Health, physiology, sleep
- **Metrics:** Body Battery, HRV, RHR, Sleep Score, Stress, VO2 Max
- **Frequency:** Daily (automatic measurements)
- **Use case:** Medical monitoring, recovery, training readiness

### Strava (strava-cli)
- **Strengths:** Workouts, activities, segments
- **Metrics:** Distance, pace, heart rate (during activity), elevation, splits
- **Frequency:** After each workout
- **Use case:** Workout analysis, progress tracking, KOM hunting

## Possible Scenarios

### 1. Recovery vs Load Correlation
```bash
# Garmin: Body Battery low
# Strava: Intense run scheduled
# → Recommendation: light recovery workout
```

### 2. Automatic Post-Workout Report
```bash
# Trigger: new activity in Strava (webhook)
# Read: average HR, pace, zones
# Read: Garmin data for the day (HRV, stress, last night's sleep)
# Generate: report "How workout affected recovery"
```

### 3. Weekly Summary
```bash
# Strava: mileage, time in zones, number of workouts
# Garmin: average Body Battery, sleep quality, RHR trend
# Result: "You ran 50 km, but slept poorly → more rest next week"
```

### 4. Training Planner
```bash
# Input: target distance (marathon/half-marathon)
# Garmin: current recovery level
# Strava: recent workouts and progress
# Output: adaptive plan accounting for physiology
```

## Technical Implementation

### Unified Data Layer
```
memory/workouts/
  2026-02-17-morning-run.md  # Strava activity + Garmin health
```

Entry format:
```markdown
# Morning Run - 2026-02-17

## Strava
- Distance: 10.2 km
- Time: 54:23
- Avg Pace: 5:20/km
- Avg HR: 152 bpm
- Max HR: 168 bpm

## Garmin (same day)
- Sleep Score: 78 (7.4h)
- Body Battery: 65 → 33 (after workout)
- HRV: 45 ms
- Stress: 38 → 72 (spike during run)

## Analysis
Workout in zone 3-4, heart rate higher than usual.  
Body Battery dropped significantly - possible under-recovery from yesterday.  
Recommendation: light base run or rest tomorrow.
```

### Script Example
```bash
#!/bin/bash
# scripts/workout_report.sh

# Get latest activity from Strava
ACTIVITY=$(strava-cli activities --limit 1 --json)

# Get Garmin data for today
HEALTH=$(gccli health today --json)

# Combine and send to AI for analysis
# → sessions_send + prompt with both JSONs
# → AI writes to memory/workouts/YYYY-MM-DD-activity.md
```

## Roadmap

- [ ] Strava webhook listener (auto-update on new activity)
- [ ] Unified JSON export (strava + garmin in one file)
- [ ] Dashboard (HTML + charts for visualizing both sources)
- [ ] Predictive analytics (ML model on historical data)
- [ ] Integration with clawd-coach for automatic training plan adjustments

## See Also

- [garmin-cli skill](/home/node/.openclaw/workspace/skills/garmin-cli/SKILL.md)
- [coach skill](/home/node/.openclaw/workspace/skills/clawd-coach/SKILL.md)
- [Strava API docs](https://developers.strava.com/docs/reference/)
