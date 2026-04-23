package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/gameplay"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
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
	case *pb.Message_Heartbeat:
		a.handleIPCHeartbeatMsg(msg.GetHeartbeat())

	case *pb.Message_LeaveTab:
		a.handleIPCLeaveTabMsg(msg.GetLeaveTab())

	case *pb.Message_LeaveSite:
		a.handleIPCLeaveSiteMsg(msg.GetLeaveSite())

	case *pb.Message_InitializeChannels:
		a.handleIPCInitializeChannelsMsg(msg.GetInitializeChannels())

	case *pb.Message_RequestInitialChannelsInfo:
		a.handleIPCRequestInitialChannelsInfoMsg(msg.GetRequestInitialChannelsInfo())
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

	_, _, err := a.persistor.Presence().RefreshPresence(context.Background(), uuid.MustParse(data.UserId), uuid.MustParse(data.ConnId), username, data.Guest)
	if err != nil {
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

	_, _, _, err := a.persistor.Presence().ClearPresence(context.Background(), uuid.MustParse(data.UserId), uuid.MustParse(data.ConnId), username, data.Guest)
	if err != nil {
		a.Log.Error("ClearPresence", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.Any("error", err))
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

	initializedChannelsMsg := &pb.Message{
		Event: &pb.Message_InitializedChannels{InitializedChannels: &pb.InitializedChannels{
			Channels: channels,
		}},
	}

	initializedChannelsMsgBytes, err := protojson.Marshal(initializedChannelsMsg)
	if err != nil {
		a.Log.Error("protojson marshal Message_InitializedChannels", slog.String("user_id", data.UserId), slog.Any("error", err))
		return
	}

	topic := "reply-initial-channels." + data.UserId + "." + data.ConnId
	if err := a.bus.rdb.Publish(context.Background(), topic, initializedChannelsMsgBytes).Err(); err != nil {
		a.Log.Error("hub publish Message_InitializedChannels", slog.String("user_id", data.UserId), slog.String("topic", "ipc"), slog.Any("error", err))
		return
	}
}

func (a *ApiHandler) handleIPCRequestInitialChannelsInfoMsg(data *pb.RequestInitialChannelsInfo) {
	var username string

	if data.Guest {
		username = "guest-" + data.UserId
	} else {
		uname, err := a.GetUsername(context.Background(), data.UserId)
		if err != nil {
			a.Log.Error("handleIPCRequestInitialChannelsInfoMsg get username", slog.Any("error", err))
			return
		}
		username = uname
	}

	for _, channel := range data.GetChannels() {
		oldChannels, newChannels, err := a.persistor.Presence().SetPresence(context.Background(), uuid.MustParse(data.UserId), uuid.MustParse(data.ConnId), username, data.Guest, channel)
		if err != nil {
			a.Log.Error("SetPresence", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.String("channel", channel), slog.Any("error", err))
			return
		}

		if err := a.broadcastPresenceDiff(context.Background(), oldChannels, newChannels, data.UserId, username); err != nil {
			a.Log.Error("broadcastPresenceDiff", slog.String("user_id", data.UserId), slog.String("conn_id", data.ConnId), slog.String("channel", channel), slog.Any("error", err))
			return
		}

		if channel == "lobby" {
			if err := a.sendLobbyInfo(data.UserId, data.ConnId); err != nil {
				a.Log.Error("sendLobbyInfo", slog.Any("error", err))
			}
		}
		if channel == "lobby.chat" {
			if err := a.sendLobbyChatInfo(data.UserId, data.ConnId); err != nil {
				a.Log.Error("sendLobbyChatInfo", slog.Any("error", err))
			}
		}
		if strings.HasPrefix(channel, "game.") || strings.HasPrefix(channel, "gametv.") {
			gameIDStr := strings.Split(channel, ".")
			if len(gameIDStr) != 3 {
				return
			}
			gameID, err := strconv.ParseInt(gameIDStr[2], 10, 64)
			if err != nil {
				a.Log.Error("gameid parseint", slog.Any("error", err))
				return
			}
			if err := a.sendGameInfo(gameID, data.UserId, data.ConnId); err != nil {
				a.Log.Error("sendGameInfo", slog.Any("error", err))
			}
		}
	}
}

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
	ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ(a)

	bb, _ := protojson.Marshal(&pb.Message{Event: &pb.Message_Echo{Echo: &pb.Echo{Message: strings.ToUpper(data.Message)}}})
	toUser, toConn, toLobby := "user."+authInfo.userID, "conn."+authInfo.connID, "lobby.chat"
	_ = []any{toUser, toConn, toLobby}
	a.bus.rdb.Publish(context.Background(), toUser, bb)
}

