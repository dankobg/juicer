package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/dankobg/juicer/persistence"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/dankobg/juicer/ws"
	"github.com/goforj/godump"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
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

	userLeft := make([]*pb.Presence, len(channelsDiff.UserLeft))

	for i, leftChannel := range channelsDiff.UserLeft {
		userLeft[i] = &pb.Presence{
			UserId:   data.UserId,
			Username: username,
			Guest:    data.Guest,
			Channel:  leftChannel,
		}
	}

	// pub to "presence.diff" for followers, and friends later

	presenceDiffMsg := &pb.Message{Event: &pb.Message_PresenceDiff{PresenceDiff: &pb.PresenceDiff{Left: userLeft}}}
	presenceDiffMsgBytes, err := protojson.Marshal(presenceDiffMsg)
	if err != nil {
		a.Log.Error("Message_PresenceDiff protojson marshal", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
		return
	}

	for _, leftChannel := range channelsDiff.UserLeft {
		if err := a.bus.rdb.Publish(context.Background(), leftChannel, presenceDiffMsgBytes).Err(); err != nil {
			a.Log.Error("publish Message_PresenceDiff", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.String("channel", leftChannel), slog.Any("error", err))
			return
		}
	}

	// clear presence
	// broadcast presence change
	// delete active seeks for connID
}

