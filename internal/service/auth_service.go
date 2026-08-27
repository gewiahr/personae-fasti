package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/mail"
	"personae-fasti/configs"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/validation"
	"personae-fasti/internal/repo"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authConfig *configs.AuthConfig
	playerRepo repo.PlayerRepository
}

func NewAuthService(authConfig *configs.AuthConfig, playerRepo repo.PlayerRepository) *AuthService {
	return &AuthService{
		authConfig: authConfig,
		playerRepo: playerRepo,
	}
}

func (s *AuthService) AuthenticateToken(ctx context.Context, tokenString string) (*domain.Player, error) {
	if tokenString == "" {
		return nil, e.NewUnauthorizedError("token is required")
	}

	tokenHash := sha256.Sum256([]byte(tokenString))
	tokenHashHex := hex.EncodeToString(tokenHash[:])

	player, err := s.playerRepo.GetByToken(ctx, tokenHashHex)
	if err != nil {
		if errors.Is(err, e.ErrNotFound) {
			return nil, e.NewUnauthorizedError("invalid token")
		}
		return nil, e.NewInternalError("authentication failed", err)
	}

	return player, nil
}

func (s *AuthService) AuthenticatePlayerWeb(ctx context.Context, req dto.LoginRequest) (*domain.Player, error) {
	player, err := s.playerRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, e.NewInternalError("", err)
	}
	if player == nil {
		return nil, e.NewValidationError("пользователя с таким логином не существует")
	}

	// ** TEMP for restoring password ** //
	if player.PasswordHash == "" {
		if valid, message := s.validatePlayerPassword(req.LoginData); !valid {
			return nil, e.NewValidationError(message)
		}
		newPWHash, err := s.generateHash(req.LoginData)
		if err != nil {
			return nil, e.NewInternalError("", err)
		}
		updatedPlayer, err := s.playerRepo.SetPlayerPassword(ctx, player.ID, newPWHash)
		if err != nil {
			return nil, e.NewInternalError("не удалось изменить пароль", err)
		}
		player.PasswordHash = updatedPlayer.PasswordHash
	} else {
		valid, err := s.validateHash(req.LoginData, player.PasswordHash)
		if err != nil || !valid {
			return nil, e.NewUnauthorizedError("Неверный пароль")
		}
	}

	return player, nil
}

func (s *AuthService) CreatePlayer(ctx context.Context, signup dto.SignUpRequest) (*domain.Player, error) {
	if valid, message := s.validatePlayerUsername(signup.Username); !valid {
		return nil, e.NewValidationError(message)
	}

	if available, err := s.checkUsernameAvailability(ctx, signup.Username); err != nil {
		return nil, e.NewInternalError("не удалось проверить доступность имени пользователя", err)
	} else if !available {
		return nil, e.NewValidationError("Логин занят")
	}

	if valid, message := s.validatePlayerPassword(signup.Password); !valid {
		return nil, e.NewValidationError(message)
	}

	if _, err := mail.ParseAddress(signup.Email); err != nil {
		return nil, e.NewValidationError("Неверный формат почты")
	}

	hash, err := s.generateHash(signup.Password)
	if err != nil {
		return nil, e.NewInternalError("не удалось сгенерировать хэш пароля", err)
	}

	player, err := s.playerRepo.CreatePlayer(ctx, &domain.Player{
		Username:     signup.Username,
		PasswordHash: hash,
		Email:        signup.Email,
	})
	if err != nil {
		return nil, e.NewInternalError("не удалось создать игрока", err)
	}

	return player, nil
}

func (s *AuthService) GetLoginUsernameAvailability(ctx context.Context, username string) (bool, error) {
	if valid, message := s.validatePlayerUsername(username); !valid {
		return false, e.NewValidationError(message)
	}

	available, err := s.checkUsernameAvailability(ctx, username)
	if err != nil {
		return false, e.NewInternalError("ошибка проверки имени пользователя", err)
	}

	return available, nil
}

func (s *AuthService) EmitWebToken(ctx context.Context, player *domain.Player) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.authConfig.JWTTokenLifetimeHours) * time.Hour)

	tokenString := s.GeneratePlayerToken()
	tokenHash := sha256.Sum256([]byte(tokenString))
	tokenHashHex := hex.EncodeToString(tokenHash[:])

	dbToken := &domain.Token{
		PlayerID:  player.ID,
		TokenHash: tokenHashHex,
		ExpiresAt: expirationTime,
		Revoked:   false,
	}

	if _, err := s.playerRepo.InsertToken(ctx, dbToken); err != nil {
		return "", e.NewInternalError("ошибка создания токена", err)
	}

	return tokenString, nil
}

func (s *AuthService) validatePlayerUsername(username string) (bool, string) {
	if valid, count := validation.IsValidLength(username, 6, 30); !valid {
		if count > 0 {
			return false, "Логин не может быть больше 30 символов"
		} else if count < 0 {
			return false, "Логин не может быть меньше 6 символов"
		}
	}

	if valid := validation.IsValidString(username, true, true, []rune{'.', '-', '_'}); !valid {
		return false, "Логин может содержать только латинские буквы, цифры, точку, дефис или нижнее подчёркивание"
	}

	return true, "Логин валиден"
}

func (s *AuthService) validatePlayerPassword(password string) (bool, string) {
	if valid, count := validation.IsValidLength(password, 8, 64); !valid {
		if count > 0 {
			return false, "Пароль не может быть больше 64 символов"
		} else if count < 0 {
			return false, "Пароль не может быть меньше 8 символов"
		}
	}

	if valid := validation.IsValidString(password, true, true, []rune{'.', ',', '-', '_', '!', '@', '#', '$', '%', '^', '&', '*', '+', '=', '?', '/', '\\'}); !valid {
		return false, "Пароль содержит некорректные символы - возможно введены не латинские буквы или скобки"
	}

	return true, "Пароль валиден"
}

func (s *AuthService) checkUsernameAvailability(ctx context.Context, usernameToCheck string) (bool, error) {
	return s.playerRepo.IsUsernameFree(ctx, usernameToCheck)
}

func (s *AuthService) GeneratePlayerToken() string {
	data := make([]byte, 16)
	rand.Read(data[0:16])

	return hex.EncodeToString(data)
}

func (s *AuthService) generateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) validateHash(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

// // chat in @chat format
// func (s *AuthService) checkTGUserChatMembership(userID int64, chat string, botToken string) (bool, error) {
// 	url := fmt.Sprintf(
// 		"https://api.telegram.org/bot%s/getChatMember?chat_id=%s&user_id=%d",
// 		botToken, chat, userID,
// 	)
//
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer resp.Body.Close()
//
// 	var data map[string]any
// 	json.NewDecoder(resp.Body).Decode(&data)
//
// 	if !data["ok"].(bool) {
// 		return false, fmt.Errorf("API error: %s", data["description"])
// 	}
//
// 	result := data["result"].(map[string]any)
// 	status := result["status"].(string)
//
// 	return status == "creator" || status == "administrator" ||
// 		status == "member" || status == "restricted", nil
// }
