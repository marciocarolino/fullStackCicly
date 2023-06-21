package transformer


import (

	"github.com/marciocarolino/fullStackCycle/internal/market/dto"
	"github.com/marciocarolino/fullStackCycle/internal/market/entity"
)

func TransformerInput(input dto.TradeInput) *entity.Order{
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entity.NewInvestor(input.InvestorID)
	order := entity.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType)

	if input.CurrentShares > 0 {
		assetPosition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}
	return order
}


func TransformOutPut(order *entity.Order) *dto.OrderOutPut{
	output := &dto.OrderOutPut{
		OrderID: order.ID,
		InvestorID: order.Investor.ID,
		AssetID: order.Asset.ID,
		OrderType: order.OrderType,
		Status: order.Staus,
		Partil: order.PendingShares,
		Shares: order.Shares,
	}

	var transactionsOutPut []*dto.TransactionOutPut
	for _, t := range order.Transaction {
		TransactionOutPut := &dto.TransactionOutPut{
			TransactionID:	t.ID,
			BuyerID:	t.BuyingOrder.ID,
			SellerID: 	t.SellingOrder.ID,
			AssetID:	t.SellingOrder.Asset.ID,
			Price:		t.Price,
			Shares:	t.SellingOrder.Shares - t.SellingOrder.PendingShares,
		}
		TransactionOutPut = append(transactionsOutPut, TransactionOutPut)
	}
	output.TransactionOutput = transactionOutput
	return output
}