func (a *ApiHandler) handleIPCLeaveSiteMsg(data *pb.LeaveSite) {
	// delete active seeks for userID
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

		game, err := a.persistor.Game().GetGameByID(context.Background(), gameID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			a.Log.Error("handleIPCRequestInitialChannelsMsg GetGameByID", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.String("path", data.Path), slog.Any("error", err))
			return
		}

		if game.ID != 0 {
			switch game.GameStateID {
			case a.protoGameStateToID(pb.GameState_GAME_STATE_IN_PROGRESS):
				channels = append(channels, fmt.Sprintf("game.%d", game.ID), fmt.Sprintf("game.%d.chat", game.ID))
			case a.protoGameStateToID(pb.GameState_GAME_STATE_FINISHED),
				a.protoGameStateToID(pb.GameState_GAME_STATE_INTERRUPTED):
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

	for _, channel := range data.GetChannels() {
		channelsDiff, err := a.persistor.Presence().SetPresence(context.Background(), uuid.MustParse(data.UserId), uuid.MustParse(data.ConnId), username, data.Guest, channel)
		if err != nil && !errors.Is(err, redis.Nil) {
			a.Log.Error("SetPresence", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.String("channel", channel), slog.Any("error", err))
			return
		}

		_ = channelsDiff
		// a.broadcastPresenceDiff()

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

		if gameIDStr, found := strings.CutPrefix(channel, "game."); found {
			gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
			if err != nil {
				a.Log.Error("gameid parseint", slog.Any("error", err))
				return
			}

			if err := a.sendGameInfo(gameID, data.UserId, data.ConnId, data.Guest); err != nil {
				a.Log.Error("sendGameInfo", slog.Any("error", err))
			}
		}

		if gametvIDStr, found := strings.CutPrefix(channel, "gametv."); found {
			gameID, err := strconv.ParseInt(gametvIDStr, 10, 64)
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
	}
}

func (a *ApiHandler) handleWSCEchoMsg(authInfo clientAuthInfo, data *pb.Echo) {
	ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ(a)

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

	for _, v := range gameVariants.Data {
		switch v.Name {
		case "standard":
			a.protoMappingsCache.variants[pb.GameVariant_GAME_VARIANT_STANDARD] = v.ID
		case "atomic":
			a.protoMappingsCache.variants[pb.GameVariant_GAME_VARIANT_ATOMIC] = v.ID
		case "crazyhouse":
			a.protoMappingsCache.variants[pb.GameVariant_GAME_VARIANT_CRAZYHOUSE] = v.ID
		case "chess960":
			a.protoMappingsCache.variants[pb.GameVariant_GAME_VARIANT_CHESS960] = v.ID
		case "king-of-the-hill":
			a.protoMappingsCache.variants[pb.GameVariant_GAME_VARIANT_KING_OF_THE_HILL] = v.ID
		case "three-check":
			a.protoMappingsCache.variants[pb.GameVariant_GAME_VARIANT_THREE_CHECK] = v.ID
		case "horde":
			a.protoMappingsCache.variants[pb.GameVariant_GAME_VARIANT_HORDE] = v.ID
		case "racing-kings":
			a.protoMappingsCache.variants[pb.GameVariant_GAME_VARIANT_RACING_KINGS] = v.ID
		}
	}

	for _, v := range gameTimeKinds.Data {
		switch v.Name {
		case "realtime":
			a.protoMappingsCache.timeKinds[pb.GameTimeKind_GAME_TIME_KIND_REALTIME] = v.ID
		case "correspondance":
			a.protoMappingsCache.timeKinds[pb.GameTimeKind_GAME_TIME_KIND_CORRESPONDENCE] = v.ID
		case "unlimited":
			a.protoMappingsCache.timeKinds[pb.GameTimeKind_GAME_TIME_KIND_UNLIMITED] = v.ID
		}
	}

	for _, v := range gameTimeCategories.Data {
		switch v.Name {
		case "hyperbullet":
			a.protoMappingsCache.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_HYPERBULLET] = v.ID
		case "bullet":
			a.protoMappingsCache.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_BULLET] = v.ID
		case "blitz":
			a.protoMappingsCache.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_BLITZ] = v.ID
		case "rapid":
			a.protoMappingsCache.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_RAPID] = v.ID
		case "classical":
			a.protoMappingsCache.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_CLASSICAL] = v.ID
		}
	}

	for _, v := range gameResults.Data {
		switch v.Name {
		case "white-won":
			a.protoMappingsCache.results[pb.GameResult_GAME_RESULT_WHITE_WON] = v.ID
		case "black-won":
			a.protoMappingsCache.results[pb.GameResult_GAME_RESULT_BLACK_WON] = v.ID
		case "draw":
			a.protoMappingsCache.results[pb.GameResult_GAME_RESULT_DRAW] = v.ID
		case "interrupted":
			a.protoMappingsCache.results[pb.GameResult_GAME_RESULT_INTERRUPTED] = v.ID
		}
	}

	for _, v := range gameResultStatuses.Data {
		switch v.Name {
		case "checkmate":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_CHECKMATE] = v.ID
		case "insufficient-material":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_INSUFFICIENT_MATERIAL] = v.ID
		case "threefold-repetition":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_THREEFOLD_REPETITION] = v.ID
		case "fivefold-repetition":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_FIVEFOLD_REPETITION] = v.ID
		case "fifty-move-rule":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_FIFTY_MOVE_RULE] = v.ID
		case "seventyfive-move-rule":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_SEVENTYFIVE_MOVE_RULE] = v.ID
		case "stalemate":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_STALEMATE] = v.ID
		case "resignation":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_RESIGNATION] = v.ID
		case "draw-agreed":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_DRAW_AGREED] = v.ID
		case "flagged":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_FLAGGED] = v.ID
		case "adjudication":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_ADJUDICATION] = v.ID
		case "timed-out":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT] = v.ID
		case "aborted":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED] = v.ID
		case "interrupted":
			a.protoMappingsCache.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_INTERRUPTED] = v.ID
		}
	}

	for _, v := range gameStates.Data {
		switch v.Name {
		case "idle":
			a.protoMappingsCache.states[pb.GameState_GAME_STATE_IDLE] = v.ID
		case "waiting-start":
			a.protoMappingsCache.states[pb.GameState_GAME_STATE_WAITING_START] = v.ID
		case "in-progress":
			a.protoMappingsCache.states[pb.GameState_GAME_STATE_IN_PROGRESS] = v.ID
		case "finished":
			a.protoMappingsCache.states[pb.GameState_GAME_STATE_FINISHED] = v.ID
		case "interrupted":
			a.protoMappingsCache.states[pb.GameState_GAME_STATE_INTERRUPTED] = v.ID
		}
	}

	return nil
}

func (a *ApiHandler) protoGameVariantToID(x pb.GameVariant) int64 {
	return a.protoMappingsCache.variants[x]
}

func (a *ApiHandler) protoGameTimeKindToID(x pb.GameTimeKind) int64 {
	return a.protoMappingsCache.timeKinds[x]
}

func (a *ApiHandler) protoGameTimeCategoryToID(x pb.GameTimeCategory) int64 {
	return a.protoMappingsCache.timeCategories[x]
}

func (a *ApiHandler) protoGameResultToID(x pb.GameResult) int64 {
	return a.protoMappingsCache.results[x]
}

func (a *ApiHandler) protoGameResultStatusToID(x pb.GameResultStatus) int64 {
	return a.protoMappingsCache.resultStatuses[x]
}

