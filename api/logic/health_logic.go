package logic

type HealthLogic interface{}

type L0HealthLogic struct {
	Logic
}

func NewL0HealthLogic(l Logic) *L0HealthLogic {
	return &L0HealthLogic{
		Logic: l,
	}
}
