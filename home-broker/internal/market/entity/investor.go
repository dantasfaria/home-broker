package entity

type Investor struct {
	ID            string
	Name          string
	AssetPosition []*InvestorAssetPosition
}

type InvestorAssetPosition struct {
	AssetID string
	Shares  int
}

func NewInvestor(id string) *Investor {
	return &Investor{
		ID:            id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	i.AssetPosition = append(i.AssetPosition, assetPosition)
}

func (i *Investor) GetAssetPosition(assetId string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assetId {
			return assetPosition
		}
	}
	return nil
}

func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares:  shares,
	}
}

func (i *Investor) UpdateAssetPosition(assetID string, qtdShares int) {
	assetPosition := i.GetAssetPosition(assetID)
	if assetPosition == nil {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetID, qtdShares))
	} else {
		assetPosition.Shares += qtdShares
	}
}