func (a *ApiHandler) protoGameStateToID(x pb.GameState) int64 {
	return a.protoMappingsCache.states[x]
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

func channelsChanged(oldChannels, newChannels []string) bool {
	// they come sorted already
	if len(newChannels) != len(oldChannels) {
		return true
	}

	for i, nc := range newChannels {
		if nc != oldChannels[i] {
			return true
		}
	}

	return false
}

func (a *ApiHandler) broadcastPresenceDiff(ctx context.Context, channelsDiff persistence.PresenceChannelsDiff, userID, username string, guest bool) error {
	panic("TODO")

	// presenceDiff := &pb.PresenceDiff{}

	// userJoined := make([]*pb.Presence, len(channelsDiff.UserJoined))
	// for i, joinedChannel := range channelsDiff.UserJoined {
	// 	userJoined[i] = &pb.Presence{
	// 		UserId:   userID,
	// 		Username: username,
	// 		Guest:    guest,
	// 		Channel:  joinedChannel,
	// 	}
	// }

	// userLeft := make([]*pb.Presence, len(channelsDiff.UserLeft))
	// for i, joinedChannel := range channelsDiff.UserLeft {
	// 	userLeft[i] = &pb.Presence{
	// 		UserId:   userID,
	// 		Username: username,
	// 		Guest:    guest,
	// 		Channel:  joinedChannel,
	// 	}
	// }

	// presenceDiffMsg := &pb.Message{Event: &pb.Message_PresenceDiff{PresenceDiff: presenceDiff}}
	// presenceDiffMsgBytes, err := protojson.Marshal(presenceDiffMsg)
	// if err != nil {
	// 	a.Log.Error("Message_PresenceDiff protojson marshal", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
	// 	return
	// }
	// if err := a.bus.rdb.Publish(context.Background(), "presence.diff."+data.UserId, presenceDiffMsgBytes).Err(); err != nil {
	// 	a.Log.Error("publish Message_PresenceDiff", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
	// 	return
	// }

	// if err := a.sendPresenceInfo(context.Background(), uuid.MustParse(data.UserId), uuid.MustParse(data.ConnId), username, data.Guest, channel); err != nil {
	// 	a.Log.Error("sendPresenceInfo", slog.Any("error", err))
	// }

	// if !channelsChanged(oldChannels, newChannels) {
	// 	return nil
	// }

	// presenceDiffMsg := &pb.Message{Event: &pb.Message_PresenceDiff{PresenceDiff: &pb.PresenceEntry{
	// 	UserId:   userID,
	// 	Username: username,
	// 	Channels: newChannels,
	// }}}

	// presenceDiffMsgBytes, err := protojson.Marshal(presenceDiffMsg)
	// if err != nil {
	// 	return fmt.Errorf("protojson.Marshal Message_PresenceDiff: %w", err)
	// }

	// topic := "presence.diff." + userID
	// if err := a.Rdb.Publish(ctx, topic, presenceDiffMsgBytes).Err(); err != nil {
	// 	return fmt.Errorf("publish Message_PresenceDiff: %w", err)
	// }

	// return nil
}

func (a *ApiHandler) broadcastPresence(ctx context.Context, userID, username string, guest bool, channels []string, deleting bool) error {
	return nil
	// for _, channel := range channels {
	// 	toSend := &pb.Message{Event: &pb.Message_UserPresence{UserPresence: &pb.UserPresence{
	// 		UserId:   userID,
	// 		Username: username,
	// 		Guest:    guest,
	// 		Channel:  channel,
	// 		Deleting: deleting,
	// 	}}}
	// 	bb, err := protojson.Marshal(toSend)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Println("EEEEEEEEEEEEEEEEEEEEEEE", channel)
	// 	if err := a.bus.rdb.Publish(ctx, channel, bb).Err(); err != nil {
	// 		return err
	// 	}
	// }

	// return nil
}

func (a *ApiHandler) sendPresenceInfo(ctx context.Context, userID, connID uuid.UUID, username string, guest bool, channel string) error {
	userPresenceInfos, err := a.persistor.Presence().ListUsersInChannel(ctx, channel)
	if err != nil {
		return err
	}

	presences := make([]*pb.Presence, len(userPresenceInfos))

	for i, info := range userPresenceInfos {
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

	if err := a.Rdb.Publish(ctx, "conn."+connID.String(), presenceStateMsgBytes).Err(); err != nil {
		return err
	}

	// // send our presence to users in this channel also
	// if err := a.broadcastPresence(ctx, userID.String(), username, guest, []string{channel}, false); err != nil {
	// 	return err
	// }

	return nil
}

func (a *ApiHandler) publishToUserID(ctx context.Context, userID string, msg *pb.Message, channel *string) error {
	return nil
}

func (a *ApiHandler) publishToConnID(ctx context.Context, connID string, msg *pb.Message) error {
	return nil
}

func (a *ApiHandler) publishToChannel(ctx context.Context, channel string, msg *pb.Message) error {
	return nil
}

// func (b *Bus) pubToUser(userID string, evt *entity.EventWrapper, channel string) error {
// 	// Publish to a user, but pass in a specific channel. Only publish to those user sockets that are in this channel/realm/what-have-you.
// 	sanitized, err := sanitize(b.stores.UserStore, b.stores.GameStore, b.stores.TournamentStore, evt, userID)
// 	if err != nil {
// 		return err
// 	}
// 	bts, err := sanitized.Serialize()
// 	if err != nil {
// 		return err
// 	}
// 	var fullChannel string
// 	if channel == "" {
// 		fullChannel = "user." + userID
// 	} else {
// 		fullChannel = "user." + userID + "." + channel
// 	}

// 	return b.natsconn.Publish(fullChannel, bts)
// }

// func (b *Bus) pubToConnectionID(connID, userID string, evt *entity.EventWrapper) error {
// 	sanitized, err := sanitize(b.stores.UserStore, b.stores.GameStore, b.stores.TournamentStore, evt, userID)
// 	if err != nil {
// 		return err
// 	}
// 	bts, err := sanitized.Serialize()
// 	if err != nil {
// 		return err
// 	}
// 	return b.natsconn.Publish("connid."+connID, bts)
// }

func ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ(a *ApiHandler) {
	ids := []uuid.UUID{uuid.MustParse("cb38c8e1-2fb6-4b4c-bd10-e099953f8ee8"), uuid.MustParse("18a425bd-1325-4f38-81c4-b4fcc5ad9992")}
	_ = ids

	v1, e1 := a.persistor.Presence().ListUsersInChannel(context.TODO(), "lobby")
	v2, e2 := a.persistor.Presence().UsersCountInChannel(context.TODO(), "lobby")
	v3, e3 := a.persistor.Presence().ListChannelsForUser(context.TODO(), ids[0])
	v4, e4 := a.persistor.Presence().TotalActiveConnsCount(context.TODO())

	if e := errors.Join(e1, e2, e3, e4); e != nil {
		fmt.Println("--------------------------- KURAAAAAAC", e)
	}

	godump.DumpJSON(map[string]any{
		"lobby_users":       v1,
		"lobby_users_count": v2,
		"danko_channels":    v3,
		"active_conns":      v4,
	})

	// players := [2]gameplay.Player{
	// 	{ID: uuid.MustParse("75f751e4-737f-40de-beb8-3964cd4eeb29"), Name: "danko", Color: pb.Color_COLOR_WHITE, Guest: false},
	// 	{ID: uuid.MustParse("14120b40-aa67-4cde-8a75-fcae8e057278"), Name: "bob", Color: pb.Color_COLOR_BLACK, Guest: false},
	// }

	// gtc := &pb.GameTimeControl{ClockMs: 300_000, IncrementMs: 0}

	// thresholds := []gameplay.CategoryThreshold{}
	// for _, x := range a.categoryThresholds {
	// 	thresholds = append(thresholds, gameplay.CategoryThreshold{
	// 		UpperLimit:   x.upperLimit,
	// 		TimeCategory: x.timeCategory,
	// 	})
	// }

	// gs, err := gameplay.NewGameState(1, players, gtc, true, thresholds)
	// if err != nil {
	// 	fmt.Println("gameplay.NewGameState: ", err.Error())
	// 	return
	// }

	// gs.Start()

	// gameSetter := models.GameSetter{
	// 	// WhiteID:      omitnull.From(gs.White.ID),
	// 	// BlackID:      omitnull.From(gs.Black.ID),
	// 	WhiteIsGuest:           omit.From(gs.Guest),
	// 	BlackIsGuest:           omit.From(gs.Guest),
	// 	GuestWhiteID:           omitnull.From(gs.White.ID),
	// 	GuestBlackID:           omitnull.From(gs.Black.ID),
	// 	GameVariantID:          omit.From(a.protoGameVariantToID(gs.GameVariant)),
	// 	GameTimeKindID:         omit.From(a.protoGameTimeKindToID(gs.GameTimeKind)),
	// 	GameTimeCategoryID:     omit.From(a.protoGameTimeCategoryToID(gs.GameTimeCategory)),
	// 	GameStateID:            omit.From(a.protoGameStateToID(gs.GameState)),
	// 	TimeControlClockMS:     omit.From(gs.GameTimeControl.ClockMs),
	// 	TimeControlIncrementMS: omit.From(gs.GameTimeControl.IncrementMs),
	// 	FirstMoveTimeoutMS:     omit.From(int32(gs.FirstMoveTimeout.Milliseconds())),
	// 	ReconnectTimeoutMS:     omit.From(int32(gs.ReconnectTimeout.Milliseconds())),
	// 	WhiteGameClock:         omit.From(gs.GameTimeControl.ClockMs),
	// 	BlackGameClock:         omit.From(gs.GameTimeControl.ClockMs),
	// 	Rated:                  omit.From(gs.Rated),
	// 	StartTime:              omitnull.FromPtr(gs.StartTime),
	// 	EndTime:                omitnull.FromPtr(gs.EndTime),
	// 	LastMove:               omitnull.FromPtr(gs.LastMove),
	// 	Fen:                    omit.From(gs.Chess.Position.Fen()),
	// }
	// if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
	// 	gameSetter.GameResultID = omitnull.From(a.protoGameResultToID(gs.GameResult))
	// }
	// if gs.GameResultStatus != pb.GameResultStatus_GAME_RESULT_STATUS_UNSPECIFIED {
	// 	gameSetter.GameResultStatusID = omitnull.From(a.protoGameResultStatusToID(gs.GameResultStatus))
	// }

	// godump.DumpJSON(gameSetter, "GAME_SETTER")

	// game, err := a.persistor.Game().CreateGame(context.TODO(), gameSetter, nil)
	// if err != nil {
	// 	fmt.Println("CreateGame err: ", err)
	// 	return
	// }
	// godump.DumpJSON(game, "GAME_FINAL_RESULT")

	// gs.GameID = game.ID

	// fmt.Printf("game_id: %d\n", gs.GameID)
	// fmt.Printf("rated: %v\n", gs.Rated)
	// fmt.Printf("white: %s\n", gs.White.Name)
	// fmt.Printf("black: %s\n", gs.Black.Name)
	// fmt.Printf("variant: %s\n", gs.GameVariant.String())
	// fmt.Printf("time_category: %s\n", gs.GameTimeCategory.String())
	// fmt.Printf("time_kind: %s\n", gs.GameTimeKind.String())
	// fmt.Printf("time_control_clock_ms: %d\n", gs.GameTimeControl.GetClockMs())
	// fmt.Printf("time_control_increment_ms: %d\n", gs.GameTimeControl.GetIncrementMs())
	// fmt.Printf("state: %s\n", gs.GameState.String())
	// fmt.Printf("result: %s\n", gs.GameResult.String())
	// fmt.Printf("result_status: %s\n", gs.GameResultStatus.String())
	// fmt.Printf("start_time: %s\n", gs.StartTime)
	// fmt.Printf("last_move: %s\n", gs.LastMove)
	// fmt.Printf("history_moves: %v\n", gs.HistoryMoveInfos)
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

	res, err := a.persistor.Pool().ListPoolPlayers(ctx, pool)
	if err != nil {
		a.Log.Error("tryMatchPoolPlayersForPool listpoolplayers", slog.String("pool", pool.Name()), slog.Any("error", err))
		return
	}

	queueSize := len(res)
	if queueSize < 2 {
		// a.Log.Debug("pool queue size < 2")
		return
	}

	pairs := make(chan [2]string, queueSize/2)

	var wg sync.WaitGroup

	go func() {
		for i := 0; i+1 < queueSize; i += 2 {
			pairs <- [2]string{res[i], res[i+1]}
		}

		close(pairs)
	}()

	workers := queueSize / 2

	for range workers {
		wg.Go(func() {
			for pair := range pairs {
				a.processPoolUsersPair(ctx, pair, clockMS, incrementMS, rated)
			}
		})
	}

	wg.Wait()

	if queueSize%2 != 0 {
		a.Log.Debug("unmatched player remains in queue", slog.String("pool", pool.Name()), slog.String("user_id", res[queueSize-1]))
	}
}

func (a *ApiHandler) processPoolUsersPair(ctx context.Context, pair [2]string, clockMS, incrementMS int32, rated bool) {
	fmt.Println("----------------- PROCESS PAIR ------------------: ", pair, clockMS, incrementMS, rated)

	// userID1, err1 := uuid.Parse(pair[0])
	// userID2, err2 := uuid.Parse(pair[1])
	// if err1 != nil || err2 != nil {
	// 	a.Log.Error("failed to parse user UUIDs", slog.String("user_id1", pair[0]), slog.String("user_id2", pair[1]))
	// 	return
	// }

	// username1, err3 := a.GetUsername(ctx, pair[0])
	// username2, err4 := a.GetUsername(ctx, pair[1])
	// if err3 != nil || err4 != nil {
	// 	a.Log.Error("failed to get pair of usernames", slog.String("user_id1", pair[0]), slog.String("user_id2", pair[1]))
	// 	return
	// }
}
