package repositories

import (
	"context"
	"errors"
	"github.com/henrion-y/base.services/domain/repository"
	"gorm.io/gorm"
	"time"

	"user.services/app/domain/models/mysql_models"
	"user.services/pkg/utils"
)

type UserRepository interface {
	FindUserBaseByPhone(c context.Context, phone string) (mysql_models.TUserBase, error)
	FindUserBaseByUserId(c context.Context, userId uint64) (mysql_models.TUserBase, error)
	FindUserBaseByPassword(c context.Context, phone string, userId uint64, password string) (mysql_models.TUserBase, error)
	CreateUserBaseByPhone(c context.Context, phone string) (mysql_models.TUserBase, error)
	FindOauthUserBaseByCode(c context.Context, code string) (mysql_models.TUserBase, error)
	CreateOauthUserBaseByCode(c context.Context, code string) (mysql_models.TUserBase, error)
	UpdateUserBase(c context.Context, userBase *mysql_models.TUserBase) error
}

type userRepository struct {
	gormDb         *gorm.DB
	baseRepository repository.BaseRepository
}

func NewUserRepository(gormDb *gorm.DB, baseRepository repository.BaseRepository) UserRepository {
	gormDb.AutoMigrate(&mysql_models.TUserBase{}, &mysql_models.TOauthUser{})

	return &userRepository{
		gormDb:         gormDb,
		baseRepository: baseRepository,
	}
}

func (r *userRepository) FindUserBaseByPhone(c context.Context, phone string) (mod mysql_models.TUserBase, err error) {
	userBase := mysql_models.TUserBase{}
	err = r.baseRepository.FindOne(c, &mod, nil,
		repository.NewFilterGroup().Equals("phone", phone).IsNull("deleted_at"), nil)
	return userBase, err
}

func (r *userRepository) FindUserBaseByUserId(c context.Context, userId uint64) (mysql_models.TUserBase, error) {
	mod := mysql_models.TUserBase{}
	err := r.baseRepository.FindOne(c, &mod, nil,
		repository.NewFilterGroup().Equals("user_id", userId).IsNull("deleted_at"), nil)

	return mod, err
}

func (r *userRepository) FindUserBaseByPassword(c context.Context, phone string, userId uint64, password string) (mysql_models.TUserBase, error) {
	mod := mysql_models.TUserBase{}
	filterGroup := repository.NewFilterGroup().Equals("password", password).IsNull("deleted_at")
	if phone != "" {
		filterGroup.Equals("phone", phone)
	}
	if userId != 0 {
		filterGroup.Equals("user_id", userId)
	}

	err := r.baseRepository.FindOne(c, &mod, nil, filterGroup, nil)
	return mod, err
}

func (r *userRepository) CreateUserBaseByPhone(c context.Context, phone string) (mysql_models.TUserBase, error) {
	nowTime := time.Now()
	username, err := utils.GenerateUsername()
	if err != nil {
		return mysql_models.TUserBase{}, err
	}
	mod := mysql_models.TUserBase{
		Phone:        phone,
		Username:     username,
		HeadPortrait: mysql_models.DefaultHeadPortrait,
		UpdatedAt:    nowTime,
		CreatedAt:    nowTime,
	}
	err = r.baseRepository.Create(c, &mod)
	return mod, err
}

// UpdateUserBase 这里修改用户信息可能需要同步给其他服务
func (r *userRepository) UpdateUserBase(c context.Context, userBase *mysql_models.TUserBase) error {
	if userBase.UserId == 0 {
		return errors.New("用户不存在")
	}
	err := r.gormDb.Table(userBase.TableName()).Where("user_id=?", userBase.UserId).Updates(userBase).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindOauthUserBaseByCode(c context.Context, code string) (mod mysql_models.TUserBase, err error) {
	oauthUser := mysql_models.TOauthUser{}
	err = r.baseRepository.FindOne(c, &oauthUser, nil, repository.NewFilterGroup().Equals("code", code), nil)
	if err != nil {
		return mod, err
	}
	if oauthUser.UserId == 0 {
		return mod, nil
	}
	return r.FindUserBaseByUserId(c, oauthUser.UserId)
}

func (r *userRepository) CreateOauthUserBaseByCode(c context.Context, code string) (mod mysql_models.TUserBase, err error) {
	userBase, err := r.CreateUserBaseByPhone(c, "")
	if err != nil {
		return mod, err
	}

	nowTime := time.Now()
	oauthUserId, err := utils.GenerateId()
	if err != nil {
		return mod, err
	}
	username, err := utils.GenerateUsername()
	if err != nil {
		return mod, err
	}
	oauthUser := &mysql_models.TOauthUser{
		UserId:       userBase.UserId,
		OauthUserId:  oauthUserId,
		Code:         code,
		Username:     username,
		HeadPortrait: mysql_models.DefaultHeadPortrait,
		UpdatedAt:    nowTime,
		CreatedAt:    nowTime,
	}
	err = r.gormDb.Table(oauthUser.TableName()).Create(oauthUser).Error
	if err != nil {
		return mod, err
	}
	return userBase, nil
}
