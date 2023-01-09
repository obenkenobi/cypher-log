package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	bos "github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/businessobjects"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type AppSecretService interface {
	// GetAppSecret gets the app secret marked by the KID
	GetAppSecret(ctx context.Context, kid string) (bos.AppSecretBo, error)
	// GetPrimaryAppSecret gets the primary app secret
	GetPrimaryAppSecret(ctx context.Context) (bos.AppSecretBo, error)
	// GeneratePrimaryAppSecret generates a new primary app secret
	GeneratePrimaryAppSecret(ctx context.Context) (bos.AppSecretBo, error)
}

type AppSecretServiceImpl struct {
	primaryAppSecretRefRepository repositories.PrimaryAppSecretRefRepository
	appSecretRepository           repositories.AppSecretRepository
	keyConf                       conf.KeyConf
	errorService                  sharedservices.ErrorService
}

func (a AppSecretServiceImpl) GetPrimaryAppSecret(ctx context.Context) (bos.AppSecretBo, error) {
	refFind, err := a.primaryAppSecretRefRepository.Get(ctx)
	if err != nil {
		return bos.AppSecretBo{}, err
	}
	if ref, ok := refFind.Get(); ok {
		return a.GetAppSecret(ctx, ref.Kid)
	}
	return a.GeneratePrimaryAppSecret(ctx)
}

func (a AppSecretServiceImpl) GetAppSecret(ctx context.Context, kid string) (bos.AppSecretBo, error) {
	appSecretFind, err := a.appSecretRepository.Get(ctx, kid)
	if err != nil {
		return bos.AppSecretBo{}, err
	}
	appSecretBoMaybe := option.Map(appSecretFind, func(appSecret models.AppSecret) bos.AppSecretBo {
		return bos.NewAppSecretBo(kid, appSecret.SecretKey)
	})
	return appSecretServiceReadMaybeModel(a, appSecretBoMaybe)

}

func (a AppSecretServiceImpl) GeneratePrimaryAppSecret(ctx context.Context) (bos.AppSecretBo, error) {
	kidGuid, err := uuid.NewRandom()
	if err != nil {
		return bos.AppSecretBo{}, err
	}
	kid := kidGuid.String()
	if utils.StringIsBlank(kid) {
		return bos.AppSecretBo{}, errors.New("generated KID is blank")
	}
	key, err := cipherutils.GenerateRandomKeyAES()
	if err != nil {
		return bos.AppSecretBo{}, err
	}

	ref := models.PrimaryAppSecretRef{Kid: kid}
	appSecret := models.AppSecret{SecretKey: key}

	if _, err := a.appSecretRepository.Set(ctx, ref.Kid, appSecret, a.keyConf.GetSecretDuration()); err != nil {
		return bos.AppSecretBo{}, err
	}
	if _, err = a.primaryAppSecretRefRepository.Set(ctx, ref, a.keyConf.GetPrimaryAppSecretDuration()); err != nil {
		return bos.AppSecretBo{}, err
	}
	return bos.NewAppSecretBo(ref.Kid, appSecret.SecretKey), nil
}

func appSecretServiceReadMaybeModel[T any](a AppSecretServiceImpl, maybe option.Maybe[T]) (T, error) {
	val, ok := maybe.Get()
	var err error = nil
	if !ok {
		ruleErr := a.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound)
		err = apperrors.NewBadReqErrorFromRuleError(ruleErr)
	}
	return val, err
}

func NewAppSecretServiceImpl(
	primaryAppSecretRefRepository repositories.PrimaryAppSecretRefRepository,
	appSecretRepository repositories.AppSecretRepository,
	keyConf conf.KeyConf,
	errorService sharedservices.ErrorService,
) *AppSecretServiceImpl {
	return &AppSecretServiceImpl{
		primaryAppSecretRefRepository: primaryAppSecretRefRepository,
		appSecretRepository:           appSecretRepository,
		keyConf:                       keyConf,
		errorService:                  errorService,
	}
}
