package idp

import (
	"context"

	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/pagination"
	"github.com/google/uuid"
)

type UserPersistor interface {
	ListUsers(ctx context.Context, filters ListUsersFilters) (pagination.WithTotal[models.User], error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	CreateUser(ctx context.Context, in models.UserSetter) (models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, in models.UserSetter) (models.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)

	CreateFriendRequest(ctx context.Context, in models.FriendshipSetter) (FriendRequest, error)
	ListFriendRequests(ctx context.Context, filters ListFriendRequestsFilters) (pagination.WithTotal[FriendRequest], error)
	AcceptFriendRequest(ctx context.Context, userID, friendID uuid.UUID) error
	DeclineFriendRequest(ctx context.Context, userID, friendID uuid.UUID) error
	CancelFriendRequest(ctx context.Context, userID, friendID uuid.UUID) error
	ListFriends(ctx context.Context, userID uuid.UUID, filters ListFriendsFilters) (pagination.WithTotal[models.User], error)
	GetFriend(ctx context.Context, userID, friendID uuid.UUID) (models.User, error)
	AddFriend(ctx context.Context, in models.FriendshipSetter) error
	DeleteFriend(ctx context.Context, userID, friendID uuid.UUID) error
	DeleteFriends(ctx context.Context, userID uuid.UUID, friendIDs []uuid.UUID) error

	ListFollowings(ctx context.Context, userID uuid.UUID, filters ListFollowingsFilters) (pagination.WithTotal[models.User], error)
	GetFollowing(ctx context.Context, userID, followingID uuid.UUID) (models.User, error)
	FollowUser(ctx context.Context, userID uuid.UUID, in models.FollowingSetter) error
	ListBlockedUsers(ctx context.Context, userID uuid.UUID, filters ListBlockedUsersFilters) (pagination.WithTotal[models.User], error)
	GetBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) (models.User, error)
	BlockUser(ctx context.Context, userID uuid.UUID, in models.BlocklistSetter) error
}
