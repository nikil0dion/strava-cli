package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://www.strava.com/api/v3"

type Client struct {
	accessToken string
	httpClient  *http.Client
}

func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		httpClient:  &http.Client{},
	}
}

func (c *Client) get(endpoint string) ([]byte, error) {
	url := baseURL + endpoint
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, body)
	}

	return body, nil
}

type Athlete struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	FirstName string  `json:"firstname"`
	LastName  string  `json:"lastname"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Sex       string  `json:"sex"`
	Weight    float64 `json:"weight"`
}

func (c *Client) GetAthlete() (*Athlete, error) {
	data, err := c.get("/athlete")
	if err != nil {
		return nil, err
	}

	var athlete Athlete
	if err := json.Unmarshal(data, &athlete); err != nil {
		return nil, fmt.Errorf("failed to parse athlete: %w", err)
	}

	return &athlete, nil
}

type Gear struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Distance float64 `json:"distance"`
	Primary  bool    `json:"primary"`
}

type BestEffort struct {
	Name        string `json:"name"`
	Distance    float64 `json:"distance"`
	ElapsedTime int    `json:"elapsed_time"`
	MovingTime  int    `json:"moving_time"`
	PRRank      *int   `json:"pr_rank"`
}

type Split struct {
	Split            int     `json:"split"`
	Distance         float64 `json:"distance"`
	MovingTime       int     `json:"moving_time"`
	AverageSpeed     float64 `json:"average_speed"`
	AverageHeartrate float64 `json:"average_heartrate"`
}

type Activity struct {
	ID                 int64        `json:"id"`
	Name               string       `json:"name"`
	Type               string       `json:"type"`
	SportType          string       `json:"sport_type"`
	Distance           float64      `json:"distance"`
	MovingTime         int          `json:"moving_time"`
	ElapsedTime        int          `json:"elapsed_time"`
	TotalElevationGain float64      `json:"total_elevation_gain"`
	StartDate          string       `json:"start_date"`
	StartDateLocal     string       `json:"start_date_local"`
	AverageSpeed       float64      `json:"average_speed"`
	MaxSpeed           float64      `json:"max_speed"`
	AverageHeartrate   float64      `json:"average_heartrate"`
	MaxHeartrate       float64      `json:"max_heartrate"`
	AverageCadence     float64      `json:"average_cadence"`
	Calories           float64      `json:"calories"`
	SufferScore        float64      `json:"suffer_score"`
	Description        string       `json:"description"`
	PRCount            int          `json:"pr_count"`
	Gear               *Gear        `json:"gear"`
	BestEfforts        []BestEffort `json:"best_efforts"`
	SplitsMetric       []Split      `json:"splits_metric"`
}

func (c *Client) GetActivities(perPage, page int) ([]Activity, error) {
	endpoint := fmt.Sprintf("/athlete/activities?per_page=%d&page=%d", perPage, page)
	
	data, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var activities []Activity
	if err := json.Unmarshal(data, &activities); err != nil {
		return nil, fmt.Errorf("failed to parse activities: %w", err)
	}

	return activities, nil
}

type Stats struct {
	RecentRunTotals struct {
		Count    int     `json:"count"`
		Distance float64 `json:"distance"`
		MovingTime float64   `json:"moving_time"`
		ElapsedTime float64  `json:"elapsed_time"`
		ElevationGain float64 `json:"elevation_gain"`
	} `json:"recent_run_totals"`
	RecentRideTotals struct {
		Count    int     `json:"count"`
		Distance float64 `json:"distance"`
		MovingTime float64   `json:"moving_time"`
		ElapsedTime float64  `json:"elapsed_time"`
		ElevationGain float64 `json:"elevation_gain"`
	} `json:"recent_ride_totals"`
	YTDRunTotals struct {
		Count    int     `json:"count"`
		Distance float64 `json:"distance"`
		MovingTime float64   `json:"moving_time"`
		ElapsedTime float64  `json:"elapsed_time"`
		ElevationGain float64 `json:"elevation_gain"`
	} `json:"ytd_run_totals"`
	YTDRideTotals struct {
		Count    int     `json:"count"`
		Distance float64 `json:"distance"`
		MovingTime float64   `json:"moving_time"`
		ElapsedTime float64  `json:"elapsed_time"`
		ElevationGain float64 `json:"elevation_gain"`
	} `json:"ytd_ride_totals"`
}

func (c *Client) GetStats(athleteID int64) (*Stats, error) {
	endpoint := fmt.Sprintf("/athletes/%d/stats", athleteID)

	data, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var stats Stats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}

	return &stats, nil
}

// GetActivity returns detailed info for a single activity
func (c *Client) GetActivity(activityID int64) (*Activity, error) {
	endpoint := fmt.Sprintf("/activities/%d", activityID)

	data, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var activity Activity
	if err := json.Unmarshal(data, &activity); err != nil {
		return nil, fmt.Errorf("failed to parse activity: %w", err)
	}

	return &activity, nil
}

// HR Zones types
type ZoneBucket struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Time float64 `json:"time"`
}

type ZoneDistribution struct {
	Type    string       `json:"type"`
	Buckets []ZoneBucket `json:"distribution_buckets"`
}

// GetZones returns HR zones for an activity
func (c *Client) GetZones(activityID int64) ([]ZoneDistribution, error) {
	endpoint := fmt.Sprintf("/activities/%d/zones", activityID)

	data, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var zones []ZoneDistribution
	if err := json.Unmarshal(data, &zones); err != nil {
		return nil, fmt.Errorf("failed to parse zones: %w", err)
	}

	return zones, nil
}

// Laps/Splits
type Lap struct {
	Name             string  `json:"name"`
	Distance         float64 `json:"distance"`
	MovingTime       int     `json:"moving_time"`
	ElapsedTime      int     `json:"elapsed_time"`
	AverageSpeed     float64 `json:"average_speed"`
	AverageHeartrate float64 `json:"average_heartrate"`
	MaxHeartrate     float64 `json:"max_heartrate"`
	LapIndex         int     `json:"lap_index"`
	PaceZone         int     `json:"pace_zone"`
}

// GetLaps returns splits/laps for an activity
func (c *Client) GetLaps(activityID int64) ([]Lap, error) {
	endpoint := fmt.Sprintf("/activities/%d/laps", activityID)

	data, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var laps []Lap
	if err := json.Unmarshal(data, &laps); err != nil {
		return nil, fmt.Errorf("failed to parse laps: %w", err)
	}

	return laps, nil
}
