package store

import (
	"context"
	"errors"
	"github.com/JohnKucharsky/DevGroup/domain"
	"github.com/JohnKucharsky/DevGroup/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type UserStore struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewUserStore(db *pgxpool.Pool, redis *redis.Client) *UserStore {
	return &UserStore{
		db:    db,
		redis: redis,
	}
}

func (us *UserStore) Create(u domain.SignUpInput) (*domain.User, error) {
	ctx := context.Background()

	rows, err := us.db.Query(
		ctx, `
        INSERT INTO users (id, name, email, password)
        VALUES (@id, @name, @email, @password)
        RETURNING id, name, email, password, updated_at`,
		pgx.NamedArgs{
			"id":       uuid.Must(uuid.NewV7()),
			"name":     u.Name,
			"email":    u.Email,
			"password": u.Password,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[domain.User],
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (us *UserStore) SetAccessToken(id uuid.UUID) (*string, error) {
	ctx := context.TODO()
	now := time.Now()

	var accessTokenExpiresInString = os.Getenv("ACCESS_TOKEN_EXPIRED_IN")
	var accessTokenPrivateKey = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessTokenExpiresIn, _ := time.ParseDuration(accessTokenExpiresInString)

	accessTokenDetails, err := utils.CreateToken(
		id.String(),
		accessTokenExpiresIn,
		accessTokenPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	if err = us.redis.Set(
		ctx,
		accessTokenDetails.TokenUUID,
		id.String(),
		time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(now),
	).Err(); err != nil {
		return nil, err
	}

	return accessTokenDetails.Token, nil
}

func (us *UserStore) SetRefreshToken(id uuid.UUID) (*string, error) {
	ctx := context.TODO()
	now := time.Now()

	var refreshTokenExpiresInString = os.Getenv("REFRESH_TOKEN_EXPIRED_IN")
	var refreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")

	refreshTokenExpiresIn, _ := time.ParseDuration(refreshTokenExpiresInString)

	refreshTokenDetails, err := utils.CreateToken(
		id.String(),
		refreshTokenExpiresIn,
		refreshTokenPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	if err = us.redis.Set(
		ctx,
		refreshTokenDetails.TokenUUID,
		id.String(),
		time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(now),
	).Err(); err != nil {
		return nil, err
	}

	return refreshTokenDetails.Token, nil
}

func (us *UserStore) GetOne(email, id string) (*domain.User, error) {
	ctx := context.Background()

	var query, param string
	if email == "" {
		query = `select id,name,email,password,updated_at from users where id::text = @param`
		param = id
	} else if id == "" {
		query = `select id,name,email,password,updated_at from users where email = @param`
		param = email
	} else {
		return nil, errors.New("you have to provide either email or id to run this thing")
	}

	rows, err := us.db.Query(
		ctx, query,
		pgx.NamedArgs{"param": param},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.User],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (us *UserStore) GetByRefreshTokenRedis(token string) (string, error) {
	ctx := context.TODO()

	var refreshTokenPublicKey = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	tokenClaims, err := utils.ValidateToken(
		token,
		refreshTokenPublicKey,
	)
	if err != nil {
		return "", err
	}

	userID, err := us.redis.Get(
		ctx,
		tokenClaims.TokenUUID,
	).Result()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (us *UserStore) GetByAccessTokenRedis(token string) (
	string,
	string,
	error,
) {
	ctx := context.TODO()

	var accessTokenPublicKey = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	tokenClaims, err := utils.ValidateToken(
		token,
		accessTokenPublicKey,
	)
	if err != nil {
		return "", "", err
	}

	userID, err := us.redis.Get(
		ctx,
		tokenClaims.TokenUUID,
	).Result()
	if errors.Is(err, redis.Nil) {
		return "", "", err
	}

	return userID, tokenClaims.TokenUUID, nil
}

func (us *UserStore) DeleteTokensRedis(refreshToken, accessToken string) error {
	ctx := context.TODO()

	var refreshTokenPublicKey = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	tokenClaims, err := utils.ValidateToken(
		refreshToken,
		refreshTokenPublicKey,
	)
	if err != nil {
		return err
	}

	_, err = us.redis.Del(
		ctx,
		tokenClaims.TokenUUID,
		accessToken,
	).Result()
	if err != nil {
		return err
	}

	return nil
}
