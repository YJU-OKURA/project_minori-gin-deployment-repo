package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

// GroupCodeService はグループコードのサービスです。
type GroupCodeService interface {
	CheckSecretExists(code string) (bool, error)
	VerifyGroupCode(code, secret string) (bool, error)
}

// groupCodeServiceImpl はGroupCodeServiceの実装です。
type groupCodeServiceImpl struct {
	Repo *repositories.GroupCodeRepository
}

// NewGroupCodeService はGroupCodeServiceを生成します。
func NewGroupCodeService(repo *repositories.GroupCodeRepository) GroupCodeService {
	return &groupCodeServiceImpl{Repo: repo}
}

// CheckSecretExists は指定されたグループコードにシークレットがあるかどうかをチェックします。
func (service *groupCodeServiceImpl) CheckSecretExists(code string) (bool, error) {
	groupCode, err := service.Repo.FindByCode(code)
	if err != nil {
		return false, err
	}
	return groupCode.Secret != nil, nil
}

// VerifyGroupCode はグループコードと、該当する場合はそのシークレットを確認します。
func (service *groupCodeServiceImpl) VerifyGroupCode(code, secret string) (bool, error) {
	groupCode, err := service.Repo.FindByCode(code)
	if err != nil {
		return false, err
	}

	// シークレットがない場合は常にtrueを返す
	if groupCode.Secret == nil || *groupCode.Secret != secret {
		return false, nil
	}

	return true, nil
}
