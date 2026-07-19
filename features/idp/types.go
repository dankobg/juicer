package idp

import (
	"time"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/google/uuid"
)

type ListUsersFilters struct {
	Page     *api.PaginationPage
	PageSize *api.PaginationPageSize
	Sort     *[]string
}

type ListFriendsFilters struct {
	api.ListFriendsParams
}

type ListFollowingsFilters struct {
	api.ListFollowingsParams
}

type ListBlockedUsersFilters struct {
	api.ListBlockedUsersParams
}

type ListFriendRequestsFilters struct {
	api.ListFriendRequestsParams
}

type FriendRequest struct {
	UserID     uuid.UUID
	FriendID   uuid.UUID
	Status     api.FriendRequestStatus
	CreatedAt  time.Time
	AnsweredAt *time.Time
}
