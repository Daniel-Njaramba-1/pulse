package pricing

import (
	"fmt"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/sajari/regression"
)

type PricingModel struct {
	regression *regression.Regression
	version string
	trained bool
	sampleSize int
}

func NewPricingModel(version string) *PricingModel {
	r := new(regression.Regression)
	r.SetObserved("price")
	r.SetVar(0, "salesCount")
	r.SetVar(1, "salesValue")
	r.SetVar(2, "salesVelocity")
	r.SetVar(3, "daysSinceLastSale")
	r.SetVar(4, "categoryRank")
	r.SetVar(5, "categoryPercentile")
	r.SetVar(6, "reviewScore")
	r.SetVar(7, "wishlistRatio")
	r.SetVar(8, "daysInStock")
	r.SetVar(9, "seasonalFactor")

	return &PricingModel{
		regression: r,
		version: version,
		trained: false,
		sampleSize: 0,
	}
}

func (pm *PricingModel) Train(data []repo.PricingFeatures, prices []float64) error {
	if len(data) != len(prices) {
		return fmt.Errorf("data and prices must have same length")
	}
	
	pm.sampleSize = len(data)

	for i, feature := range data {
		pm.regression.Train(
			regression.DataPoint(prices[i],
				[]float64{
					float64(feature.TotalSalesCount),
					feature.TotalSalesValue,
					feature.SalesVelocity,
					float64(feature.DaysSinceLastSale),
					float64(feature.CategoryRank),
					feature.CategoryPercentile,
					feature.ReviewScore,
					feature.WishlistToSalesRatio,
					float64(feature.DaysInStock),
					feature.SeasonalFactor,
				},
			),
		)
	}

	err := pm.regression.Run()
	if err != nil {
		return err
	}

	pm.trained = true
	return nil
}

func (pm *PricingModel) GetCoefficients() repo.PriceModelCoefficients {
	return repo.PriceModelCoefficients{
		ModelVersion:           pm.version,
        TrainingDate:           time.Now(),
		SampleSize: 			pm.sampleSize,
        RSquared:               pm.regression.R2,
        Intercept:              pm.regression.Coeff(0),
        SalesCountCoef:         pm.regression.Coeff(1),
        SalesValueCoef:         pm.regression.Coeff(2),
        SalesVelocityCoef:      pm.regression.Coeff(3),
        DaysSinceSaleCoef:      pm.regression.Coeff(4),
        CategoryRankCoef:       pm.regression.Coeff(5),
        CategoryPercentileCoef: pm.regression.Coeff(6),
        ReviewScoreCoef:        pm.regression.Coeff(7),
        WishlistRatioCoef:      pm.regression.Coeff(8),
        DaysInStockCoef:        pm.regression.Coeff(9),
        SeasonalFactorCoef:     pm.regression.Coeff(10),
	}
}