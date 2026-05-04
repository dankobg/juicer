package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/engine"
	"github.com/dankobg/juicer/gameplay"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/dankobg/juicer/persistence"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/dankobg/juicer/ws"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type clientAuthInfo struct {
	userID    string
	connID    string
	authState ws.ClientAuthState
}

func (a *ApiHandler) PubsubProcess(ctx context.Context) {
	a.Log.Info("gameserver pubsub started")

	for {
		select {
		case msg := <-a.bus.subMessages["ipc"]:
			a.onIPCMsg(msg)
		case msg := <-a.bus.subMessages["wsc.*"]:
			a.onWSCMsg(msg)

		case <-ctx.Done():
			a.Log.Debug("gameserver pubsub ctx done")
			return
		}
	}
}

func (a *ApiHandler) onIPCMsg(m *redis.Message) {
	msg := &pb.Message{}
	if err := protojson.Unmarshal([]byte(m.Payload), msg); err != nil {
		a.Log.Error("protojson.Unmarshal IPC Message")
		return
	}

	switch msg.GetEvent().(type) {
	case *pb.Message_ClientConnected:
		a.handleIPCClientConnectedMsg(msg.GetClientConnected())

	case *pb.Message_ClientDisconnected:
		a.handleIPCClientDisconnectedMsg(msg.GetClientDisconnected())

	case *pb.Message_Heartbeat:
		a.handleIPCHeartbeatMsg(msg.GetHeartbeat())

	case *pb.Message_LeaveTab:
		a.handleIPCLeaveTabMsg(msg.GetLeaveTab())

	case *pb.Message_LeaveSite:
		a.handleIPCLeaveSiteMsg(msg.GetLeaveSite())

	case *pb.Message_InitializeChannels:
		a.handleIPCInitializeChannelsMsg(msg.GetInitializeChannels())
	}
}

