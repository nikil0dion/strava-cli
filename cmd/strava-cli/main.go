package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/nikilodion/strava-cli/internal/api"
	"github.com/nikilodion/strava-cli/internal/auth"
	"github.com/nikilodion/strava-cli/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "profile":
		cmdProfile()
	case "activities":
		cmdActivities()
	case "activity":
		cmdActivity()
	case "zones":
		cmdZones()
	case "laps":
		cmdLaps()
	case "stats":
		cmdStats()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func getClient() *api.Client {
	creds, err := config.LoadCredentials()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading credentials: %v\n", err)
		os.Exit(1)
	}

	if err := auth.EnsureValidToken(creds); err != nil {
		fmt.Fprintf(os.Stderr, "Error refreshing token: %v\n", err)
		os.Exit(1)
	}

	return api.NewClient(creds.AccessToken)
}

func cmdProfile() {
	client := getClient()

	athlete, err := client.GetAthlete()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Athlete Profile\n")
	fmt.Printf("---------------\n")
	fmt.Printf("ID:       %d\n", athlete.ID)
	fmt.Printf("Name:     %s %s\n", athlete.FirstName, athlete.LastName)
	fmt.Printf("Username: %s\n", athlete.Username)
	fmt.Printf("Location: %s, %s, %s\n", athlete.City, athlete.State, athlete.Country)
	if athlete.Weight > 0 {
		fmt.Printf("Weight:   %.1f kg\n", athlete.Weight)
	}
}

func cmdActivities() {
	fs := flag.NewFlagSet("activities", flag.ExitOnError)
	limit := fs.Int("limit", 10, "Number of activities to show")
	fs.Parse(os.Args[2:])

	client := getClient()

	activities, err := client.GetActivities(*limit, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching activities: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Recent Activities (last %d)\n", *limit)
	fmt.Printf("---------------------------\n\n")

	for i, act := range activities {
		fmt.Printf("%d. %s\n", i+1, act.Name)
		fmt.Printf("   Type:     %s\n", act.Type)
		fmt.Printf("   Date:     %s\n", act.StartDate)
		fmt.Printf("   Distance: %.2f km\n", act.Distance/1000)
		fmt.Printf("   Time:     %d min\n", act.MovingTime/60)
		if act.AverageHeartrate > 0 {
			fmt.Printf("   Avg HR:   %.0f bpm\n", act.AverageHeartrate)
		}
		fmt.Println()
	}
}

func cmdStats() {
	client := getClient()

	athlete, err := client.GetAthlete()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching athlete: %v\n", err)
		os.Exit(1)
	}

	stats, err := client.GetStats(athlete.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching stats: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Athlete Statistics\n")
	fmt.Printf("------------------\n\n")

	fmt.Printf("Year to Date - Running\n")
	fmt.Printf("  Activities: %d\n", stats.YTDRunTotals.Count)
	fmt.Printf("  Distance:   %.2f km\n", stats.YTDRunTotals.Distance/1000)
	fmt.Printf("  Time:       %.1f hours\n", stats.YTDRunTotals.MovingTime/3600)
	fmt.Printf("  Elevation:  %.0f m\n\n", stats.YTDRunTotals.ElevationGain)

	fmt.Printf("Year to Date - Cycling\n")
	fmt.Printf("  Activities: %d\n", stats.YTDRideTotals.Count)
	fmt.Printf("  Distance:   %.2f km\n", stats.YTDRideTotals.Distance/1000)
	fmt.Printf("  Time:       %.1f hours\n", stats.YTDRideTotals.MovingTime/3600)
	fmt.Printf("  Elevation:  %.0f m\n\n", stats.YTDRideTotals.ElevationGain)

	fmt.Printf("Recent (4 weeks) - Running\n")
	fmt.Printf("  Activities: %d\n", stats.RecentRunTotals.Count)
	fmt.Printf("  Distance:   %.2f km\n", stats.RecentRunTotals.Distance/1000)
	fmt.Printf("  Time:       %.1f hours\n\n", stats.RecentRunTotals.MovingTime/3600)

	fmt.Printf("Recent (4 weeks) - Cycling\n")
	fmt.Printf("  Activities: %d\n", stats.RecentRideTotals.Count)
	fmt.Printf("  Distance:   %.2f km\n", stats.RecentRideTotals.Distance/1000)
	fmt.Printf("  Time:       %.1f hours\n", stats.RecentRideTotals.MovingTime/3600)
}

func cmdActivity() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: strava-cli activity <activity_id>")
		os.Exit(1)
	}

	activityID, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid activity ID: %s\n", os.Args[2])
		os.Exit(1)
	}

	client := getClient()

	activity, err := client.GetActivity(activityID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching activity: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🏃 %s (%s)\n", activity.Name, activity.Type)
	fmt.Printf("📅 %s\n\n", activity.StartDateLocal)

	fmt.Printf("📊 STATS\n")
	fmt.Printf("• Distance: %.2f km\n", activity.Distance/1000)
	fmt.Printf("• Time: %d:%02d\n", activity.MovingTime/60, activity.MovingTime%60)
	if activity.AverageSpeed > 0 && activity.Type == "Run" {
		pace := 1000 / activity.AverageSpeed / 60
		paceMin := int(pace)
		paceSec := int((pace - float64(paceMin)) * 60)
		fmt.Printf("• Pace: %d:%02d /km\n", paceMin, paceSec)
	}
	if activity.TotalElevationGain > 0 {
		fmt.Printf("• Elevation: +%.0f m\n", activity.TotalElevationGain)
	}

	if activity.AverageHeartrate > 0 {
		fmt.Printf("\n❤️ HEART RATE\n")
		fmt.Printf("• Avg: %.0f bpm\n", activity.AverageHeartrate)
		fmt.Printf("• Max: %.0f bpm\n", activity.MaxHeartrate)
	}
	if activity.AverageCadence > 0 {
		fmt.Printf("• Cadence: %.0f spm\n", activity.AverageCadence*2)
	}

	if activity.Calories > 0 || activity.SufferScore > 0 {
		fmt.Println()
		if activity.Calories > 0 {
			fmt.Printf("🔥 Calories: %.0f kcal\n", activity.Calories)
		}
		if activity.SufferScore > 0 {
			fmt.Printf("💪 Suffer Score: %.0f\n", activity.SufferScore)
		}
	}

	// Gear
	if activity.Gear != nil && activity.Gear.Name != "" {
		fmt.Printf("\n👟 Gear: %s (%.0f km total)\n", activity.Gear.Name, activity.Gear.Distance/1000)
	}

	// Best Efforts
	if len(activity.BestEfforts) > 0 {
		fmt.Printf("\n🏆 BEST EFFORTS\n")
		for _, be := range activity.BestEfforts {
			mins := be.MovingTime / 60
			secs := be.MovingTime % 60
			prText := ""
			if be.PRRank != nil && *be.PRRank > 0 {
				prText = fmt.Sprintf(" (PR #%d)", *be.PRRank)
			}
			fmt.Printf("• %s: %d:%02d%s\n", be.Name, mins, secs, prText)
		}
	}

	// Splits
	if len(activity.SplitsMetric) > 1 {
		fmt.Printf("\n📏 SPLITS\n")
		for _, s := range activity.SplitsMetric {
			if s.Distance < 500 {
				continue // skip partial km
			}
			pace := 0.0
			if s.AverageSpeed > 0 {
				pace = 1000 / s.AverageSpeed / 60
			}
			paceMin := int(pace)
			paceSec := int((pace - float64(paceMin)) * 60)
			hrText := ""
			if s.AverageHeartrate > 0 {
				hrText = fmt.Sprintf(" | %.0f bpm", s.AverageHeartrate)
			}
			fmt.Printf("• Km %d: %d:%02d%s\n", s.Split, paceMin, paceSec, hrText)
		}
	}
}