func (a *ApiHandler) handleWSCSeekGameMsg(authInfo clientAuthInfo, data *pb.SeekGame) {
	if data == nil {
		return
	}

	ctx := context.Background()

	if _, err := a.bus.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := fmt.Sprintf("seek_game.%d.%d-%d", authInfo.authState, data.GetTimeControl().GetClockMs(), data.GetTimeControl().GetIncrementMs())
		if err := p.ZAdd(ctx, key, redis.Z{Member: authInfo.userID, Score: float64(time.Now().UnixNano())}).Err(); err != nil {
			a.Log.Error("SeekGame add to queue", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
		}

		// 	if err := p.HSet(ctx, "clients_seeking", clientID, key).Err(); err != nil {
		// 		h.log.Error("seek_game add seek key for client", slog.String("user_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
		// 		return err
		// 	}

		if err := p.Publish(ctx, key, authInfo.userID).Err(); err != nil {
			a.Log.Error("SeekGame publish joined", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
			return err
		}

		return nil
	}); err != nil {
		a.Log.Error("SeekGame pipeline", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
	}

	// h.broadcastHubInfoToClient(uuid.MustParse(clientID))
}

func (a *ApiHandler) handleWSCCancelSeekGameMsg(authInfo clientAuthInfo, data *pb.CancelSeekGame) {
	fmt.Println(data, "handleWSCCancelSeekGameMsg")
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

func (a *ApiHandler) sendLobbyInfo(userID, connID string) error {
	return nil
}

func (a *ApiHandler) sendLobbyChatInfo(userID, connID string) error {
	return nil
}

func (a *ApiHandler) sendGameInfo(gameID int64, userID, connID string) error {
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

func (a *ApiHandler) broadcastPresenceDiff(ctx context.Context, oldChannels, newChannels []string, userID, username string) error {
	return nil
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

func (a *ApiHandler) getPresence(ctx context.Context, channel string) (*pb.Message, error) {
	return nil, nil
	// users, err := a.persistor.Presence().GetUsersInChannel(ctx, channel)
	// if err != nil {
	// 	return nil, err
	// }

	// presenceList := make([]*pb.UserPresence, 0)
	// for _, u := range users {
	// 	presenceList = append(presenceList, &pb.UserPresence{
	// 		UserId:   u.ID,
	// 		Username: u.Username,
	// 		Guest:    u.Guest,
	// 		Channel:  channel,
	// 	})
	// }
	// upsMsg := &pb.Message{Event: &pb.Message_UserPresences{UserPresences: &pb.UserPresences{Presences: presenceList}}}
	// return upsMsg, nil
}

func (a *ApiHandler) sendPresenceInfo(ctx context.Context, userID, connID uuid.UUID, username string, guest bool, channel string) error {
	presMsg, err := a.getPresence(ctx, channel)
	if err != nil {
		return err
	}
	_ = presMsg

	// if err := a.publishToConnID(ctx, connID, userID, presMsg); err != nil {
	// 	return err
	// }

	// send our presence to users in this channel also
	if err := a.broadcastPresence(ctx, userID.String(), username, guest, []string{channel}, false); err != nil {
		return err
	}

	return nil
}

func (a *ApiHandler) publishToUserID(ctx context.Context, userID string, msg *pb.Message, channel *string) error {
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

func ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ(a *ApiHandler) {
	players := [2]gameplay.Player{
		{ID: uuid.MustParse("75f751e4-737f-40de-beb8-3964cd4eeb29"), Name: "danko", Color: pb.Color_COLOR_WHITE, Guest: false},
		{ID: uuid.MustParse("14120b40-aa67-4cde-8a75-fcae8e057278"), Name: "bob", Color: pb.Color_COLOR_BLACK, Guest: false},
	}

	gtc := &pb.GameTimeControl{ClockMs: 300_000, IncrementMs: 0}

	thresholds := []gameplay.CategoryThreshold{}
	for _, x := range a.categoryThresholds {
		thresholds = append(thresholds, gameplay.CategoryThreshold{
			UpperLimit:   x.upperLimit,
			TimeCategory: x.timeCategory,
		})
	}

	gs, err := gameplay.NewGameState(1, players, gtc, true, thresholds)
	if err != nil {
		fmt.Println("gameplay.NewGameState: ", err.Error())
		return
	}

	gs.Start()

	gameSetter := models.GameSetter{
		// WhiteID:      omitnull.From(gs.White.ID),
		// BlackID:      omitnull.From(gs.Black.ID),
		WhiteIsGuest:           omit.From(gs.Guest),
		BlackIsGuest:           omit.From(gs.Guest),
		GuestWhiteID:           omitnull.From(gs.White.ID),
		GuestBlackID:           omitnull.From(gs.Black.ID),
		GameVariantID:          omit.From(a.protoGameVariantToID(gs.GameVariant)),
		GameTimeKindID:         omit.From(a.protoGameTimeKindToID(gs.GameTimeKind)),
		GameTimeCategoryID:     omit.From(a.protoGameTimeCategoryToID(gs.GameTimeCategory)),
		GameStateID:            omit.From(a.protoGameStateToID(gs.GameState)),
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
		PGN:                    omitnull.From("TODO"),
	}
	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		gameSetter.GameResultID = omitnull.From(a.protoGameResultToID(gs.GameResult))
	}
	if gs.GameResultStatus != pb.GameResultStatus_GAME_RESULT_STATUS_UNSPECIFIED {
		gameSetter.GameResultStatusID = omitnull.From(a.protoGameResultStatusToID(gs.GameResultStatus))
	}

	godump.DumpJSON(gameSetter, "GAME_SETTER")

	game, err := a.persistor.Game().CreateGame(context.TODO(), gameSetter, nil)
	if err != nil {
		fmt.Println("CreateGame err: ", err)
		return
	}
	godump.DumpJSON(game, "GAME_FINAL_RESULT")

	gs.GameID = game.ID

	fmt.Printf("game_id: %d\n", gs.GameID)
	fmt.Printf("rated: %v\n", gs.Rated)
	fmt.Printf("white: %s\n", gs.White.Name)
	fmt.Printf("black: %s\n", gs.Black.Name)
	fmt.Printf("variant: %s\n", gs.GameVariant.String())
	fmt.Printf("time_category: %s\n", gs.GameTimeCategory.String())
	fmt.Printf("time_kind: %s\n", gs.GameTimeKind.String())
	fmt.Printf("time_control_clock_ms: %s\n", gs.GameTimeControl.GetClockMs())
	fmt.Printf("time_control_increment_ms: %s\n", gs.GameTimeControl.GetIncrementMs())
	fmt.Printf("state: %s\n", gs.GameState.String())
	fmt.Printf("result: %s\n", gs.GameResult.String())
	fmt.Printf("result_status: %s\n", gs.GameResultStatus.String())
	fmt.Printf("start_time: %s\n", gs.StartTime)
	fmt.Printf("last_move: %s\n", gs.LastMove)
	fmt.Printf("history_moves: %v\n", gs.HistoryMoveInfos)
}
