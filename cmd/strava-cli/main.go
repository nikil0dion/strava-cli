package main

import (
	"flag"
	"fmt"
	"os"

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

func printUsage() {
	fmt.Println(`strava-cli - Strava API command-line tool

Usage:
  strava-cli <command> [options]

Commands:
  profile       Show athlete profile
  activities    List recent activities
  stats         Show athlete statistics
  help          Show this help message

Options:
  activities:
    --limit <n>   Number of activities to show (default: 10)

Examples:
  strava-cli profile
  strava-cli activities --limit 20
  strava-cli stats`)
}
