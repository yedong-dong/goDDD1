package services

import (
	"goDDD1/config"
	"goDDD1/models"

	"github.com/jinzhu/gorm"
)

// RewardFlowService 奖励流水服务接口
type RewardFlowService interface {
	CreateRewardFlow(tx *gorm.DB, userID uint, itemType models.RewardFlowType, itemID uint, quantity int64, source string) error
	GetUserRewardFlows(userID uint, page, pageSize int) ([]*models.RewardFlow, int64, error)
}

// rewardFlowService 奖励流水服务实现
type rewardFlowService struct{}

// NewRewardFlowService 创建奖励流水服务实例
func NewRewardFlowService() RewardFlowService {
	return &rewardFlowService{}
}

// CreateRewardFlow 创建奖励流水记录
func (s *rewardFlowService) CreateRewardFlow(tx *gorm.DB, userID uint, itemType models.RewardFlowType, itemID uint, quantity int64, source string) error {
	flow := &models.RewardFlow{
		UserID:   userID,
		ItemType: itemType,
		ItemID:   itemID,
		Quantity: quantity,
		Source:   source,
	}

	// 根据是否传入事务决定使用哪个数据库连接
	if tx != nil {
		return tx.Create(flow).Error
	}
	return config.Database.Create(flow).Error
}

// GetUserRewardFlows 获取用户的奖励流水记录
func (s *rewardFlowService) GetUserRewardFlows(userID uint, page, pageSize int) ([]*models.RewardFlow, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	var flows []*models.RewardFlow
	var total int64

	query := config.Database.Where("user_id = ?", userID)

	// 获取总数
	if err := query.Model(&models.RewardFlow{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(pageSize).Order("ctime desc").Find(&flows).Error; err != nil {
		return nil, 0, err
	}

	return flows, total, nil
}
