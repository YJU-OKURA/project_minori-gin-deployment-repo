package repositories

import (
	"errors"
	"log"
	"strconv"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type ClassCodeRepository interface {
	FindByCode(code string) (*models.ClassCode, error)
	FindByClassID(cid uint) (*models.ClassCode, error)
	SaveClassCode(classCode *models.ClassCode) error
}

// ClassCodeRepository はグループコードのリポジトリです。
type classCodeRepository struct {
	db *gorm.DB
}

// NewClassCodeRepository はClassCodeRepositoryを生成します。
func NewClassCodeRepository(db *gorm.DB) ClassCodeRepository {
	return &classCodeRepository{db: db}
}

// FindByCode は指定されたコードのグループコードを取得します。
func (r *classCodeRepository) FindByCode(code string) (*models.ClassCode, error) {
	var classCode models.ClassCode
	result := r.db.Where("code = ?", code).First(&classCode)
	if result.Error != nil {
		// レコードが見つからない場合、nilを返します。
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// それ以外のエラーの場合はエラーを返します。
		return nil, result.Error
	}
	return &classCode, nil
}

// FindByClassID は指定されたクラスIDのクラスコードを取得します。
func (r *classCodeRepository) FindByClassID(cid uint) (*models.ClassCode, error) {
	var classCode models.ClassCode
	result := r.db.Where("cid = ?", cid).First(&classCode)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("ClassCode not found for ClassID: %d", cid)
		return nil, nil
	}

	if result.Error != nil {
		log.Printf("Error retrieving ClassCode for ClassID %d: %v", cid, result.Error)
		return nil, result.Error
	}

	log.Printf("ClassCode retrieved for ClassID %d: %+v", cid, classCode)
	return &classCode, nil
}

func (r *classCodeRepository) SaveClassCode(classCode *models.ClassCode) error {
	var class models.Class
	if err := r.db.First(&class, "id = ?", classCode.CID).Error; err != nil {
		return errors.New("invalid class ID: " + strconv.Itoa(int(classCode.CID)))
	}

	var user models.User
	if err := r.db.First(&user, "id = ?", classCode.UID).Error; err != nil {
		return errors.New("invalid user ID: " + strconv.Itoa(int(classCode.UID)))
	}

	return r.db.Create(classCode).Error
}