func cmdZones() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: strava-cli zones <activity_id>")
		os.Exit(1)
	}

	activityID, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid activity ID: %s\n", os.Args[2])
		os.Exit(1)
	}

	client := getClient()

	zones, err := client.GetZones(activityID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching zones: %v\n", err)
		os.Exit(1)
	}

	for _, z := range zones {
		if z.Type == "heartrate" {
			fmt.Printf("❤️ HR ZONES\n\n")
			var totalTime float64
			for _, b := range z.Buckets {
				totalTime += b.Time
			}
			for i, b := range z.Buckets {
				pct := 0.0
				if totalTime > 0 {
					pct = b.Time / totalTime * 100
				}
				mins := int(b.Time / 60)
				secs := int(b.Time) % 60
				maxStr := fmt.Sprintf("%.0f", b.Max)
				if b.Max < 0 {
					maxStr = "max"
				}
				fmt.Printf("Zone %d (%.0f-%s bpm): %d:%02d (%.0f%%)\n", i+1, b.Min, maxStr, mins, secs, pct)
			}
		}
	}
}

func cmdLaps() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: strava-cli laps <activity_id>")
		os.Exit(1)
	}

	activityID, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid activity ID: %s\n", os.Args[2])
		os.Exit(1)
	}

	client := getClient()

	laps, err := client.GetLaps(activityID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching laps: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📏 SPLITS\n\n")
	for _, lap := range laps {
		pace := 0.0
		if lap.AverageSpeed > 0 {
			pace = 1000 / lap.AverageSpeed / 60
		}
		paceMin := int(pace)
		paceSec := int((pace - float64(paceMin)) * 60)
		fmt.Printf("%d. %.2f km | %d:%02d /km | %.0f bpm\n",
			lap.LapIndex, lap.Distance/1000, paceMin, paceSec, lap.AverageHeartrate)
	}
}

func printUsage() {
	fmt.Println(`strava-cli - Strava API command-line tool

Usage:
  strava-cli <command> [options]

Commands:
  profile           Show athlete profile
  activities        List recent activities
  activity <id>     Show activity details
  zones <id>        Show HR zones for activity
  laps <id>         Show splits/laps for activity
  stats             Show athlete statistics
  help              Show this help message

Options:
  activities:
    --limit <n>   Number of activities to show (default: 10)

Examples:
  strava-cli profile
  strava-cli activities --limit 20
  strava-cli activity 17462649839
  strava-cli zones 17462649839
  strava-cli laps 17462649839
  strava-cli stats`)
}
