package services

import (
	"context"
	"github.com/barweiss/go-tuple"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/businessrules"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	cDTOs "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/keydtos"
	nDTOs "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/notedtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type NoteService interface {
	AddNoteTransaction(
		ctx context.Context,
		userBo userbos.UserBo,
		dto cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto],
	) single.Single[cDTOs.SuccessDto]
	UpdateNoteTransaction(
		ctx context.Context,
		userBo userbos.UserBo,
		dto cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto],
	) single.Single[cDTOs.SuccessDto]
	DeleteNoteTransaction(ctx context.Context, userBo userbos.UserBo, noteId string) single.Single[cDTOs.SuccessDto]
	GetNoteById(ctx context.Context, userBo userbos.UserBo, noteId string) single.Single[nDTOs.NoteDetailsDto]
	GetNotes(
		ctx context.Context,
		userBo userbos.UserBo,
		pageRequest pagination.PageRequest,
	) single.Single[pagination.Page[nDTOs.NoteReadDto]]
}

type NoteServiceImpl struct {
	noteRepository repositories.NoteRepository
	userKeyService externalservices.ExtUserKeyService
	crudDSHandler  dshandlers.CrudDSHandler
	errorService   sharedservices.ErrorService
	noteBr         businessrules.NoteBr
}

func (n NoteServiceImpl) AddNoteTransaction(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto],
) single.Single[cDTOs.SuccessDto] {
	return dshandlers.TransactionalSingle(ctx, n.crudDSHandler,
		func(_ dshandlers.Session, ctx context.Context) single.Single[cDTOs.SuccessDto] {
			sessDto, noteCreateDto := sessReqDto.Session, sessReqDto.Value
			keyDtoSrc := n.userKeyService.GetKeyFromSession(ctx, sessDto).ScheduleLazyAndCache(ctx)
			keySrc := single.MapWithError(keyDtoSrc, keydtos.UserKeyDto.GetKey).ScheduleLazyAndCache(ctx)
			titleCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteCreateDto.Title))
			})
			textCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteCreateDto.GetText()))
			})
			noteSaveSrc := single.FlatMap(single.Zip3(keyDtoSrc, titleCipherSrc, textCipherSrc),
				func(t tuple.T3[keydtos.UserKeyDto, []byte, []byte]) single.Single[models.Note] {
					keyDto, titleCipher, textCipher := t.V1, t.V2, t.V3
					note := models.Note{
						UserId:      userBo.Id,
						TextCipher:  textCipher,
						TitleCipher: titleCipher,
						KeyVersion:  keyDto.KeyVersion,
					}
					return n.noteRepository.Create(ctx, note)
				})
			return single.Map(noteSaveSrc, func(_ models.Note) cDTOs.SuccessDto {
				return cDTOs.NewSuccessTrue()
			})
		})

}

func (n NoteServiceImpl) UpdateNoteTransaction(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto],
) single.Single[cDTOs.SuccessDto] {
	return dshandlers.TransactionalSingle(ctx, n.crudDSHandler,
		func(_ dshandlers.Session, ctx context.Context) single.Single[cDTOs.SuccessDto] {
			sessDto, noteUpdateDto := sessReqDto.Session, sessReqDto.Value
			existingSrc := n.getExistingNote(ctx, noteUpdateDto.Id).ScheduleLazyAndCache(ctx)
			keyDtoSrc := n.userKeyService.GetKeyFromSession(ctx, sessDto).ScheduleLazyAndCache(ctx)
			keySrc := single.FlatMap(single.Zip2(existingSrc, keyDtoSrc),
				func(t tuple.T2[models.Note, keydtos.UserKeyDto]) single.Single[[]byte] {
					existing, keyDto := t.V1, t.V2
					noteUpdateValidationSrc := n.noteBr.ValidateNoteUpdate(userBo, sessDto, existing)
					return single.MapWithError(noteUpdateValidationSrc, func(_ []apperrors.RuleError) ([]byte, error) {
						return keyDto.GetKey()
					})
				}).ScheduleLazyAndCache(ctx)
			titleCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteUpdateDto.Title))
			})
			textCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteUpdateDto.GetText()))
			})
			noteSaveSrc := single.FlatMap(single.Zip4(existingSrc, keyDtoSrc, titleCipherSrc, textCipherSrc),
				func(t tuple.T4[models.Note, keydtos.UserKeyDto, []byte, []byte]) single.Single[models.Note] {
					existing, keyDto, titleCipher, textCipher := t.V1, t.V2, t.V3, t.V4
					existing.TitleCipher = titleCipher
					existing.TextCipher = textCipher
					existing.KeyVersion = keyDto.KeyVersion
					return n.noteRepository.Update(ctx, existing)
				})
			return single.Map(noteSaveSrc, func(_ models.Note) cDTOs.SuccessDto {
				return cDTOs.NewSuccessTrue()
			})
		})
}

func (n NoteServiceImpl) DeleteNoteTransaction(
	ctx context.Context,
	userBo userbos.UserBo,
	noteId string,
) single.Single[cDTOs.SuccessDto] {
	return dshandlers.TransactionalSingle(ctx, n.crudDSHandler,
		func(_ dshandlers.Session, ctx context.Context) single.Single[cDTOs.SuccessDto] {
			existingSrc := n.getExistingNote(ctx, noteId).ScheduleLazyAndCache(ctx)
			validateDeleteSrc := single.FlatMap(existingSrc,
				func(existing models.Note) single.Single[[]apperrors.RuleError] {
					return n.noteBr.ValidateNoteDelete(userBo, existing)
				})
			noteDeleteSrc := single.FlatMap(single.Zip2(existingSrc, validateDeleteSrc),
				func(t tuple.T2[models.Note, []apperrors.RuleError]) single.Single[models.Note] {
					existing := t.V1
					return n.noteRepository.Delete(ctx, existing)
				})
			return single.Map(noteDeleteSrc, func(_ models.Note) cDTOs.SuccessDto {
				return cDTOs.NewSuccessTrue()
			})
		})
}

func (n NoteServiceImpl) GetNoteById(
	ctx context.Context,
	userBo userbos.UserBo,
	noteId string,
) single.Single[nDTOs.NoteDetailsDto] {
	//TODO implement me
	panic("implement me")
}

func (n NoteServiceImpl) GetNotes(
	ctx context.Context,
	userBo userbos.UserBo,
	pageRequest pagination.PageRequest,
) single.Single[pagination.Page[nDTOs.NoteReadDto]] {
	//TODO implement me
	panic("implement me")
}

func (n NoteServiceImpl) getExistingNote(ctx context.Context, id string) single.Single[models.Note] {
	findSrc := n.noteRepository.FindById(ctx, id)
	return single.FlatMap(findSrc, func(m option.Maybe[models.Note]) single.Single[models.Note] {
		return option.Map(m, single.Just[models.Note]).
			OrElseGet(func() single.Single[models.Note] {
				ruleErr := n.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound)
				return single.Error[models.Note](apperrors.NewBadReqErrorFromRuleError(ruleErr))
			})
	})
}

func NewNoteServiceImpl(
	noteRepository repositories.NoteRepository,
	userKeyService externalservices.ExtUserKeyService,
	crudDSHandler dshandlers.CrudDSHandler,
	errorService sharedservices.ErrorService,
	noteBr businessrules.NoteBr,
) *NoteServiceImpl {
	return &NoteServiceImpl{
		noteRepository: noteRepository,
		userKeyService: userKeyService,
		crudDSHandler:  crudDSHandler,
		errorService:   errorService,
		noteBr:         noteBr,
	}
}
