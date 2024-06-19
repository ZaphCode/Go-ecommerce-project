package auth

import (
	"fmt"
	"strings"

	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/gofiber/fiber/v2"
)

// * Get OAuth url handler
// @Summary      Get OAuth url
// @Description  Get OAuth url from provider (google/github/discord)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param 		 provider path string true "OAuth provider" Enums(google, discord, github)
// @Success      200  {object}  dtos.URLRespOKDTO
// @Failure      400  {object}  dtos.DetailRespErrDTO
// @Router       /auth/{provider}/url [get]
func (h *AuthHandler) GetOAuthUrl(c *fiber.Ctx) error {
	provider := c.Params("provider")
	providers := utils.GetOAuthProviders()

	if !utils.ItemInSlice(provider, providers) {
		return h.RespErr(c, 400, "invalid oauth provider", fmt.Sprintf("The avalible proveders are: %s", strings.Join(providers, ", ")))
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

	return h.RespOK(c, 200, "OAuth url for "+provider, oauthSvc.GetOAuthUrl())
}
