package services

import (
	"goDDD1/config"
	"goDDD1/models"
)

type UserCurrencyFlowServiceInterface interface {
	CreateUserCurrencyFlow(*models.UserCurrencyFlow) error
	GetUserCurrencyFlow(string) (map[string]interface{}, error)
	GetAllUserCurrencyFlow() ([]models.UserCurrencyFlow, error)
}

type UserCurrencyFlowService struct {
}

func NewUserCurrencyFlowService() UserCurrencyFlowServiceInterface {
	return &UserCurrencyFlowService{}
}

func (s *UserCurrencyFlowService) CreateUserCurrencyFlow(userCurrencyFlow *models.UserCurrencyFlow) error {
	return config.Database.Create(userCurrencyFlow).Error
}

func (s *UserCurrencyFlowService) GetUserCurrencyFlow(userID string) (map[string]interface{}, error) {
	var userCurrencyFlows []models.UserCurrencyFlow
	err := config.Database.Where("user_id = ?", userID).Find(&userCurrencyFlows).Error
	if err != nil {
		return nil, err
	}

	// 构建map集合返回结果
	result := map[string]interface{}{
		"user_id":     userID,
		"total_count": len(userCurrencyFlows),
		"flows":       make([]map[string]interface{}, 0),
	}

	// 将每条记录转换为map
	flows := make([]map[string]interface{}, 0)
	for _, flow := range userCurrencyFlows {
		flowMap := map[string]interface{}{
			"user_id":   flow.UserID,
			"store_id":  flow.StoreID,
			"cost_type": flow.CostType,
			"price":     flow.Price,
			"ctime":     flow.Ctime,
		}
		flows = append(flows, flowMap)
	}

	result["flows"] = flows
	return result, nil
}

func (s *UserCurrencyFlowService) GetAllUserCurrencyFlow() ([]models.UserCurrencyFlow, error) {
	var userCurrencyFlows []models.UserCurrencyFlow
	err := config.Database.Find(&userCurrencyFlows).Error
	if err != nil {
		return nil, err
	}
	return userCurrencyFlows, nil
}
