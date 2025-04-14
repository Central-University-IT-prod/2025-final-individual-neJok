package campaign

import (
	"neJok/solution/internal/model"
)

func CalculateCampaignScore(campaign model.CampaignForUser, maxScore int, maxEndDate int32, minEndDate int32, maxCostPerImpression float64, maxCostPerClick float64) float64 {
	var profit float64
	if maxCostPerImpression != 0 {
		profit = campaign.CostPerImpression / maxCostPerImpression
	} else {
		profit = 0
	}
	var clickProbability float64

	if maxScore > 0 {
		clickProbability = campaign.Score / float64(maxScore)
	} else {
		clickProbability = 0.0
	}

	if maxCostPerClick != 0 {
		profit += campaign.CostPerClick / maxCostPerClick * clickProbability
	}

	var endDateProbability float64
	if maxEndDate-minEndDate != 0 {
		endDateProbability = float64(maxEndDate-campaign.EndDate) / float64(maxEndDate-minEndDate)
	} else {
		endDateProbability = 0.0
	}

	totalScore := 0.5*profit + clickProbability*0.25 + 0.1*endDateProbability
	return totalScore
}
