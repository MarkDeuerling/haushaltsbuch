package refreshtoken

// import (
// 	"fmt"
// 	"time"

// )

// // Presenter Interface
// type Presenter interface {
// 	TokenSuccessful(userID user.ID, refreshToken, accessToken string)
// 	TokenFailed(msg string)
// }

// // TokenGenerator usecase Interface
// type TokenGenerator interface {
// 	RefreshAuthTokens(refreshToken string)
// }

// type refreshTokenUseCase struct {
// 	presenter            Presenter
// 	tokenGenerateService port.TokenGenerator
// 	tokenValidateService port.TokenParser
// 	tokenRepository      port.TokenRepository
// 	scopeRepository      port.ScopeRepository
// 	accessTokenExpire    time.Duration
// 	refreshTokenExpire   time.Duration
// }

// // NewRefreshTokenUseCase creates a new RefreshTokenUseCase
// func NewRefreshTokenUseCase(
// 	presenter Presenter,
// 	tokenGenerateService port.TokenGenerator,
// 	tokenValidateService port.TokenParser,
// 	tokenRepository port.TokenRepository,
// 	scopeRepository port.ScopeRepository,
// 	accessTokenExpire time.Duration,
// 	refreshTokenExpire time.Duration,
// ) TokenGenerator {
// 	return &refreshTokenUseCase{
// 		presenter:            presenter,
// 		tokenGenerateService: tokenGenerateService,
// 		tokenValidateService: tokenValidateService,
// 		tokenRepository:      tokenRepository,
// 		scopeRepository:      scopeRepository,
// 		accessTokenExpire:    accessTokenExpire,
// 		refreshTokenExpire:   refreshTokenExpire,
// 	}
// }

// // RefreshAuthTokens usecase signature
// func (r *refreshTokenUseCase) RefreshAuthTokens(refreshToken string) {
// 	claims, err := r.tokenValidateService.Parse(refreshToken)
// 	if err != nil {
// 		r.presenter.TokenFailed(fmt.Sprintf("invalid token, %s", err.Error()))
// 		return
// 	}

// 	userID := claims.Sub

// 	exist, err := r.tokenRepository.Exists(userID, refreshToken)
// 	if err != nil {
// 		r.presenter.TokenFailed(err.Error())
// 		return
// 	}

// 	if !exist {
// 		r.presenter.TokenFailed("invalid token")
// 		return
// 	}

// 	scope, err := r.scopeRepository.FindScopesByUserID(userID)
// 	if err != nil {
// 		r.presenter.TokenFailed(err.Error())
// 		return
// 	}

// 	accessToken, err := r.tokenGenerateService.GenerateAccessToken(userID, r.accessTokenExpire*time.Second, scope)
// 	if err != nil {
// 		r.presenter.TokenFailed(err.Error())
// 		return
// 	}

// 	refreshToken, err = r.tokenGenerateService.GenerateRefreshToken(userID, r.refreshTokenExpire*time.Second)
// 	if err != nil {
// 		r.presenter.TokenFailed(err.Error())
// 		return
// 	}

// 	err = r.tokenRepository.Delete(userID, refreshToken)
// 	if err != nil {
// 		r.presenter.TokenFailed(err.Error())
// 		return
// 	}

// 	err = r.tokenRepository.Save(userID, refreshToken)
// 	if err != nil {
// 		r.presenter.TokenFailed(err.Error())
// 		return
// 	}

// 	r.presenter.TokenSuccessful(userID, accessToken, refreshToken)
// }