func (a *ApiHandler) handleIPCHeartbeatMsg(data *pb.Heartbeat) {
	var username string

	if data.Guest {
		username = "guest-" + data.UserId
	} else {
		uname, err := a.GetUsername(context.Background(), data.UserId)
		if err != nil {
			a.Log.Error("handleIPCHeartbeatMsg get username", slog.Any("error", err))
			return
		}

		username = uname
	}

	err := a.persistor.Presence().RefreshPresence(context.Background(), uuid.MustParse(data.UserId), uuid.MustParse(data.ConnId), username, data.Guest)
	if err != nil && !errors.Is(err, redis.Nil) {
		a.Log.Error("RefreshPresence", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
	}

	// refresh presence
	// broadcast presence change
}

func (a *ApiHandler) handleIPCLeaveTabMsg(data *pb.LeaveTab) {
	var username string

	if data.Guest {
		username = "guest-" + data.UserId
	} else {
		uname, err := a.GetUsername(context.Background(), data.UserId)
		if err != nil {
			a.Log.Error("handleIPCLeaveTabMsg get username", slog.Any("error", err))
			return
		}

		username = uname
	}

	channelsDiff, err := a.persistor.Presence().ClearPresence(context.Background(), uuid.MustParse(data.UserId), uuid.MustParse(data.ConnId), username, data.Guest)
	if err != nil && !errors.Is(err, redis.Nil) {
		a.Log.Error("ClearPresence", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
	}

	if err := a.publishPresenceDiff(context.Background(), channelsDiff, data.UserId, data.ConnId, username, data.Guest); err != nil {
		a.Log.Error("broadcastPresenceDiff", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
	}

	for _, leftChannel := range channelsDiff.UserLeft {
		if err := a.sendUserPresenceDiffToChannel(context.Background(), channelsDiff, leftChannel, data.UserId, data.ConnId, username, data.Guest); err != nil {
			a.Log.Error("sendUserPresenceDiffToChannel", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
		}
	}
}

func (a *ApiHandler) handleIPCLeaveSiteMsg(data *pb.LeaveSite) {
}

func (a *ApiHandler) handleIPCInitializeChannelsMsg(data *pb.InitializeChannels) {
	channels := make([]string, 0)

	if data.Path == "" || data.Path == "/" {
		channels = append(channels, "lobby", "lobby.chat")
	} else {
		var gameID int64

		if gameIDStr, ok := strings.CutPrefix(data.Path, "/game/"); ok {
			n, err := strconv.ParseInt(gameIDStr, 10, 64)
			if err != nil {
				a.Log.Error("gameid parseint", slog.Any("error", err))
				return
			}

			gameID = n
		}

		if gametvIDStr, ok := strings.CutPrefix(data.Path, "/gametv/"); ok {
			n, err := strconv.ParseInt(gametvIDStr, 10, 64)
			if err != nil {
				a.Log.Error("gametvid parseint", slog.Any("error", err))
				return
			}

			gameID = n
		}

		game, err := a.persistor.Game().GetGameByID(context.Background(), gameID, dbtype.GetGameByIDFilters{})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			a.Log.Error("handleIPCRequestInitialChannelsMsg GetGameByID", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.String("path", data.Path), slog.Any("error", err))
			return
		}

		if game.ID != 0 {
			switch game.GameStateID {
			// @TODO: fix later...
			case a.gameStateProtoToID(pb.GameState_GAME_STATE_WAITING_START),
				a.gameStateProtoToID(pb.GameState_GAME_STATE_IN_PROGRESS):
				channels = append(channels, fmt.Sprintf("game.%d", game.ID), fmt.Sprintf("game.%d.chat", game.ID))
			case a.gameStateProtoToID(pb.GameState_GAME_STATE_FINISHED),
				a.gameStateProtoToID(pb.GameState_GAME_STATE_INTERRUPTED):
				channels = append(channels, fmt.Sprintf("gametv.%d", game.ID), fmt.Sprintf("gametv.%d.chat", game.ID))
			default:
				return
			}
		}
	}

	initialChannelsMsg := &pb.Message{
		Event: &pb.Message_InitialChannels{InitialChannels: &pb.InitialChannels{
			Channels: channels,
		}},
	}

	initialChannelsMsgBytes, err := protojson.Marshal(initialChannelsMsg)
	if err != nil {
		a.Log.Error("protojson marshal Message_InitialChannels", slog.String("user_id", data.UserId), slog.Any("error", err))
		return
	}

	topic := "reply-initial-channels." + data.UserId + "." + data.ConnId
	if err := a.bus.rdb.Publish(context.Background(), topic, initialChannelsMsgBytes).Err(); err != nil {
		a.Log.Error("hub publish Message_InitialChannels", slog.String("user_id", data.UserId), slog.String("topic", "ipc"), slog.Any("error", err))
		return
	}
}

func (a *ApiHandler) handleIPCClientConnectedMsg(data *pb.ClientConnected) {
	var username string

	if data.Guest {
		username = "guest-" + data.UserId
	} else {
		uname, err := a.GetUsername(context.Background(), data.UserId)
		if err != nil {
			a.Log.Error("handleIPCClientConnectedMsg get username", slog.Any("error", err))
			return
		}

		username = uname
	}

	channels := data.GetChannels()

	channelsDiff, err := a.persistor.Presence().SetPresence(context.Background(), uuid.MustParse(data.UserId), uuid.MustParse(data.ConnId), username, data.Guest, channels)
	if err != nil && !errors.Is(err, redis.Nil) {
		a.Log.Error("SetPresence", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
		return
	}

	if err := a.publishPresenceDiff(context.Background(), channelsDiff, data.UserId, data.ConnId, username, data.Guest); err != nil {
		a.Log.Error("broadcastPresenceDiff", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
	}

	for _, channel := range channels {
		if err := a.sendChannelPresenceStateToConn(context.Background(), channel, data.ConnId); err != nil {
			a.Log.Error("sendChannelPresenceStateToConn", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
		}

		if err := a.sendUserPresenceDiffToChannel(context.Background(), channelsDiff, channel, data.UserId, data.ConnId, username, data.Guest); err != nil {
			a.Log.Error("sendUserPresenceDiffToChannel", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
		}

		if channel == "lobby" {
			if err := a.sendLobbyInfo(data.UserId, data.ConnId, data.Guest); err != nil {
				a.Log.Error("sendLobbyInfo", slog.Any("error", err))
			}
		}

		if channel == "lobby.chat" {
			if err := a.sendLobbyChatInfo(data.UserId, data.ConnId, data.Guest); err != nil {
				a.Log.Error("sendLobbyChatInfo", slog.Any("error", err))
			}
		}

		if after, foundPrefix := strings.CutPrefix(channel, "game."); foundPrefix {
			if before, foundSuffix := strings.CutSuffix(after, ".chat"); foundSuffix {
				fmt.Println("SEND GAME CHAT FOR GAME:", before)
			} else {
				gameID, err := strconv.ParseInt(after, 10, 64)

				if err != nil {
					a.Log.Error("gameid parseint", slog.Any("error", err))
					return
				}

				if err := a.sendGameInfo(gameID, data.UserId, data.ConnId, data.Guest); err != nil {
					a.Log.Error("sendGameInfo", slog.Any("error", err))
				}
			}
		}

		if after, foundPrefix := strings.CutPrefix(channel, "gametv."); foundPrefix {
			if before, foundSuffix := strings.CutSuffix(after, ".chat"); foundSuffix {
				fmt.Println("SEND GAMETV CHAT FOR GAME:", before)
			} else {
				gameID, err := strconv.ParseInt(after, 10, 64)
				if err != nil {
					a.Log.Error("gametvid parseint", slog.Any("error", err))
					return
				}

				if err := a.sendGameTvInfo(gameID, data.UserId, data.ConnId, data.Guest); err != nil {
					a.Log.Error("sendGameTvInfo", slog.Any("error", err))
				}
			}
		}
	}
}

func (a *ApiHandler) handleIPCClientDisconnectedMsg(data *pb.ClientDisconnected) {}

func (a *ApiHandler) onWSCMsg(m *redis.Message) {
	msg := &pb.Message{}
	if err := protojson.Unmarshal([]byte(m.Payload), msg); err != nil {
		a.Log.Error("protojson.Unmarshal WSC Message")
		return
	}

	clientAuthInfo, err := extractWSCTopicParts(m.Channel)
	if err != nil {
		a.Log.Error("extractWSCTopicParts", slog.String("channel", m.Channel), slog.String("pattern", m.Pattern), slog.String("payload", m.Payload), slog.Any("error", err))
		return
	}

	switch msg.GetEvent().(type) {
	case *pb.Message_Echo:
		a.handleWSCEchoMsg(clientAuthInfo, msg.GetEcho())
	case *pb.Message_SeekGame:
		a.handleWSCSeekGameMsg(clientAuthInfo, msg.GetSeekGame())
	case *pb.Message_CancelSeekGame:
		a.handleWSCCancelSeekGameMsg(clientAuthInfo, msg.GetCancelSeekGame())
	case *pb.Message_AbortGame:
		a.handleWSCAbortGame(clientAuthInfo, msg.GetAbortGame())
	case *pb.Message_ResignGame:
		a.handleWSCResignGame(clientAuthInfo, msg.GetResignGame())
	case *pb.Message_OfferDraw:
		a.handleWSCOfferDraw(clientAuthInfo, msg.GetOfferDraw())
	case *pb.Message_AcceptDraw:
		a.handleWSCAcceptDraw(clientAuthInfo, msg.GetAcceptDraw())
	case *pb.Message_DeclineDraw:
		a.handleWSCDeclineDraw(clientAuthInfo, msg.GetDeclineDraw())
	case *pb.Message_PlayMoveUci:
		a.handleWSCPlayMoveUCI(clientAuthInfo, msg.GetPlayMoveUci())
	case *pb.Message_SendLobbyChat:
		a.handleWSCSendLobbyChat(clientAuthInfo, msg.GetSendLobbyChat())
	}
}

func (a *ApiHandler) handleWSCEchoMsg(authInfo clientAuthInfo, data *pb.Echo) {
	bb, _ := protojson.Marshal(&pb.Message{Event: &pb.Message_Echo{Echo: &pb.Echo{Message: strings.ToUpper(data.Message)}}})
	toUser, toConn, toLobby := "user."+authInfo.userID, "conn."+authInfo.connID, "lobby.chat"
	_ = []any{toUser, toConn, toLobby}
	a.bus.rdb.Publish(context.Background(), toUser, bb)
}

func (a *ApiHandler) handleWSCSeekGameMsg(authInfo clientAuthInfo, data *pb.SeekGame) {
	if data == nil {
		return
	}

	pool := dbtype.Pool{
		ClockMS:     data.GetTimeControl().GetClockMs(),
		IncrementMS: data.GetTimeControl().GetIncrementMs(),
		Rated:       authInfo.authState == ws.ClientAuth,
	}
	if err := a.persistor.Pool().JoinPool(context.Background(), uuid.MustParse(authInfo.userID), pool); err != nil {
		a.Log.Error("SeekGame join pool failed", slog.String("pool", pool.Name()), slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
		return
	}
}

func (a *ApiHandler) handleWSCCancelSeekGameMsg(authInfo clientAuthInfo, data *pb.CancelSeekGame) {
	if err := a.persistor.Pool().LeavePool(context.Background(), uuid.MustParse(authInfo.userID)); err != nil {
		a.Log.Error("CancelSeekGame leave pool failed", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
		return
	}
}

func (a *ApiHandler) handleWSCAbortGame(authInfo clientAuthInfo, data *pb.AbortGame) {
	fmt.Println(data, "handleWSCAbortGame")
}

func (a *ApiHandler) handleWSCResignGame(authInfo clientAuthInfo, data *pb.ResignGame) {
	fmt.Println(data, "handleWSCResignGame")
}

func (a *ApiHandler) handleWSCOfferDraw(authInfo clientAuthInfo, data *pb.OfferDraw) {
	fmt.Println(data, "handleWSCOfferDraw")
}

func (a *ApiHandler) handleWSCAcceptDraw(authInfo clientAuthInfo, data *pb.AcceptDraw) {
	fmt.Println(data, "handleWSCAcceptDraw")
}

func (a *ApiHandler) handleWSCDeclineDraw(authInfo clientAuthInfo, data *pb.DeclineDraw) {
	fmt.Println(data, "handleWSCDeclineDraw")
}

func (a *ApiHandler) handleWSCPlayMoveUCI(authInfo clientAuthInfo, data *pb.PlayMoveUCI) {
	fmt.Println(data, "handleWSCPlayMoveUCI")
}

func (a *ApiHandler) handleWSCSendLobbyChat(authInfo clientAuthInfo, data *pb.SendLobbyChat) {
	fmt.Println(data, "handleWSCSendLobbyChat")
	lobbyChatReceivedMsg := &pb.Message{Event: &pb.Message_LobbyChat{LobbyChat: &pb.LobbyChat{
		MessageId: "1",
		UserId:    authInfo.userID,
		PostedAt:  time.Now().Format(time.RFC3339),
		Message:   data.Message,
	}}}

	lobbyChatReceivedMsgBytes, err := protojson.Marshal(lobbyChatReceivedMsg)
	if err != nil {
	} else {
		if err := a.bus.rdb.Publish(context.Background(), "lobby.chat", lobbyChatReceivedMsgBytes).Err(); err != nil {
			a.Log.Error("LobbyChat publish", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
			return
		}
	}
}

// extractWSCTopicParts extracts the user_id, conn_id and auth_state
func extractWSCTopicParts(topic string) (clientAuthInfo, error) {
	parts := strings.Split(topic, ".")
	if len(parts) != 4 {
		return clientAuthInfo{}, fmt.Errorf("invalid parts length, expected 4, got: %d", len(parts))
	}

	clientID, connID, authStateStr := parts[1], parts[2], parts[3]
	if authStateStr != "0" && authStateStr != "1" {
		return clientAuthInfo{}, fmt.Errorf("invalid parts auth_state. must be 0 or 1")
	}

	authState := ws.ClientGuest
	if authStateStr == "1" {
		authState = ws.ClientAuth
	}

	return clientAuthInfo{
		userID:    clientID,
		connID:    connID,
		authState: authState,
	}, nil
}

func (a *ApiHandler) FetchCategoryThresholds(ctx context.Context) error {
	gameTimeCategories, err := a.persistor.GameTimeCategory().ListGameTimeCategories(ctx, dbtype.ListGameTimeCategoriesFilters{})
	if err != nil {
		return err
	}

	for _, v := range gameTimeCategories.Data {
		var limit time.Duration = math.MaxUint32
		if v.UpperTimeLimitSecs.IsValue() {
			limit = time.Second * time.Duration(v.UpperTimeLimitSecs.MustGet())
		}

		switch v.Name {
		case "hyperbullet":
			a.categoryThresholds = append(a.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_HYPERBULLET, upperLimit: limit})
		case "bullet":
			a.categoryThresholds = append(a.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_BULLET, upperLimit: limit})
		case "blitz":
			a.categoryThresholds = append(a.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_BLITZ, upperLimit: limit})
		case "rapid":
			a.categoryThresholds = append(a.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_RAPID, upperLimit: limit})
		case "classical":
			a.categoryThresholds = append(a.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_CLASSICAL, upperLimit: limit})
		}
	}

	return nil
}

func (a *ApiHandler) FetchProtoMappingsCacheLookups(ctx context.Context) error {
	gameVariants, e1 := a.persistor.GameVariant().ListGameVariants(ctx, dbtype.ListGameVariantsFilters{})
	gameTimeKinds, e2 := a.persistor.GameTimeKind().ListGameTimeKinds(ctx, dbtype.ListGameTimeKindsFilters{})
	gameTimeCategories, e3 := a.persistor.GameTimeCategory().ListGameTimeCategories(ctx, dbtype.ListGameTimeCategoriesFilters{})
	gameResults, e4 := a.persistor.GameResult().ListGameResults(ctx, dbtype.ListGameResultsFilters{})
	gameResultStatuses, e5 := a.persistor.GameResultStatus().ListGameResultStatuses(ctx, dbtype.ListGameResultStatusesFilters{})

	gameStates, e6 := a.persistor.GameState().ListGameStates(ctx, dbtype.ListGameStatesFilters{})
	if err := errors.Join(e1, e2, e3, e4, e5, e6); err != nil {
		return err
	}

	gameVariantNameToProto := map[string]pb.GameVariant{
		"standard":         pb.GameVariant_GAME_VARIANT_STANDARD,
		"atomic":           pb.GameVariant_GAME_VARIANT_ATOMIC,
		"crazyhouse":       pb.GameVariant_GAME_VARIANT_CRAZYHOUSE,
		"chess960":         pb.GameVariant_GAME_VARIANT_CHESS960,
		"king-of-the-hill": pb.GameVariant_GAME_VARIANT_KING_OF_THE_HILL,
		"three-check":      pb.GameVariant_GAME_VARIANT_THREE_CHECK,
		"horde":            pb.GameVariant_GAME_VARIANT_HORDE,
		"racing-kings":     pb.GameVariant_GAME_VARIANT_RACING_KINGS,
	}

	gameTimeKindNameToProto := map[string]pb.GameTimeKind{
		"realtime":       pb.GameTimeKind_GAME_TIME_KIND_REALTIME,
		"correspondence": pb.GameTimeKind_GAME_TIME_KIND_CORRESPONDENCE,
		"unlimited":      pb.GameTimeKind_GAME_TIME_KIND_UNLIMITED,
	}

	gameTimeCategoryNameToProto := map[string]pb.GameTimeCategory{
		"hyperbullet": pb.GameTimeCategory_GAME_TIME_CATEGORY_HYPERBULLET,
		"bullet":      pb.GameTimeCategory_GAME_TIME_CATEGORY_BULLET,
		"blitz":       pb.GameTimeCategory_GAME_TIME_CATEGORY_BLITZ,
		"rapid":       pb.GameTimeCategory_GAME_TIME_CATEGORY_RAPID,
		"classical":   pb.GameTimeCategory_GAME_TIME_CATEGORY_CLASSICAL,
	}

	gameResultNameToProto := map[string]pb.GameResult{
		"white-won":   pb.GameResult_GAME_RESULT_WHITE_WON,
		"black-won":   pb.GameResult_GAME_RESULT_BLACK_WON,
		"draw":        pb.GameResult_GAME_RESULT_DRAW,
		"interrupted": pb.GameResult_GAME_RESULT_INTERRUPTED,
	}

	gameResultStatusNameToProto := map[string]pb.GameResultStatus{
		"checkmate":             pb.GameResultStatus_GAME_RESULT_STATUS_CHECKMATE,
		"insufficient-material": pb.GameResultStatus_GAME_RESULT_STATUS_INSUFFICIENT_MATERIAL,
		"threefold-repetition":  pb.GameResultStatus_GAME_RESULT_STATUS_THREEFOLD_REPETITION,
		"fivefold-repetition":   pb.GameResultStatus_GAME_RESULT_STATUS_FIVEFOLD_REPETITION,
		"fifty-move-rule":       pb.GameResultStatus_GAME_RESULT_STATUS_FIFTY_MOVE_RULE,
		"seventyfive-move-rule": pb.GameResultStatus_GAME_RESULT_STATUS_SEVENTYFIVE_MOVE_RULE,
		"stalemate":             pb.GameResultStatus_GAME_RESULT_STATUS_STALEMATE,
		"resignation":           pb.GameResultStatus_GAME_RESULT_STATUS_RESIGNATION,
		"draw-agreed":           pb.GameResultStatus_GAME_RESULT_STATUS_DRAW_AGREED,
		"flagged":               pb.GameResultStatus_GAME_RESULT_STATUS_FLAGGED,
		"adjudication":          pb.GameResultStatus_GAME_RESULT_STATUS_ADJUDICATION,
		"timed-out":             pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT,
		"aborted":               pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED,
		"interrupted":           pb.GameResultStatus_GAME_RESULT_STATUS_INTERRUPTED,
	}

	gameStateNameToProto := map[string]pb.GameState{
		"idle":          pb.GameState_GAME_STATE_IDLE,
		"waiting-start": pb.GameState_GAME_STATE_WAITING_START,
		"in-progress":   pb.GameState_GAME_STATE_IN_PROGRESS,
		"finished":      pb.GameState_GAME_STATE_FINISHED,
		"interrupted":   pb.GameState_GAME_STATE_INTERRUPTED,
	}

	for _, v := range gameVariants.Data {
		protoEnum, ok := gameVariantNameToProto[v.Name]
		if !ok {
			a.Log.Warn("proto mappings cache game_variant not exist", slog.String("name", v.Name))
			continue
		}

		a.protoMappingsCache.gameVariantsProtoToDB[protoEnum] = v.ID
		a.protoMappingsCache.gameVariantsDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameTimeKinds.Data {
		protoEnum, ok := gameTimeKindNameToProto[v.Name]
		if !ok {
			a.Log.Warn("proto mappings cache game_time_kind not exist", slog.String("name", v.Name))
			continue
		}

		a.protoMappingsCache.gameTimeKindsProtoToDB[protoEnum] = v.ID
		a.protoMappingsCache.gameTimeKindsDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameTimeCategories.Data {
		protoEnum, ok := gameTimeCategoryNameToProto[v.Name]
		if !ok {
			a.Log.Warn("proto mappings cache game_time_category not exist", slog.String("name", v.Name))
			continue
		}

		a.protoMappingsCache.gameTimeCategoriesProtoToDB[protoEnum] = v.ID
		a.protoMappingsCache.gameTimeCategoriesDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameResults.Data {
		protoEnum, ok := gameResultNameToProto[v.Name]
		if !ok {
			a.Log.Warn("proto mappings cache game_result not exist", slog.String("name", v.Name))
			continue
		}

		a.protoMappingsCache.gameResultsProtoToDB[protoEnum] = v.ID
		a.protoMappingsCache.gameResultsDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameResultStatuses.Data {
		protoEnum, ok := gameResultStatusNameToProto[v.Name]
		if !ok {
			a.Log.Warn("proto mappings cache game_result_status not exist", slog.String("name", v.Name))
			continue
		}

		a.protoMappingsCache.gameResultStatusesProtoToDB[protoEnum] = v.ID
		a.protoMappingsCache.gameResultStatusesDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameStates.Data {
		protoEnum, ok := gameStateNameToProto[v.Name]
		if !ok {
			a.Log.Warn("proto mappings cache game_state not exist", slog.String("name", v.Name))
			continue
		}

		a.protoMappingsCache.gameStatesProtoToDB[protoEnum] = v.ID
		a.protoMappingsCache.gameStatesDBToProto[v.ID] = protoEnum
	}

	return nil
}

func (a *ApiHandler) gameVariantProtoToID(x pb.GameVariant) int64 {
	return a.protoMappingsCache.gameVariantsProtoToDB[x]
}

func (a *ApiHandler) gameTimeKindProtoToID(x pb.GameTimeKind) int64 {
	return a.protoMappingsCache.gameTimeKindsProtoToDB[x]
}

func (a *ApiHandler) gameTimeCategoryProtoToID(x pb.GameTimeCategory) int64 {
	return a.protoMappingsCache.gameTimeCategoriesProtoToDB[x]
}

func (a *ApiHandler) gameResultProtoToID(x pb.GameResult) int64 {
	return a.protoMappingsCache.gameResultsProtoToDB[x]
}

func (a *ApiHandler) gameResultStatusProtoToID(x pb.GameResultStatus) int64 {
	return a.protoMappingsCache.gameResultStatusesProtoToDB[x]
}

func (a *ApiHandler) gameStateProtoToID(x pb.GameState) int64 {
	return a.protoMappingsCache.gameStatesProtoToDB[x]
}

func (a *ApiHandler) gameVariantIDToProto(id int64) pb.GameVariant {
	return a.protoMappingsCache.gameVariantsDBToProto[id]
}

func (a *ApiHandler) gameTimeKindIDToProto(id int64) pb.GameTimeKind {
	return a.protoMappingsCache.gameTimeKindsDBToProto[id]
}

func (a *ApiHandler) gameTimeCategoryIDToProto(id int64) pb.GameTimeCategory {
	return a.protoMappingsCache.gameTimeCategoriesDBToProto[id]
}

func (a *ApiHandler) gameResultIDToProto(id int64) pb.GameResult {
	return a.protoMappingsCache.gameResultsDBToProto[id]
}

func (a *ApiHandler) gameResultStatusIDToProto(id int64) pb.GameResultStatus {
	return a.protoMappingsCache.gameResultStatusesDBToProto[id]
}

func (a *ApiHandler) gameStateIDToProto(id int64) pb.GameState {
	return a.protoMappingsCache.gameStatesDBToProto[id]
}

func (a *ApiHandler) sendLobbyInfo(userID, connID string, guest bool) error {
	return nil
}

func (a *ApiHandler) sendLobbyChatInfo(userID, connID string, guest bool) error {
	return nil
}

func (a *ApiHandler) sendGameInfo(gameID int64, userID, connID string, guest bool) error {
	return nil
}

func (a *ApiHandler) sendGameTvInfo(gameID int64, userID, connID string, guest bool) error {
	return nil
}

func (a *ApiHandler) publishPresenceDiff(ctx context.Context, channelsDiff persistence.PresenceChannelsDiff, userID, connID, username string, guest bool) error {
	presenceDiff := &pb.PresenceDiff{}

	if len(channelsDiff.UserJoined) > 0 {
		presenceDiff.Joined = make([]*pb.Presence, len(channelsDiff.UserJoined))
		for i, joinedChannel := range channelsDiff.UserJoined {
			presenceDiff.Joined[i] = &pb.Presence{
				UserId:   userID,
				Username: username,
				Guest:    guest,
				Channel:  joinedChannel,
			}
		}
	}

	if len(channelsDiff.UserLeft) > 0 {
		presenceDiff.Left = make([]*pb.Presence, len(channelsDiff.UserLeft))
		for i, leftChannel := range channelsDiff.UserLeft {
			presenceDiff.Left[i] = &pb.Presence{
				UserId:   userID,
				Username: username,
				Guest:    guest,
				Channel:  leftChannel,
			}
		}
	}

	presenceDiffMsg := &pb.Message{Event: &pb.Message_PresenceDiff{PresenceDiff: presenceDiff}}

	presenceDiffMsgBytes, err := protojson.Marshal(presenceDiffMsg)
	if err != nil {
		a.Log.Error("Message_PresenceDiff protojson marshal", slog.String("user_id", userID), slog.String("conn_id", connID), slog.Any("error", err))
		return err
	}

	if err := a.bus.rdb.Publish(ctx, "presence.diff."+userID, presenceDiffMsgBytes).Err(); err != nil {
		a.Log.Error("publish Message_PresenceDiff", slog.String("user_id", userID), slog.String("conn_id", connID), slog.Any("error", err))
		return err
	}

	return nil
}

func (a *ApiHandler) sendUserPresenceDiffToChannel(ctx context.Context, channelsDiff persistence.PresenceChannelsDiff, channel, userID, connID, username string, guest bool) error {
	presenceDiff := &pb.PresenceDiff{}
	if len(channelsDiff.UserJoined) > 0 {
		presenceDiff.Joined = make([]*pb.Presence, len(channelsDiff.UserJoined))
		for i, joinedChannel := range channelsDiff.UserJoined {
			presenceDiff.Joined[i] = &pb.Presence{
				UserId:   userID,
				Username: username,
				Guest:    guest,
				Channel:  joinedChannel,
			}
		}
	}

	if len(channelsDiff.UserLeft) > 0 {
		presenceDiff.Left = make([]*pb.Presence, len(channelsDiff.UserLeft))
		for i, leftChannel := range channelsDiff.UserLeft {
			presenceDiff.Left[i] = &pb.Presence{
				UserId:   userID,
				Username: username,
				Guest:    guest,
				Channel:  leftChannel,
			}
		}
	}

	presenceDiffMsg := &pb.Message{Event: &pb.Message_PresenceDiff{PresenceDiff: presenceDiff}}

	presenceDiffMsgBytes, err := protojson.Marshal(presenceDiffMsg)
	if err != nil {
		return err
	}

	if err := a.Rdb.Publish(ctx, channel, presenceDiffMsgBytes).Err(); err != nil {
		return err
	}

	return nil
}

func (a *ApiHandler) sendChannelPresenceStateToConn(ctx context.Context, channel, connID string) error {
	users, err := a.persistor.Presence().ListUsersInChannel(ctx, channel)
	if err != nil {
		return err
	}

	presences := make([]*pb.Presence, len(users))

	for i, info := range users {
		presences[i] = &pb.Presence{
			UserId:   info.ID,
			Username: info.Username,
			Guest:    info.Guest,
			Channel:  channel,
		}
	}

	presenceStateMsg := &pb.Message{Event: &pb.Message_PresenceState{PresenceState: &pb.PresenceState{Presences: presences}}}

	presenceStateMsgBytes, err := protojson.Marshal(presenceStateMsg)
	if err != nil {
		return err
	}

	if err := a.Rdb.Publish(ctx, "conn."+connID, presenceStateMsgBytes).Err(); err != nil {
		return err
	}

	return nil
}

func (a *ApiHandler) publishToUser(ctx context.Context, userID string, msg *pb.Message, channel *string) error {
	bb, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	topic := "user." + userID
	if channel != nil && *channel != "" {
		topic = "user." + userID + "." + *channel
	}

	return a.bus.rdb.Publish(ctx, topic, bb).Err()
}

func (a *ApiHandler) publishToConn(ctx context.Context, connID string, msg *pb.Message) error {
	bb, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	topic := "conn." + connID

	return a.bus.rdb.Publish(ctx, topic, bb).Err()
}

func (a *ApiHandler) publishToChannel(ctx context.Context, channel string, msg *pb.Message) error {
	bb, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	if channel == "" {
		return errors.New("empty channel")
	}

	return a.bus.rdb.Publish(ctx, channel, bb).Err()
}

func (a *ApiHandler) StartMatchmaking(ctx context.Context) {
	a.Log.Info("gameserver matchmaking started")

	matchmakingInterval := time.Second * 15
	ticker := time.NewTicker(matchmakingInterval)

loop:
	for {
		select {
		case <-ticker.C:
			a.tryMatchPoolPlayers(ctx)
		case <-ctx.Done():
			break loop
		}
	}
}

func (a *ApiHandler) tryMatchPoolPlayers(ctx context.Context) {
	a.Log.Debug("matchmaking trying to match pool players")

	for _, quickGame := range quickGames {
		a.tryMatchPoolPlayersForPool(ctx, quickGame.ClockSecs*1000, quickGame.IncrementSecs*1000, true)
		a.tryMatchPoolPlayersForPool(ctx, quickGame.ClockSecs*1000, quickGame.IncrementSecs*1000, false)
	}
}

func (a *ApiHandler) tryMatchPoolPlayersForPool(ctx context.Context, clockMS, incrementMS int32, rated bool) {
	pool := dbtype.Pool{ClockMS: clockMS, IncrementMS: incrementMS, Rated: rated}

	// a.Log.Debug("try match pool", slog.String("pool", pool.Name()))

	matchedPairs := make([][2]string, 0)

	for {
		res, err := a.persistor.Pool().MatchPair(ctx, pool)
		if err != nil {
			break
		}

		if len(res) < 2 {
			break
		}

		matchedPairs = append(matchedPairs, [2]string{res[0], res[1]})
	}

	if len(matchedPairs) == 0 {
		return
	}

	var wg sync.WaitGroup

	for _, pair := range matchedPairs {
		wg.Go(func() {
			a.processMatchedPoolPair(ctx, pair, pool)
		})
	}

	wg.Wait()
}

func (a *ApiHandler) processMatchedPoolPair(ctx context.Context, pair [2]string, pool dbtype.Pool) {
	a.Log.Debug("processing matched pool pair", slog.String("pool", pool.Name()), slog.Any("pair", pair))

	userID1, err1 := uuid.Parse(pair[0])

	userID2, err2 := uuid.Parse(pair[1])
	if err1 != nil || err2 != nil {
		a.Log.Error("invalid matched pair user id", slog.Any("error", errors.Join(err1, err2)))
	}

	// userMeta1, err3 := a.persistor.Pool().GetPoolUserMeta(ctx, userID1)
	// userMeta2, err4 := a.persistor.Pool().GetPoolUserMeta(ctx, userID1)
	// if err3 != nil || err4 != nil {
	// 	a.Log.Error("failed to get user meta", slog.Any("error", errors.Join(err3, err4)))
	// }

	username1, username2 := "guest", "guest"

	if pool.Rated {
		uname1, err5 := a.GetUsername(ctx, userID1.String())

		uname2, err6 := a.GetUsername(ctx, userID2.String())
		if err5 != nil || err6 != nil {
			a.Log.Error("failed to get usernames", slog.Any("error", errors.Join(err5, err6)))
		}

		username1, username2 = uname1, uname2
	}

	color1, color2 := pb.Color_COLOR_WHITE, pb.Color_COLOR_BLACK
	if rand.IntN(2) == 1 {
		color1, color2 = pb.Color_COLOR_BLACK, pb.Color_COLOR_WHITE
	}

	players := [2]gameplay.Player{
		{ID: userID1, Name: username1, Color: color1, Guest: !pool.Rated},
		{ID: userID2, Name: username2, Color: color2, Guest: !pool.Rated},
	}

	gtc := &pb.GameTimeControl{ClockMs: pool.ClockMS, IncrementMs: pool.IncrementMS}

	thresholds := []gameplay.CategoryThreshold{}
	for _, x := range a.categoryThresholds {
		thresholds = append(thresholds, gameplay.CategoryThreshold{
			UpperLimit:   x.upperLimit,
			TimeCategory: x.timeCategory,
		})
	}

	gs, err := gameplay.NewGameState(-1, players, gtc, true, thresholds, gameplay.WithRated(pool.Rated))
	if err != nil {
		a.Log.Error("gameplay.NewGameState", slog.Any("error", err))
		return
	}

	gs.WaitingStart()

	gameSetter := models.GameSetter{
		GameVariantID:          omit.From(a.gameVariantProtoToID(gs.GameVariant)),
		GameTimeKindID:         omit.From(a.gameTimeKindProtoToID(gs.GameTimeKind)),
		GameTimeCategoryID:     omit.From(a.gameTimeCategoryProtoToID(gs.GameTimeCategory)),
		GameStateID:            omit.From(a.gameStateProtoToID(gs.GameState)),
		TimeControlClockMS:     omit.From(gs.GameTimeControl.ClockMs),
		TimeControlIncrementMS: omit.From(gs.GameTimeControl.IncrementMs),
		FirstMoveTimeoutMS:     omit.From(int32(gs.FirstMoveTimeout.Milliseconds())),
		ReconnectTimeoutMS:     omit.From(int32(gs.ReconnectTimeout.Milliseconds())),
		WhiteGameClock:         omit.From(gs.GameTimeControl.ClockMs),
		BlackGameClock:         omit.From(gs.GameTimeControl.ClockMs),
		Rated:                  omit.From(gs.Rated),
		StartTime:              omitnull.FromPtr(gs.StartTime),
		EndTime:                omitnull.FromPtr(gs.EndTime),
		LastMove:               omitnull.FromPtr(gs.LastMove),
		Fen:                    omit.From(gs.Chess.Position.Fen()),
		Repetitions:            omit.From(int32(gs.Chess.Repetitions)),
	}
	if pool.Rated {
		gameSetter.WhiteID = omitnull.From(gs.White.ID)
		gameSetter.BlackID = omitnull.From(gs.Black.ID)
	} else {
		gameSetter.GuestWhiteID = omitnull.From(gs.White.ID)
		gameSetter.GuestBlackID = omitnull.From(gs.Black.ID)
	}

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		gameSetter.GameResultID = omitnull.From(a.gameResultProtoToID(gs.GameResult))
	}

	if gs.GameResultStatus != pb.GameResultStatus_GAME_RESULT_STATUS_UNSPECIFIED {
		gameSetter.GameResultStatusID = omitnull.From(a.gameResultStatusProtoToID(gs.GameResultStatus))
	}

	moveSetters := make([]models.GameMoveSetter, len(gs.GameMoves))
	for i, move := range gs.GameMoves {
		moveSetter := models.GameMoveSetter{
			Fen:   omit.From(move.GetFen()),
			Uci:   omit.From(move.GetUci()),
			San:   omit.From(move.GetSan()),
			Check: omit.From(gs.Chess.Position.Check),
		}
		if move.GetPlayedAt() != nil {
			moveSetter.PlayedAt = omitnull.From(move.GetPlayedAt().AsTime())
		}

		moveSetters[i] = moveSetter
	}

	hashSetters := make([]models.GameHistoryHashSetter, len(gs.Chess.HistoryHashes))
	for i, hash := range gs.Chess.HistoryHashes {
		hashSetters[i] = models.GameHistoryHashSetter{
			Hash: omit.From(int64(hash)),
		}
	}

	game, err := a.persistor.Game().CreateGame(ctx, gameSetter, moveSetters, hashSetters)
	if err != nil {
		a.Log.Error("CreateGame", slog.Any("error", err))
		return
	}

	gs.GameID = game.ID

	a.gamestates[game.ID] = gs

	matchFoundMsg := &pb.Message{Event: &pb.Message_MatchFound{MatchFound: &pb.MatchFound{
		GameId:             gs.GameID,
		UserId:             "",
		GameState:          0,
		Color:              color1,
		Fen:                "",
		Ply:                0,
		Clocks:             &pb.Clocks{},
		LegalMoves:         []string{},
		TimeControl:        &pb.GameTimeControl{},
		OpponentInfo:       &pb.OpponentInfo{},
		ReconnectTimeoutMs: 0,
		FirstMoveTimeoutMs: 0,
		GameMoves:          []*pb.GameMove{},
		StartTime:          &timestamppb.Timestamp{},
	}}}

	matchFoundMsgBytes, err := protojson.Marshal(matchFoundMsg)
	if err != nil {
		a.Log.Error("protojson marshal Message_MatchFound", slog.Any("error", err))
	} else {
		if err := a.bus.rdb.Publish(ctx, "user."+userID1.String(), matchFoundMsgBytes).Err(); err != nil {
			a.Log.Error("publish Message_MatchFound", slog.Any("error", err))
		}
		if err := a.bus.rdb.Publish(ctx, "user."+userID2.String(), matchFoundMsgBytes).Err(); err != nil {
			a.Log.Error("publish Message_MatchFound", slog.Any("error", err))
		}
	}
}

func (a *ApiHandler) loadGameState(gameID int64) (*gameplay.GameState, error) {
	if gs, ok := a.gamestates[gameID]; ok {
		return gs, nil
	}

	filters := dbtype.GetGameByIDFilters{GetGameParams: api.GetGameParams{Embed: &[]api.GetGameParamsEmbed{api.GetGameParamsEmbedMoves}}, WithGameHashes: true}

	game, err := a.persistor.Game().GetGameByID(context.Background(), gameID, filters)
	if err != nil {
		return nil, err
	}

	return a.gameStateFromPersistence(context.Background(), game.Game, game.GameMoves.Val, game.GameHistoryHashes.Val)
}

func (a *ApiHandler) gameStateFromPersistence(ctx context.Context, game models.Game, moves *[]models.GameMove, hashes *[]models.GameHistoryHash) (*gameplay.GameState, error) {
	chess, err := engine.NewChess(game.Fen)
	if err != nil {
		return nil, err
	}

	chess.Repetitions = uint16(game.Repetitions)

	if hashes != nil && len(*hashes) > 0 {
		chess.HistoryHashes = make([]uint64, len(*hashes))
		for i, hash := range *hashes {
			chess.HistoryHashes[i] = uint64(hash.Hash)
		}
	}

	whiteID := game.GuestWhiteID.GetOr(game.WhiteID.MustGet())
	blackID := game.GuestBlackID.GetOr(game.BlackID.MustGet())

	whiteUsername, blackUsername := "guest", "guest"

	if game.Rated {
		uname1, err5 := a.GetUsername(ctx, whiteID.String())

		uname2, err6 := a.GetUsername(ctx, blackID.String())
		if err5 != nil || err6 != nil {
			a.Log.Error("failed to get usernames", slog.Any("error", errors.Join(err5, err6)))
		}

		whiteUsername, blackUsername = uname1, uname2
	}

	white := &gameplay.Player{
		ID:    whiteID,
		Name:  whiteUsername,
		Color: pb.Color_COLOR_WHITE,
		Guest: !game.Rated,
	}

	black := &gameplay.Player{
		ID:    blackID,
		Name:  blackUsername,
		Color: pb.Color_COLOR_BLACK,
		Guest: !game.Rated,
	}

	players := map[uuid.UUID]*gameplay.Player{whiteID: white, blackID: black}

	var gameMoves []*pb.GameMove

	if moves != nil && len(*moves) > 0 {
		gameMoves = make([]*pb.GameMove, len(*moves))

		for i, m := range *moves {
			move := &pb.GameMove{
				Fen:   m.Fen,
				Check: m.Check,
			}
			if m.Uci != "" {
				move.Uci = &m.Uci
			}
			if m.San != "" {
				move.San = &m.San
			}
			if m.PlayedAt.IsValue() {
				move.PlayedAt = timestamppb.New(m.PlayedAt.MustGet())
			}

			gameMoves[i] = move
		}
	}

	gs := &gameplay.GameState{
		Chess:            chess,
		GameID:           game.ID,
		White:            white,
		Black:            black,
		Players:          players,
		Guest:            !game.Rated,
		GameVariant:      a.gameVariantIDToProto(game.GameVariantID),
		GameTimeKind:     a.gameTimeKindIDToProto(game.GameTimeKindID),
		GameTimeCategory: a.gameTimeCategoryIDToProto(game.GameTimeCategoryID),
		GameTimeControl:  &pb.GameTimeControl{ClockMs: game.TimeControlClockMS, IncrementMs: game.TimeControlIncrementMS},
		GameState:        a.gameStateIDToProto(game.GameStateID),
		ReconnectTimeout: time.Duration(game.ReconnectTimeoutMS) * time.Millisecond,
		FirstMoveTimeout: time.Duration(game.FirstMoveTimeoutMS) * time.Millisecond,
		LastMove:         game.LastMove.Ptr(),
		StartTime:        game.StartTime.Ptr(),
		EndTime:          game.EndTime.Ptr(),
		Rated:            game.Rated,
		GameMoves:        gameMoves,
		// TimerAction:      make(chan gameplay.timerAction),
	}

	if game.GameResultID.IsValue() {
		gs.GameResult = a.gameResultIDToProto(game.GameResultID.MustGet())
	}

	if game.GameResultStatusID.IsValue() {
		gs.GameResultStatus = a.gameResultStatusIDToProto(game.GameResultStatusID.MustGet())
	}

	return gs, nil
}

func debug_print_game_info(gs *gameplay.GameState) {
	fmt.Printf("game_id: %d\n", gs.GameID)
	fmt.Printf("rated: %v\n", gs.Rated)
	fmt.Printf("white: %s\n", gs.White.Name)
	fmt.Printf("black: %s\n", gs.Black.Name)
	fmt.Printf("variant: %s\n", gs.GameVariant.String())
	fmt.Printf("time_category: %s\n", gs.GameTimeCategory.String())
	fmt.Printf("time_kind: %s\n", gs.GameTimeKind.String())
	fmt.Printf("time_control_clock_ms: %d\n", gs.GameTimeControl.GetClockMs())
	fmt.Printf("time_control_increment_ms: %d\n", gs.GameTimeControl.GetIncrementMs())
	fmt.Printf("state: %s\n", gs.GameState.String())
	fmt.Printf("result: %s\n", gs.GameResult.String())
	fmt.Printf("result_status: %s\n", gs.GameResultStatus.String())
	fmt.Printf("start_time: %s\n", gs.StartTime)
	fmt.Printf("last_move: %s\n", gs.LastMove)
	fmt.Printf("game_moves: %v\n", gs.GameMoves)
	fmt.Printf("repetitions: %v\n", gs.Chess.Repetitions)
	fmt.Printf("history_hashes: %v\n", gs.Chess.HistoryHashes)
}
