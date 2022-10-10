package services

import (
	"context"
	"errors"
	"github.com/barweiss/go-tuple"
	"github.com/google/uuid"
	bos "github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/businessobjects"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type AppSecretService interface {
	// GetAppSecret gets the app secret marked by the KID
	GetAppSecret(ctx context.Context, kid string) single.Single[bos.AppSecretBo]
	// GetPrimaryAppSecret gets the primary app secret
	GetPrimaryAppSecret(ctx context.Context) single.Single[bos.AppSecretBo]
	// GeneratePrimaryAppSecret generates a new primary app secret
	GeneratePrimaryAppSecret(ctx context.Context) single.Single[bos.AppSecretBo]
}

type AppSecretServiceImpl struct {
	primaryAppSecretRefRepository repositories.PrimaryAppSecretRefRepository
	appSecretRepository           repositories.AppSecretRepository
	keyConf                       conf.KeyConf
	errorService                  sharedservices.ErrorService
}

func (a AppSecretServiceImpl) GetPrimaryAppSecret(ctx context.Context) single.Single[bos.AppSecretBo] {
	refFindSrc := a.primaryAppSecretRefRepository.Get(ctx)
	return single.FlatMap(refFindSrc,
		func(maybe option.Maybe[models.PrimaryAppSecretRef]) single.Single[bos.AppSecretBo] {
			if ref, ok := maybe.Get(); ok {
				return a.GetAppSecret(ctx, ref.Kid)
			}
			return a.GeneratePrimaryAppSecret(ctx)
		},
	)
}

func (a AppSecretServiceImpl) GetAppSecret(ctx context.Context, kid string) single.Single[bos.AppSecretBo] {
	appSecretFindSrc := a.appSecretRepository.Get(ctx, kid)
	return single.MapWithError(appSecretFindSrc, func(maybe option.Maybe[models.AppSecret]) (bos.AppSecretBo, error) {
		appSecretBoMaybe := option.Map(maybe, func(appSecret models.AppSecret) bos.AppSecretBo {
			return bos.NewAppSecretBo(kid, appSecret.SecretKey)
		})
		return appSecretServiceReadMaybeModel(a, appSecretBoMaybe)
	})
}

func (a AppSecretServiceImpl) GeneratePrimaryAppSecret(ctx context.Context) single.Single[bos.AppSecretBo] {
	kidGuidSrc := single.FromSupplier(uuid.NewRandom)
	kidSrc := single.MapWithError(kidGuidSrc, func(kidGuid uuid.UUID) (string, error) {
		newKid := kidGuid.String()
		if utils.StringIsBlank(newKid) {
			return newKid, errors.New("generated KID is blank")
		}
		return newKid, nil
	})
	newKeySrc := single.FromSupplier(cipherutils.GenerateRandomKeyAES)
	kidKeySrc := single.Zip2(kidSrc, newKeySrc)
	return single.FlatMap(kidKeySrc, func(t tuple.T2[string, []byte]) single.Single[bos.AppSecretBo] {
		ref := models.PrimaryAppSecretRef{Kid: t.V1}
		appSecret := models.AppSecret{SecretKey: t.V2}
		secretSaveSrc := a.appSecretRepository.Set(ctx, ref.Kid, appSecret, a.keyConf.GetSecretDuration())
		refSavedSrc := single.FlatMap(
			secretSaveSrc,
			func(_ models.AppSecret) single.Single[models.PrimaryAppSecretRef] {
				return a.primaryAppSecretRefRepository.Set(ctx, ref, a.keyConf.GetPrimaryAppSecretDuration())
			},
		)
		return single.Map(refSavedSrc, func(_ models.PrimaryAppSecretRef) bos.AppSecretBo {
			return bos.NewAppSecretBo(ref.Kid, appSecret.SecretKey)
		})
	})
}

func appSecretServiceReadMaybeModel[T any](a AppSecretServiceImpl, maybe option.Maybe[T]) (T, error) {
	val, ok := maybe.Get()
	var err error = nil
	if !ok {
		err = apperrors.NewBadReqErrorFromRuleError(
			a.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound),
		)
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
