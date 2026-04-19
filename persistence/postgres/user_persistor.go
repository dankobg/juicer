package postgres

import (
	"context"
	"fmt"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/persistence"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

var _ persistence.UserPersistor = (*PgUserPersistor)(nil)

type PgUserPersistor struct {
	*PgPersistor
}

func NewPgUserPersistor(ps *PgPersistor) *PgUserPersistor {
	return &PgUserPersistor{
		PgPersistor: ps,
	}
}

func (pst *PgUserPersistor) ListUsers(ctx context.Context, filters dbtype.ListUsersFilters) (dbtype.PagedResult[models.User], error) {
	q := psql.Select(
		sm.Columns(models.Users.Columns),
		sm.From(models.Users.Name()),
		sm.GroupBy(models.Users.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Users.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	type ListUsersRow struct {
		models.User
		TotalCount int64
	}

	countries, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListUsersRow]())
	if err != nil {
		return dbtype.PagedResult[models.User]{}, fmt.Errorf("query users")
	}

	result := dbtype.PagedResult[models.User]{
		Data: make([]models.User, len(countries)),
	}
	for i, row := range countries {
		result.Data[i] = row.User
	}

	if len(countries) > 0 {
		result.TotalCount = countries[0].TotalCount
	}

	return result, nil
}

func (pst *PgUserPersistor) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	q := psql.Select(
		sm.Columns(models.Users.Columns),
		sm.From(models.Users.Name()),
		sm.Where(models.Users.Columns.ID.EQ(psql.Arg(userID))),
	)

	user, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.User]())
	if err != nil {
		return models.User{}, fmt.Errorf("query user")
	}

	return user, nil
}

func (pst *PgUserPersistor) DeleteUserByID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	q := models.Users.Delete(dm.Where(models.Users.Columns.ID.EQ(psql.Arg(userID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return uuid.Nil, fmt.Errorf("delete country: %w", err)
	}

	return userID, nil
}

func (pst *PgUserPersistor) CreateUser(ctx context.Context, in models.UserSetter) (models.User, error) {
	q := models.Users.Insert(&in, im.Returning(models.Users.Columns))

	user, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.User]())
	if err != nil {
		return models.User{}, fmt.Errorf("insert user")
	}

	return user, nil
}

func (pst *PgUserPersistor) UpdateUser(ctx context.Context, userID uuid.UUID, in models.UserSetter) (models.User, error) {
	q := models.Users.Update(
		in.UpdateMod(),
		um.Where(models.Users.Columns.ID.EQ(psql.Arg(userID))),
		um.Returning(models.Users.Columns),
	)

	user, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.User]())
	if err != nil {
		return models.User{}, fmt.Errorf("update user")
	}

	return user, nil
}

func (pst *PgUserPersistor) CreateFriendRequest(ctx context.Context, in models.FriendshipSetter) (dbtype.FriendRequest, error) {
	q := models.Friendships.Insert(&in, im.Returning(models.Friendships.Columns))

	friendship, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Friendship]())
	if err != nil {
		return dbtype.FriendRequest{}, fmt.Errorf("insert friend request")
	}

	friendRequest := dbtype.FriendRequest{
		UserID:     friendship.InitiatorID,
		FriendID:   friendship.ReceiverID,
		Status:     api.FriendRequestStatus(friendship.Status),
		CreatedAt:  friendship.CreatedAt,
		AnsweredAt: friendship.AnsweredAt.Ptr(),
	}

	return friendRequest, nil
}

func (pst *PgUserPersistor) ListFriendRequests(ctx context.Context, filters dbtype.ListFriendRequestsFilters) (dbtype.PagedResult[dbtype.FriendRequest], error) {
	panic("@TODO IMPLEMENT DB - ListFriendRequests")
}

func (pst *PgUserPersistor) AcceptFriendRequest(ctx context.Context, userID, friendID uuid.UUID) error {
	panic("@TODO IMPLEMENT DB - AcceptFriendRequest")
}

func (pst *PgUserPersistor) DeclineFriendRequest(ctx context.Context, userID, friendID uuid.UUID) error {
	panic("@TODO IMPLEMENT DB - DeclineFriendRequest")
}

func (pst *PgUserPersistor) CancelFriendRequest(ctx context.Context, userID, friendID uuid.UUID) error {
	panic("@TODO IMPLEMENT DB - CancelFriendRequest")
}

func (pst *PgUserPersistor) ListFriends(ctx context.Context, userID uuid.UUID, filters dbtype.ListFriendsFilters) (dbtype.PagedResult[models.User], error) {
	panic("@TODO IMPLEMENT DB - ListFriends")
}

func (pst *PgUserPersistor) GetFriend(ctx context.Context, userID, friendID uuid.UUID) (models.User, error) {
	panic("@TODO IMPLEMENT DB - GetFriend")
}

func (pst *PgUserPersistor) AddFriend(ctx context.Context, in models.FriendshipSetter) error {
	panic("@TODO IMPLEMENT DB - AddFriend")
}

func (pst *PgUserPersistor) DeleteFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	panic("@TODO IMPLEMENT DB - DeleteFriend")
}

func (pst *PgUserPersistor) DeleteFriends(ctx context.Context, userID uuid.UUID, friendIDs []uuid.UUID) error {
	panic("@TODO IMPLEMENT DB - DeleteFriends")
}

func (pst *PgUserPersistor) ListFollowings(ctx context.Context, userID uuid.UUID, filters dbtype.ListFollowingsFilters) (dbtype.PagedResult[models.User], error) {
	panic("@TODO IMPLEMENT DB - ListFollowings")
}

func (pst *PgUserPersistor) GetFollowing(ctx context.Context, userID, followingID uuid.UUID) (models.User, error) {
	panic("@TODO IMPLEMENT DB - GetFollowing")
}

func (pst *PgUserPersistor) FollowUser(ctx context.Context, userID uuid.UUID, in models.FollowingSetter) error {
	panic("@TODO IMPLEMENT DB - FollowUser")
}

func (pst *PgUserPersistor) ListBlockedUsers(ctx context.Context, userID uuid.UUID, filters dbtype.ListBlockedUsersFilters) (dbtype.PagedResult[models.User], error) {
	panic("@TODO IMPLEMENT DB - ListBlockedUsers")
}

func (pst *PgUserPersistor) GetBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) (models.User, error) {
	panic("@TODO IMPLEMENT DB - GetBlockedUser")
}

func (pst *PgUserPersistor) BlockUser(ctx context.Context, userID uuid.UUID, in models.BlocklistSetter) error {
	panic("@TODO IMPLEMENT DB - BlockUser")
}
