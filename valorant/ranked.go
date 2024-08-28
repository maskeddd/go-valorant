package valorant

import (
	"context"
	"fmt"
	"net/http"
)

type RankedService service

type RankTier string

const (
	RankTierImmortal1 RankTier = "21"
	RankTierImmortal2 RankTier = "22"
	RankTierImmortal3 RankTier = "23"
	RankTierRadiant   RankTier = "24"
)

type Leaderboard struct {
	Shard                 string                   `json:"shard"`
	ActID                 string                   `json:"actId"`
	TotalPlayers          int                      `json:"totalPlayers"`
	Players               []Player                 `json:"players"`
	ImmortalStartingIndex int                      `json:"immortalStartingIndex"`
	ImmortalStartingPage  int                      `json:"immortalStartingPage"`
	StartIndex            int                      `json:"startIndex"`
	TierDetails           map[RankTier]TierDetails `json:"tierDetails"`
	TopTierRRThreshold    int                      `json:"topTierRRThreshold"`
}

type Player struct {
	PUUID           *string `json:"puuid,omitempty"`
	GameName        *string `json:"gameName,omitempty"`
	TagLine         *string `json:"tagLine,omitempty"`
	LeaderboardRank int     `json:"leaderboardRank"`
	RankedRating    int     `json:"rankedRating"`
	NumberOfWins    int     `json:"numberOfWins"`
	CompetitiveTier int     `json:"competitiveTier"`
}

type TierDetails struct {
	RankedRatingThreshold int `json:"rankedRatingThreshold"`
	StartingIndex         int `json:"startingIndex"`
	StartingPage          int `json:"startingPage"`
}

type LeaderboardListOptions struct {
	Size       int `url:"size,omitempty"`
	StartIndex int `url:"startIndex,omitempty"`
}

// ListLeaderboardByAct gets the leaderboard for a specific act.
//
// Valorant API docs: https://developer.riotgames.com/apis#val-ranked-v1/GET_getLeaderboard
func (s *RankedService) ListLeaderboardByAct(ctx context.Context, actID string, opts *LeaderboardListOptions) (*Leaderboard, *http.Response, error) {
	u := fmt.Sprintf("ranked/v1/leaderboards/by-act/%s", actID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var leaderboard *Leaderboard
	resp, err := s.client.Do(ctx, req, &leaderboard)
	if err != nil {
		return nil, resp, err
	}

	return leaderboard, resp, nil
}
