package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/gofiber/fiber/v2"
)

// * Sing in with OAuth handler
// @Summary      Sign in OAuth
// @Description  Sing in by OAuth provider (google/github/discord)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        code query string true "OAuth code"
// @Param        provider path string true "OAuth provider" Enums(google, discord, github)
// @Success      200  {object}  dtos.SignInRespOKDTO
// @Failure      406  {object}  dtos.DetailRespErrDTO
// @Failure      400  {object}  dtos.RespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /auth/{provider}/callback [get]
func (h *AuthHandler) SignInWihOAuth(c *fiber.Ctx) error {
	code := c.Query("code")
	provider := c.Params("provider")
	providers := utils.GetOAuthProviders()
	cfg := config.Get()

	if code == "" {
		return h.RespErr(c, 400, "missing oauth code from query")
	}

	if !utils.ItemInSlice(provider, providers) {
		return h.RespErr(c, 406,
			"invalid oauth provider",
			fmt.Sprintf("The available providers are: %s", strings.Join(providers, ", ")),
		)
	}

	var oauthSvc auth.OAuthService

	switch provider {
	case utils.GoogleProvider:
		oauthSvc = auth.NewGoogleOAuthService()
	case utils.DiscordProvider:
		oauthSvc = auth.NewDiscordOAuthService()
	case utils.GithubProvider:
		oauthSvc = auth.NewGithubOAuthService()
	}

	oauthUser, err := oauthSvc.GetOAuthUser(code)

	if err != nil {
		return h.RespErr(c, 500, "error getting user from "+provider, err.Error())
	}

	user, err := h.usrSvc.GetByEmail(oauthUser.Email)

	if err != nil {
		return h.RespErr(c, 500, "searching user error", err.Error())
	}

	if user == nil {
		if err := h.usrSvc.Create(oauthUser); err != nil {
			return h.RespErr(c, 500, "creating user error", err.Error())
		}
		user = oauthUser
	}

	// if !user.VerifiedEmail {
	// 	go func() {
	// 		tokenCode, err := h.jwtSvc.CreateToken(
	// 			auth.Claims{ID: user.ID, Role: user.Role},
	// 			time.Hour*24*3, cfg.Api.VerificationSecret,
	// 		)

	// 		if err != nil {
	// 			fmt.Println("Error sending email")
	// 			return
	// 		}

	// 		err = h.emailSvc.SendEmail(email.EmailData{
	// 			Email:    user.Email,
	// 			Subject:  "Pulse | Verify your email!",
	// 			Template: "change_password.html",
	// 			Data: fiber.Map{
	// 				"Name":  user.Username,
	// 				"Email": user.Email,
	// 				"Code":  tokenCode,
	// 			},
	// 		})

	// 		if err != nil {
	// 			fmt.Println("Error sending email")
	// 			return
	// 		}

	// 		fmt.Println(">>> Email Sent to:", user.Email)
	// 	}()
	// }

	refreshToken, rtErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Hour*24*5, cfg.Api.RefreshTokenSecret,
	)

	if rtErr != nil {
		return h.RespErr(c, 500, "error creating tokens", "something went wrong")
	}

	c.Cookie(&fiber.Cookie{
		Name:     cfg.Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return c.Redirect(cfg.Api.ClientOrigin)
	
	// return h.RespOK(c, 200, "sign in successfully", fiber.Map{
	// 	"user":          user,
	// 	"access_token":  accessToken,
	// 	"refresh_token": refreshToken,
	// })
}
