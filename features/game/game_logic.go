package game

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/features/chat"
	"github.com/dankobg/juicer/gameplay"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/dankobg/juicer/ws"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (g *GameService) FetchCategoryThresholds(ctx context.Context) error {
	gameTimeCategories, err := g.pst.GameTimeCategory.ListGameTimeCategories(ctx, ListGameTimeCategoriesFilters{})
	if err != nil {
		return err
	}

	for _, v := range gameTimeCategories.Data {
		var limit time.Duration = 1<<63 - 1
		if v.UpperTimeLimitSecs.IsValue() {
			limit = time.Second * time.Duration(v.UpperTimeLimitSecs.MustGet())
		}

		switch v.Name {
		case "hyperbullet":
			g.categoryThresholds = append(g.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_HYPERBULLET, upperLimit: limit})
		case "bullet":
			g.categoryThresholds = append(g.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_BULLET, upperLimit: limit})
		case "blitz":
			g.categoryThresholds = append(g.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_BLITZ, upperLimit: limit})
		case "rapid":
			g.categoryThresholds = append(g.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_RAPID, upperLimit: limit})
		case "classical":
			g.categoryThresholds = append(g.categoryThresholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_CLASSICAL, upperLimit: limit})
		}
	}

	slices.SortFunc(g.categoryThresholds, func(a, b categoryThreshold) int {
		return cmp.Compare(a.upperLimit, b.upperLimit)
	})

	return nil
}

func (g *GameService) FetchProtoMappingsCacheLookups(ctx context.Context) error {
	gameVariants, e1 := g.pst.GameVariant.ListGameVariants(ctx, ListGameVariantsFilters{})
	gameTimeKinds, e2 := g.pst.GameTimeKind.ListGameTimeKinds(ctx, ListGameTimeKindsFilters{})
	gameTimeCategories, e3 := g.pst.GameTimeCategory.ListGameTimeCategories(ctx, ListGameTimeCategoriesFilters{})
	gameResults, e4 := g.pst.GameResult.ListGameResults(ctx, ListGameResultsFilters{})
	gameResultStatuses, e5 := g.pst.GameResultStatus.ListGameResultStatuses(ctx, ListGameResultStatusesFilters{})

	gameStates, e6 := g.pst.GameState.ListGameStates(ctx, ListGameStatesFilters{})
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
		"active":      pb.GameState_GAME_STATE_ACTIVE,
		"finished":    pb.GameState_GAME_STATE_FINISHED,
		"interrupted": pb.GameState_GAME_STATE_INTERRUPTED,
	}

	for _, v := range gameVariants.Data {
		protoEnum, ok := gameVariantNameToProto[v.Name]
		if !ok {
			g.log.Warn("proto mappings cache game_variant not exist", slog.String("name", v.Name))
			continue
		}

		g.protoMappingsCache.gameVariantsProtoToDB[protoEnum] = v.ID
		g.protoMappingsCache.gameVariantsDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameTimeKinds.Data {
		protoEnum, ok := gameTimeKindNameToProto[v.Name]
		if !ok {
			g.log.Warn("proto mappings cache game_time_kind not exist", slog.String("name", v.Name))
			continue
		}

		g.protoMappingsCache.gameTimeKindsProtoToDB[protoEnum] = v.ID
		g.protoMappingsCache.gameTimeKindsDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameTimeCategories.Data {
		protoEnum, ok := gameTimeCategoryNameToProto[v.Name]
		if !ok {
			g.log.Warn("proto mappings cache game_time_category not exist", slog.String("name", v.Name))
			continue
		}

		g.protoMappingsCache.gameTimeCategoriesProtoToDB[protoEnum] = v.ID
		g.protoMappingsCache.gameTimeCategoriesDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameResults.Data {
		protoEnum, ok := gameResultNameToProto[v.Name]
		if !ok {
			g.log.Warn("proto mappings cache game_result not exist", slog.String("name", v.Name))
			continue
		}

		g.protoMappingsCache.gameResultsProtoToDB[protoEnum] = v.ID
		g.protoMappingsCache.gameResultsDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameResultStatuses.Data {
		protoEnum, ok := gameResultStatusNameToProto[v.Name]
		if !ok {
			g.log.Warn("proto mappings cache game_result_status not exist", slog.String("name", v.Name))
			continue
		}

		g.protoMappingsCache.gameResultStatusesProtoToDB[protoEnum] = v.ID
		g.protoMappingsCache.gameResultStatusesDBToProto[v.ID] = protoEnum
	}

	for _, v := range gameStates.Data {
		protoEnum, ok := gameStateNameToProto[v.Name]
		if !ok {
			g.log.Warn("proto mappings cache game_state not exist", slog.String("name", v.Name))
			continue
		}

		g.protoMappingsCache.gameStatesProtoToDB[protoEnum] = v.ID
		g.protoMappingsCache.gameStatesDBToProto[v.ID] = protoEnum
	}

	return nil
}

func (g *GameService) gameVariantProtoToID(x pb.GameVariant) int64 {
	return g.protoMappingsCache.gameVariantsProtoToDB[x]
}

func (g *GameService) gameTimeKindProtoToID(x pb.GameTimeKind) int64 {
	return g.protoMappingsCache.gameTimeKindsProtoToDB[x]
}

func (g *GameService) gameTimeCategoryProtoToID(x pb.GameTimeCategory) int64 {
	return g.protoMappingsCache.gameTimeCategoriesProtoToDB[x]
}

func (g *GameService) gameResultProtoToID(x pb.GameResult) int64 {
	return g.protoMappingsCache.gameResultsProtoToDB[x]
}

func (g *GameService) gameResultStatusProtoToID(x pb.GameResultStatus) int64 {
	return g.protoMappingsCache.gameResultStatusesProtoToDB[x]
}

func (g *GameService) gameStateProtoToID(x pb.GameState) int64 {
	return g.protoMappingsCache.gameStatesProtoToDB[x]
}

func (g *GameService) gameVariantIDToProto(id int64) pb.GameVariant {
	return g.protoMappingsCache.gameVariantsDBToProto[id]
}

func (g *GameService) gameTimeKindIDToProto(id int64) pb.GameTimeKind {
	return g.protoMappingsCache.gameTimeKindsDBToProto[id]
}

func (g *GameService) gameTimeCategoryIDToProto(id int64) pb.GameTimeCategory {
	return g.protoMappingsCache.gameTimeCategoriesDBToProto[id]
}

func (g *GameService) gameResultIDToProto(id int64) pb.GameResult {
	return g.protoMappingsCache.gameResultsDBToProto[id]
}

func (g *GameService) gameResultStatusIDToProto(id int64) pb.GameResultStatus {
	return g.protoMappingsCache.gameResultStatusesDBToProto[id]
}

func (g *GameService) gameStateIDToProto(id int64) pb.GameState {
	return g.protoMappingsCache.gameStatesDBToProto[id]
}

func (g *GameService) onGameEvent(event gameplay.GameEvent) {
	switch ev := event.(type) {
	case gameplay.PlayMoveUCIEvent:
		g.handlePlayMoveUCIEvent(ev)

	case gameplay.PlayMoveUCIErrorEvent:
		g.handlePlayMoveUCIErrorEvent(ev)

	case gameplay.AbortEvent:
		g.handleAbortGameEvent(ev)

	case gameplay.AbortErrorEvent:
		g.handleAbortGameErrorEvent(ev)

	case gameplay.ResignEvent:
		g.handleResignGameEvent(ev)

	case gameplay.ResignErrorEvent:
		g.handleResignGameErrorEvent(ev)

	case gameplay.OfferDrawEvent:
		g.handleOfferDrawEvent(ev)

	case gameplay.OfferDrawErrorEvent:
		g.handleOfferDrawErrorEvent(ev)

	case gameplay.AcceptDrawEvent:
		g.handleAcceptDrawEvent(ev)

	case gameplay.AcceptDrawErrorEvent:
		g.handleAcceptDrawErrorEvent(ev)

	case gameplay.DeclineDrawEvent:
		g.handleDeclineDrawEvent(ev)

	case gameplay.DeclineDrawErrorEvent:
		g.handleDeclineDrawErrorEvent(ev)

	case gameplay.GameFinishedEvent:
		g.handleGameFinishedEvent(ev)

	case gameplay.PlayerDisconnected:
		g.handleGamePlayerDisconnectedEvent(ev)

	case gameplay.PlayerReconnected:
		g.handleGamePlayerReconnectedEvent(ev)
	}
}

func (g *GameService) handlePlayMoveUCIEvent(event gameplay.PlayMoveUCIEvent) {
	g.log.Debug("handlePlayMoveUCIEvent", slog.String("user_id", event.UserID.String()), slog.Int64("game_id", event.GameID))

	whiteSecs := int64(event.WhiteRemainingGameTime / time.Second)
	whiteNS := int64(event.WhiteRemainingGameTime % time.Second)
	blackSecs := int64(event.BlackRemainingGameTime / time.Second)
	blackNS := int64(event.BlackRemainingGameTime % time.Second)

	gameSetter := models.GameSetter{
		Fen:                    omit.From(event.Position.Fen()),
		LastMove:               omitnull.FromPtr(event.LastMove),
		Repetitions:            omit.From(int32(event.Repetitions)),
		WhiteGameRemainingSecs: omit.From(int32(whiteSecs)),
		WhiteGameRemainingNS:   omit.From(whiteNS),
		BlackGameRemainingSecs: omit.From(int32(blackSecs)),
		BlackGameRemainingNS:   omit.From(blackNS),
	}

	moveSetter := &models.GameMoveSetter{
		GameID:   omit.From(event.GameID),
		Fen:      omit.From(event.Position.Fen()),
		Uci:      omit.From(event.Uci),
		San:      omit.From(event.San),
		Lan:      omit.From(event.Lan),
		PlayedAt: omitnull.FromPtr(event.LastMove),
	}

	hashSetter := &models.GameHistoryHashSetter{
		GameID: omit.From(event.GameID),
		Hash:   omit.From(int64(event.Position.Hash)),
	}

	if _, err := g.pst.Game.UpdateGame(context.Background(), event.GameID, gameSetter, moveSetter, hashSetter); err != nil {
		g.log.Error("UpdateGame", slog.Int64("game_id", event.GameID), slog.Any("error", err))
	}

	moveAckMsg := &pb.Message{Event: &pb.Message_MoveAck{MoveAck: &pb.MoveAck{
		GameId:  int32(event.GameID),
		Version: int32(event.Version),
	}}}

	moveAckMsgBytes, err := protojson.Marshal(moveAckMsg)
	if err != nil {
		g.log.Error("protojson.Marshal Message_MoveAck", slog.Any("error", err))
	} else {
		topic := fmt.Sprintf("user.%s.game.%d", event.UserID.String(), event.GameID)
		if err := g.bus.Publish(context.Background(), topic, moveAckMsgBytes); err != nil {
			g.log.Error("publish Message_MoveAck", slog.Int64("game_id", event.GameID), slog.Any("error", err))
		}
	}

	moveSyncMsg := &pb.Message{Event: &pb.Message_MoveSync{MoveSync: &pb.MoveSync{
		GameId: int32(event.GameID),
		Uci:    event.Uci,
		San:    event.San,
		Lan:    event.Lan,
		Fen:    event.Position.Fen(),
		Ply:    uint32(event.Position.Ply),
		Clocks: &pb.Clocks{
			White: durationpb.New(event.WhiteRemainingGameTime),
			Black: durationpb.New(event.BlackRemainingGameTime),
		},
		LegalMoves: event.LegalMoves,
		Version:    int32(event.Version),
		PlayedAt:   timestamppb.New(event.PlayedAt),
	}}}

	moveSyncMsgBytes, err := protojson.Marshal(moveSyncMsg)
	if err != nil {
		g.log.Error("protojson.Marshal Message_MoveSync", slog.Any("error", err))
	} else {
		topic := fmt.Sprintf("game.%d", event.GameID)
		if err := g.bus.Publish(context.Background(), topic, moveSyncMsgBytes); err != nil {
			g.log.Error("publish Message_MoveSync", slog.Int64("game_id", event.GameID), slog.Any("error", err))
		}
	}
}

func (g *GameService) handlePlayMoveUCIErrorEvent(event gameplay.PlayMoveUCIErrorEvent) {
	g.log.Debug("handlePlayMoveUCIErrorEvent", slog.String("user_id", event.UserID.String()), slog.Int64("game_id", event.GameID), slog.Any("error", event.Err))
}

func (g *GameService) handleAbortGameEvent(event gameplay.AbortEvent) {
	g.log.Debug("handleAbortGameEvent", slog.Int64("game_id", event.GameID))

	gameSetter := models.GameSetter{
		EndTime: omitnull.From(event.EndTime),
	}

	if _, err := g.pst.Game.UpdateGame(context.Background(), event.GameID, gameSetter, nil, nil); err != nil {
		g.log.Error("UpdateGame", slog.Int64("game_id", event.GameID), slog.Any("error", err))
	}

	go func() {
		g.gameEvent <- gameplay.GameFinishedEvent{
			GameID:           event.GameID,
			GameResult:       event.GameResult,
			GameResultStatus: event.GameResultStatus,
			GameState:        event.GameState,
			EndTime:          time.Now(),
		}
	}()
}

func (g *GameService) handleAbortGameErrorEvent(event gameplay.AbortErrorEvent) {
	g.log.Debug("handleAbortGameErrorEvent", slog.String("user_id", event.UserID.String()), slog.Int64("game_id", event.GameID), slog.Any("error", event.Err))
}

func (g *GameService) handleResignGameEvent(event gameplay.ResignEvent) {
	g.log.Debug("handleResignGameEvent", slog.Int64("game_id", event.GameID))

	gameSetter := models.GameSetter{
		EndTime: omitnull.From(event.EndTime),
	}

	if _, err := g.pst.Game.UpdateGame(context.Background(), event.GameID, gameSetter, nil, nil); err != nil {
		g.log.Error("UpdateGame", slog.Int64("game_id", event.GameID), slog.Any("error", err))
	}

	go func() {
		g.gameEvent <- gameplay.GameFinishedEvent{
			GameID:           event.GameID,
			GameResult:       event.GameResult,
			GameResultStatus: event.GameResultStatus,
			GameState:        event.GameState,
			EndTime:          time.Now(),
		}
	}()
}

func (g *GameService) handleResignGameErrorEvent(event gameplay.ResignErrorEvent) {
	g.log.Debug("handleResignGameErrorEvent", slog.String("user_id", event.UserID.String()), slog.Int64("game_id", event.GameID), slog.Any("error", event.Err))
}

func (g *GameService) handleOfferDrawEvent(event gameplay.OfferDrawEvent) {
	g.log.Debug("handleOfferDrawEvent", slog.String("user_id", event.UserID.String()), slog.Int64("game_id", event.GameID))

	// @TODO: update active game pending draw offers

	drawOfferMsg := &pb.Message{Event: &pb.Message_DrawOffer{DrawOffer: &pb.DrawOffer{
		GameId:    int32(event.GameID),
		Ply:       uint32(event.Ply),
		OfferedBy: event.UserID.String(),
		OfferedAt: timestamppb.New(event.OfferedAt),
	}}}

	drawOfferMsgBytes, err := protojson.Marshal(drawOfferMsg)
	if err != nil {
		g.log.Error("protojson.Marshal Message_DrawOffer", slog.Any("error", err))
	} else {
		topic := fmt.Sprintf("user.%s.game.%d", event.OtherPlayer.String(), event.GameID)
		if err := g.bus.Publish(context.Background(), topic, drawOfferMsgBytes); err != nil {
			g.log.Error("publish Message_DrawOffer", slog.Int64("game_id", event.GameID), slog.Any("error", err))
		}
	}
}

func (g *GameService) handleOfferDrawErrorEvent(event gameplay.OfferDrawErrorEvent) {
	g.log.Debug("handleOfferDrawErrorEvent", slog.String("user_id", event.UserID.String()), slog.Int64("game_id", event.GameID), slog.Any("error", event.Err))
}

func (g *GameService) handleAcceptDrawEvent(event gameplay.AcceptDrawEvent) {
	g.log.Debug("handleAcceptDrawEvent", slog.Int64("game_id", event.GameID))

	gameSetter := models.GameSetter{
		EndTime: omitnull.From(event.EndTime),
	}

	if _, err := g.pst.Game.UpdateGame(context.Background(), event.GameID, gameSetter, nil, nil); err != nil {
		g.log.Error("UpdateGame", slog.Int64("game_id", event.GameID), slog.Any("error", err))
	}

	go func() {
		g.gameEvent <- gameplay.GameFinishedEvent{
			GameID:           event.GameID,
			GameResult:       event.GameResult,
			GameResultStatus: event.GameResultStatus,
			GameState:        event.GameState,
			EndTime:          time.Now(),
		}
	}()
}

func (g *GameService) handleAcceptDrawErrorEvent(event gameplay.AcceptDrawErrorEvent) {
	g.log.Debug("handleAcceptDrawErrorEvent", slog.String("user_id", event.UserID.String()), slog.Int64("game_id", event.GameID), slog.Any("error", event.Err))
}

func (g *GameService) handleDeclineDrawEvent(event gameplay.DeclineDrawEvent) {
	g.log.Debug("handleDeclineDrawEvent", slog.String("user_id", event.UserID.String()), slog.Int64("game_id", event.GameID))

	drawDeclinedMsg := &pb.Message{Event: &pb.Message_DrawDeclined{DrawDeclined: &pb.DrawDeclined{
		GameId:     int32(event.GameID),
		DeclinedBy: event.UserID.String(),
	}}}

	drawDeclinedMsgBytes, err := protojson.Marshal(drawDeclinedMsg)
	if err != nil {
		g.log.Error("protojson.Marshal Message_DrawDeclined", slog.Any("error", err))
	} else {
		topic := fmt.Sprintf("user.%s.game.%d", event.OtherPlayer.String(), event.GameID)
		if err := g.bus.Publish(context.Background(), topic, drawDeclinedMsgBytes); err != nil {
			g.log.Error("publish Message_DrawDeclined", slog.Int64("game_id", event.GameID), slog.Any("error", err))
		}
	}
}

func (g *GameService) handleDeclineDrawErrorEvent(event gameplay.DeclineDrawErrorEvent) {
	g.log.Debug("handleDeclineDrawErrorEvent", slog.String("user_id", event.UserID.String()), slog.Int64("game_id", event.GameID), slog.Any("error", event.Err))
}

func (g *GameService) handleGamePlayerDisconnectedEvent(event gameplay.PlayerDisconnected) {
	playerLeftMsg := &pb.Message{Event: &pb.Message_PlayerLeft{
		PlayerLeft: &pb.PlayerLeft{
			GameId: int32(event.GameID),
			UserId: event.UserID.String(),
			LeftAt: timestamppb.New(event.DisconnectedAt),
		},
	}}

	playerLeftMsgBytes, err := protojson.Marshal(playerLeftMsg)
	if err != nil {
		g.log.Error("Message_PlayerLeft protojson marshal", slog.String("user_id", event.UserID.String()), slog.Any("error", err))
		return
	} else {
		if err := g.bus.Publish(context.Background(), "user."+event.OtherUserID.String(), playerLeftMsgBytes); err != nil {
			g.log.Error("PlayerLeft publish", slog.String("user_id", event.UserID.String()), slog.Any("error", err))
			return
		}
	}
}

func (g *GameService) handleGamePlayerReconnectedEvent(event gameplay.PlayerReconnected) {
	playerRejoinedMsg := &pb.Message{Event: &pb.Message_PlayerRejoined{
		PlayerRejoined: &pb.PlayerRejoined{
			GameId:     int32(event.GameID),
			UserId:     event.UserID.String(),
			RejoinedAt: timestamppb.New(event.ReconnectedAt),
		},
	}}

	playerRejoinedMsgBytes, err := protojson.Marshal(playerRejoinedMsg)
	if err != nil {
		g.log.Error("Message_PlayerRejoined protojson marshal", slog.String("user_id", event.UserID.String()), slog.Any("error", err))
		return
	} else {
		if err := g.bus.Publish(context.Background(), "user."+event.OtherUserID.String(), playerRejoinedMsgBytes); err != nil {
			g.log.Error("PlayerRejoined publish", slog.String("user_id", event.UserID.String()), slog.Any("error", err))
			return
		}
	}
}

func (g *GameService) handleGameFinishedEvent(event gameplay.GameFinishedEvent) {
	g.log.Debug("handleGameFinishedEvent", slog.Int64("game_id", event.GameID))

	gameSetter := models.GameSetter{
		EndTime:            omitnull.From(event.EndTime),
		GameResultID:       omitnull.From(g.gameResultProtoToID(event.GameResult)),
		GameResultStatusID: omitnull.From(g.gameResultStatusProtoToID(event.GameResultStatus)),
		GameStateID:        omit.From(g.gameStateProtoToID(event.GameState)),
	}

	if _, err := g.pst.Game.UpdateGame(context.Background(), event.GameID, gameSetter, nil, nil); err != nil {
		g.log.Error("UpdateGame", slog.Int64("game_id", event.GameID), slog.Any("error", err))
	}

	gameFinishedMsg := &pb.Message{Event: &pb.Message_GameFinished{GameFinished: &pb.GameFinished{
		GameId:           int32(event.GameID),
		GameResult:       event.GameResult,
		GameResultStatus: event.GameResultStatus,
		GameState:        event.GameState,
	}}}

	gameFinishedMsgBytes, err := protojson.Marshal(gameFinishedMsg)
	if err != nil {
		g.log.Error("protojson.Marshal Message_GameFinished", slog.Int64("game_id", event.GameID), slog.Any("error", err))
	} else {
		topic := fmt.Sprintf("game.%d", event.GameID)
		if err := g.bus.Publish(context.Background(), topic, gameFinishedMsgBytes); err != nil {
			g.log.Error("publish Message_GameFinished", slog.Int64("game_id", event.GameID), slog.Any("error", err))
		}
	}

	if err := g.pst.ActiveGame.DeleteActiveGameByID(context.Background(), event.GameID); err != nil {
		g.log.Error("DeleteActiveGameByID", slog.Int64("game_id", event.GameID), slog.Any("error", err))
	}
}

func (g *GameService) onIPCMsg(m *redis.Message) {
	msg := &pb.Message{}
	if err := protojson.Unmarshal([]byte(m.Payload), msg); err != nil {
		g.log.Error("protojson.Unmarshal IPC Message")
		return
	}

	switch msg.GetEvent().(type) {
	case *pb.Message_ClientConnected:
		g.handleIPCClientConnectedMsg(msg.GetClientConnected())

	case *pb.Message_ClientDisconnected:
		g.handleIPCClientDisconnectedMsg(msg.GetClientDisconnected())

	case *pb.Message_Heartbeat:
		g.handleIPCHeartbeatMsg(msg.GetHeartbeat())

	case *pb.Message_LeaveTab:
		g.handleIPCLeaveTabMsg(msg.GetLeaveTab())

	case *pb.Message_LeaveSite:
		g.handleIPCLeaveSiteMsg(msg.GetLeaveSite())

	case *pb.Message_InitializeChannels:
		g.handleIPCInitializeChannelsMsg(msg.GetInitializeChannels())
	}
}

func (g *GameService) handleIPCHeartbeatMsg(data *pb.Heartbeat) {
	var username string

	userUid := uuid.MustParse(data.GetUserId())

	if data.GetGuest() {
		username = "guest-" + userUid.String()
	} else {
		uname, err := g.usrRdr.GetUsername(context.Background(), userUid)
		if err != nil {
			g.log.Error("handleIPCHeartbeatMsg get username", slog.Any("error", err))
			return
		}

		username = uname
	}

	err := g.presenceSvc.RefreshPresence(context.Background(), userUid, uuid.MustParse(data.GetConnId()), username, data.GetGuest())
	if err != nil && !errors.Is(err, redis.Nil) {
		g.log.Error("RefreshPresence", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.Any("error", err))
	}

	// refresh presence
	// broadcast presence change
}

func (g *GameService) handleIPCLeaveTabMsg(data *pb.LeaveTab) {
	leftAt := time.Now()

	userUid := uuid.MustParse(data.GetUserId())

	var username string

	if data.GetGuest() {
		username = "guest-" + userUid.String()
	} else {
		uname, err := g.usrRdr.GetUsername(context.Background(), userUid)
		if err != nil {
			g.log.Error("handleIPCLeaveTabMsg get username", slog.Any("error", err))
			return
		}

		username = uname
	}

	channelsDiff, err := g.presenceSvc.ClearPresence(context.Background(), userUid, uuid.MustParse(data.GetConnId()), username, data.GetGuest())
	if err != nil && !errors.Is(err, redis.Nil) {
		g.log.Error("ClearPresence", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.Any("error", err))
	}

	for _, ch := range channelsDiff.UserLeft {
		if strings.HasPrefix(ch, "game.") && !strings.HasSuffix(ch, ".chat") {
			parts := strings.Split(ch, "game.")

			gameID, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return
			}

			userID := uuid.MustParse(data.UserId)

			// @TODO: remove later, gs has player ids anyway
			isInActiveGame, err := g.pst.ActiveGame.IsUserInActiveGame(context.Background(), userID, gameID)
			if err != nil {
				g.log.Error("IsUserInActiveGame", slog.Any("error", err))
				return
			}

			if isInActiveGame {
				gs, err := g.loadGameState(gameID)
				if err != nil {
					g.log.Error("loadGameState", slog.Any("error", err))
					return
				}

				gs.GameCommand <- gameplay.LeftGame{
					GameID: gameID,
					UserID: userID,
					LefAt:  leftAt,
				}
			}
		}
	}

	if err := g.presenceSvc.PublishPresenceDiff(context.Background(), channelsDiff, userUid.String(), data.GetConnId(), username, data.GetGuest()); err != nil {
		g.log.Error("broadcastPresenceDiff", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.Any("error", err))
	}

	for _, leftChannel := range channelsDiff.UserLeft {
		if err := g.presenceSvc.SendUserPresenceDiffToChannel(context.Background(), channelsDiff, leftChannel, userUid.String(), data.GetConnId(), username, data.GetGuest()); err != nil {
			g.log.Error("sendUserPresenceDiffToChannel", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.Any("error", err))
		}
	}

	if err := g.pst.Pool.LeavePool(context.Background(), userUid); err != nil {
		g.log.Error("leave tab LeavePool", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.Any("error", err))
	}
}

func (g *GameService) handleIPCLeaveSiteMsg(data *pb.LeaveSite) {
}

func (g *GameService) handleIPCInitializeChannelsMsg(data *pb.InitializeChannels) {
	userUid := uuid.MustParse(data.GetUserId())
	channels := make([]string, 0)

	if data.GetPath() == "" || data.GetPath() == "/" {
		channels = append(channels, "lobby", "lobby.chat")
	} else {
		var gameID int64

		if gameIDStr, ok := strings.CutPrefix(data.GetPath(), "/game/"); ok {
			n, err := strconv.ParseInt(gameIDStr, 10, 64)
			if err != nil {
				g.log.Error("gameid parseint", slog.Any("error", err))
				return
			}

			gameID = n
		}

		if gametvIDStr, ok := strings.CutPrefix(data.GetPath(), "/gametv/"); ok {
			n, err := strconv.ParseInt(gametvIDStr, 10, 64)
			if err != nil {
				g.log.Error("gametvid parseint", slog.Any("error", err))
				return
			}

			gameID = n
		}

		game, err := g.pst.Game.GetGameByID(context.Background(), gameID, GetGameByIDFilters{})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			g.log.Error("handleIPCRequestInitialChannelsMsg GetGameByID", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.String("path", data.GetPath()), slog.Any("error", err))
			return
		}

		if game.ID != 0 {
			switch game.GameStateID {
			case g.gameStateProtoToID(pb.GameState_GAME_STATE_ACTIVE):
				channels = append(channels, fmt.Sprintf("game.%d", game.ID), fmt.Sprintf("game.%d.chat", game.ID))
			case g.gameStateProtoToID(pb.GameState_GAME_STATE_FINISHED),
				g.gameStateProtoToID(pb.GameState_GAME_STATE_INTERRUPTED):
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
		g.log.Error("protojson marshal Message_InitialChannels", slog.String("user_id", userUid.String()), slog.Any("error", err))
		return
	}

	topic := "reply-initial-channels." + userUid.String() + "." + data.GetConnId()
	if err := g.bus.Publish(context.Background(), topic, initialChannelsMsgBytes); err != nil {
		g.log.Error("hub publish Message_InitialChannels", slog.String("user_id", userUid.String()), slog.String("topic", "ipc"), slog.Any("error", err))
		return
	}
}

func (g *GameService) handleIPCClientConnectedMsg(data *pb.ClientConnected) {
	connectedAt := time.Now()

	userUid := uuid.MustParse(data.GetUserId())

	var username string

	if data.GetGuest() {
		username = "guest-" + userUid.String()
	} else {
		uname, err := g.usrRdr.GetUsername(context.Background(), userUid)
		if err != nil {
			g.log.Error("handleIPCClientConnectedMsg get username", slog.Any("error", err))
			return
		}

		username = uname
	}

	channelsDiff, err := g.presenceSvc.SetPresence(context.Background(), userUid, uuid.MustParse(data.GetConnId()), username, data.GetGuest(), data.GetChannels())
	if err != nil && !errors.Is(err, redis.Nil) {
		g.log.Error("SetPresence", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.Any("error", err))
		return
	}

	for _, ch := range channelsDiff.UserJoined {
		if strings.HasPrefix(ch, "game.") && !strings.HasSuffix(ch, ".chat") {
			parts := strings.Split(ch, "game.")

			gameID, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return
			}

			userID := uuid.MustParse(data.UserId)

			// @TODO: remove later, gs has player ids anyway
			isInActiveGame, err := g.pst.ActiveGame.IsUserInActiveGame(context.Background(), userID, gameID)
			if err != nil {
				g.log.Error("IsUserInActiveGame", slog.Any("error", err))
				return
			}

			if isInActiveGame {
				gs, err := g.loadGameState(gameID)
				if err != nil {
					g.log.Error("loadGameState", slog.Any("error", err))
					return
				}

				player := gs.GetPlayerByID(userID)
				if player == nil {
					return
				}

				if (player.Color == pb.Color_COLOR_WHITE && gs.WhiteDisconnectedAt != nil) || (player.Color == pb.Color_COLOR_BLACK && gs.BlackDisconnectedAt != nil) {
					gs.GameCommand <- gameplay.RejoinedGame{
						GameID:     gameID,
						UserID:     userID,
						RejoinedAt: connectedAt,
					}
				}
			}
		}
	}

	if err := g.presenceSvc.PublishPresenceDiff(context.Background(), channelsDiff, userUid.String(), data.GetConnId(), username, data.GetGuest()); err != nil {
		g.log.Error("broadcastPresenceDiff", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.Any("error", err))
	}

	for _, channel := range data.GetChannels() {
		if err := g.presenceSvc.SendChannelPresenceStateToConn(context.Background(), channel, data.GetConnId()); err != nil {
			g.log.Error("sendChannelPresenceStateToConn", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.Any("error", err))
		}

		if err := g.presenceSvc.SendUserPresenceDiffToChannel(context.Background(), channelsDiff, channel, userUid.String(), data.GetConnId(), username, data.GetGuest()); err != nil {
			g.log.Error("sendUserPresenceDiffToChannel", slog.String("user_id", userUid.String()), slog.String("conn_id", data.GetConnId()), slog.Any("error", err))
		}

		if channel == "lobby" {
			if err := g.sendLobbyInfo(userUid.String(), data.GetConnId(), data.GetGuest()); err != nil {
				g.log.Error("sendLobbyInfo", slog.Any("error", err))
			}
		}

		if channel == "lobby.chat" {
			if err := g.sendLobbyChatInfo(userUid.String(), data.GetConnId(), data.GetGuest()); err != nil {
				g.log.Error("sendLobbyChatInfo", slog.Any("error", err))
			}
		}

		if afterGameStr, isGame := strings.CutPrefix(channel, "game."); isGame {
			parts := strings.Split(afterGameStr, ".")

			var (
				gameIDStr  string
				isGameChat bool
			)

			if len(parts) > 0 {
				gameIDStr = parts[0]
			}

			if len(parts) == 2 && parts[1] == "chat" {
				isGameChat = true
			}

			gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
			if err != nil {
				return
			}

			if isGameChat {
				if err := g.sendGameChatInfo(gameID, userUid.String(), data.GetConnId(), data.GetGuest()); err != nil {
					g.log.Error("sendGameChatInfo", slog.Any("error", err))
				}
			} else {
				if err := g.sendGameInfo(gameID, userUid.String(), data.GetConnId(), data.GetGuest()); err != nil {
					g.log.Error("sendGameInfo", slog.Any("error", err))
				}
			}
		}

		if afterGameTvStr, isGameTv := strings.CutPrefix(channel, "gametv."); isGameTv {
			parts := strings.Split(afterGameTvStr, ".")

			var (
				gameIDStr    string
				isGameTvChat bool
			)

			if len(parts) > 0 {
				gameIDStr = parts[0]
			}

			if len(parts) == 2 && parts[1] == "chat" {
				isGameTvChat = true
			}

			gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
			if err != nil {
				return
			}

			if isGameTvChat {
				if err := g.sendGameTvChatInfo(gameID, userUid.String(), data.GetConnId(), data.GetGuest()); err != nil {
					g.log.Error("sendGameTvChatInfo", slog.Any("error", err))
				}
			} else {
				if err := g.sendGameTvInfo(gameID, userUid.String(), data.GetConnId(), data.GetGuest()); err != nil {
					g.log.Error("sendGameTvInfo", slog.Any("error", err))
				}
			}
		}
	}
}

func (g *GameService) handleIPCClientDisconnectedMsg(data *pb.ClientDisconnected) {}

func (g *GameService) onWSCMsg(m *redis.Message) {
	msg := &pb.Message{}
	if err := protojson.Unmarshal([]byte(m.Payload), msg); err != nil {
		g.log.Error("protojson.Unmarshal WSC Message")
		return
	}

	clientAuthInfo, err := extractWSCTopicParts(m.Channel)
	if err != nil {
		g.log.Error("extractWSCTopicParts", slog.String("channel", m.Channel), slog.String("pattern", m.Pattern), slog.String("payload", m.Payload), slog.Any("error", err))
		return
	}

	switch msg.GetEvent().(type) {
	case *pb.Message_Echo:
		g.handleWSCEchoMsg(clientAuthInfo, msg.GetEcho())
	case *pb.Message_SeekGame:
		g.handleWSCSeekGameMsg(clientAuthInfo, msg.GetSeekGame())
	case *pb.Message_CancelSeekGame:
		g.handleWSCCancelSeekGameMsg(clientAuthInfo, msg.GetCancelSeekGame())
	case *pb.Message_AbortGame:
		g.handleWSCAbortGame(clientAuthInfo, msg.GetAbortGame())
	case *pb.Message_ResignGame:
		g.handleWSCResignGame(clientAuthInfo, msg.GetResignGame())
	case *pb.Message_OfferDraw:
		g.handleWSCOfferDraw(clientAuthInfo, msg.GetOfferDraw())
	case *pb.Message_AcceptDraw:
		g.handleWSCAcceptDraw(clientAuthInfo, msg.GetAcceptDraw())
	case *pb.Message_DeclineDraw:
		g.handleWSCDeclineDraw(clientAuthInfo, msg.GetDeclineDraw())
	case *pb.Message_PlayMoveUci:
		g.handleWSCPlayMoveUCI(clientAuthInfo, msg.GetPlayMoveUci())
	case *pb.Message_SendLobbyChat:
		g.handleWSCSendLobbyChat(clientAuthInfo, msg.GetSendLobbyChat())
	case *pb.Message_ListLobbyChats:
		g.handleWSCListLobbyChats(clientAuthInfo, msg.GetListLobbyChats())
	case *pb.Message_SendGameChat:
		g.handleWSCSendGameChat(clientAuthInfo, msg.GetSendGameChat())
	case *pb.Message_ListGameChats:
		g.handleWSCListGameChats(clientAuthInfo, msg.GetListGameChats())
	}
}

func (g *GameService) handleWSCEchoMsg(authInfo clientAuthInfo, data *pb.Echo) {
	bb, _ := protojson.Marshal(&pb.Message{Event: &pb.Message_Echo{Echo: &pb.Echo{Message: strings.ToUpper(data.GetMessage())}}})
	toUser, toConn, toLobby := "user."+authInfo.userID, "conn."+authInfo.connID, "lobby.chat"
	_ = []any{toUser, toConn, toLobby}
	_ = g.bus.Publish(context.Background(), toUser, bb)
}

func (g *GameService) handleWSCSeekGameMsg(authInfo clientAuthInfo, data *pb.SeekGame) {
	if data == nil {
		return
	}

	pool := Pool{
		ClockMS:     data.GetGameTimeControl().GetClockMs(),
		IncrementMS: data.GetGameTimeControl().GetIncrementMs(),
		Rated:       authInfo.authState == ws.ClientAuth,
	}
	if err := g.pst.Pool.JoinPool(context.Background(), uuid.MustParse(authInfo.userID), pool); err != nil {
		g.log.Error("SeekGame join pool failed", slog.String("pool", pool.Name()), slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
		return
	}
}

func (g *GameService) handleWSCCancelSeekGameMsg(authInfo clientAuthInfo, data *pb.CancelSeekGame) {
	if err := g.pst.Pool.LeavePool(context.Background(), uuid.MustParse(authInfo.userID)); err != nil {
		g.log.Error("CancelSeekGame leave pool failed", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
		return
	}
}

func (g *GameService) handleWSCAbortGame(authInfo clientAuthInfo, data *pb.AbortGame) {
	userID := uuid.MustParse(authInfo.userID)

	// @TODO: remove later, gs has player ids anyway
	isInActiveGame, err := g.pst.ActiveGame.IsUserInActiveGame(context.Background(), userID, int64(data.GetGameId()))
	if err != nil {
		g.log.Error("IsUserInActiveGame", slog.Any("error", err))
		return
	}

	if !isInActiveGame {
		return
	}

	gs, err := g.loadGameState(int64(data.GetGameId()))
	if err != nil {
		g.log.Error("loadGameState", slog.Any("error", err))
		return
	}

	gs.GameCommand <- gameplay.AbortGameCmd{
		GameID: int64(data.GetGameId()),
		UserID: userID,
	}
}

func (g *GameService) handleWSCResignGame(authInfo clientAuthInfo, data *pb.ResignGame) {
	userID := uuid.MustParse(authInfo.userID)

	// @TODO: remove later, gs has player ids anyway
	isInActiveGame, err := g.pst.ActiveGame.IsUserInActiveGame(context.Background(), userID, int64(data.GetGameId()))
	if err != nil {
		g.log.Error("IsUserInActiveGame", slog.Any("error", err))
		return
	}

	if !isInActiveGame {
		return
	}

	gs, err := g.loadGameState(int64(data.GetGameId()))
	if err != nil {
		g.log.Error("loadGameState", slog.Any("error", err))
		return
	}

	gs.GameCommand <- gameplay.ResignGameCmd{
		GameID: int64(data.GetGameId()),
		UserID: userID,
	}
}

func (g *GameService) handleWSCOfferDraw(authInfo clientAuthInfo, data *pb.OfferDraw) {
	userID := uuid.MustParse(authInfo.userID)

	// @TODO: remove later, gs has player ids anyway
	isInActiveGame, err := g.pst.ActiveGame.IsUserInActiveGame(context.Background(), userID, int64(data.GetGameId()))
	if err != nil {
		g.log.Error("IsUserInActiveGame", slog.Any("error", err))
		return
	}

	if !isInActiveGame {
		return
	}

	gs, err := g.loadGameState(int64(data.GetGameId()))
	if err != nil {
		g.log.Error("loadGameState", slog.Any("error", err))
		return
	}

	gs.GameCommand <- gameplay.OfferDrawCmd{
		GameID: int64(data.GetGameId()),
		UserID: userID,
		Ply:    int(gs.Chess.Position.Ply),
	}
}

func (g *GameService) handleWSCAcceptDraw(authInfo clientAuthInfo, data *pb.AcceptDraw) {
	userID := uuid.MustParse(authInfo.userID)

	// @TODO: remove later, gs has player ids anyway
	isInActiveGame, err := g.pst.ActiveGame.IsUserInActiveGame(context.Background(), userID, int64(data.GetGameId()))
	if err != nil {
		g.log.Error("IsUserInActiveGame", slog.Any("error", err))
		return
	}

	if !isInActiveGame {
		return
	}

	gs, err := g.loadGameState(int64(data.GetGameId()))
	if err != nil {
		g.log.Error("loadGameState", slog.Any("error", err))
		return
	}

	gs.GameCommand <- gameplay.AcceptDrawCmd{
		GameID: int64(data.GetGameId()),
		UserID: userID,
		Ply:    int(gs.Chess.Position.Ply),
	}
}

func (g *GameService) handleWSCDeclineDraw(authInfo clientAuthInfo, data *pb.DeclineDraw) {
	userID := uuid.MustParse(authInfo.userID)

	// @TODO: remove later, gs has player ids anyway
	isInActiveGame, err := g.pst.ActiveGame.IsUserInActiveGame(context.Background(), userID, int64(data.GetGameId()))
	if err != nil {
		g.log.Error("IsUserInActiveGame", slog.Any("error", err))
		return
	}

	if !isInActiveGame {
		return
	}

	gs, err := g.loadGameState(int64(data.GetGameId()))
	if err != nil {
		g.log.Error("loadGameState", slog.Any("error", err))
		return
	}

	gs.GameCommand <- gameplay.DeclineDrawCmd{
		GameID: int64(data.GetGameId()),
		UserID: userID,
	}
}

func (g *GameService) handleWSCPlayMoveUCI(authInfo clientAuthInfo, data *pb.PlayMoveUCI) {
	userID := uuid.MustParse(authInfo.userID)

	// @TODO: remove later, gs has player ids anyway
	isInActiveGame, err := g.pst.ActiveGame.IsUserInActiveGame(context.Background(), userID, int64(data.GetGameId()))
	if err != nil {
		g.log.Error("IsUserInActiveGame", slog.Any("error", err))
		return
	}

	if !isInActiveGame {
		return
	}

	gs, err := g.loadGameState(int64(data.GetGameId()))
	if err != nil {
		g.log.Error("loadGameState", slog.Any("error", err))
		return
	}

	gs.GameCommand <- gameplay.PlayMoveUCICmd{
		GameID: int64(data.GetGameId()),
		UserID: userID,
		UCI:    data.GetUci(),
		Ack:    data.GetAck(),
	}
}

func (g *GameService) handleWSCSendLobbyChat(authInfo clientAuthInfo, data *pb.SendLobbyChat) {
	channel := "lobby.chat"
	lobbyMsg := data.GetMessage()

	username := "guest"

	if authInfo.authState == ws.ClientAuth {
		uname, err := g.usrRdr.GetUsername(context.Background(), uuid.MustParse(authInfo.userID))
		if err != nil {
			g.log.Error("lobby chat GetUsername", slog.String("user_id", authInfo.userID), slog.Any("error", err))
			return
		}

		username = uname
	}

	msg, err := g.chatSvc.AddChatMessage(context.Background(), channel, chat.ChatMessage{
		Channel:  channel,
		UserID:   uuid.MustParse(authInfo.userID),
		Username: username,
		Message:  lobbyMsg,
		PostedAt: time.Now(),
	})
	if err != nil {
		g.log.Error("lobby chat AddChatMessage", slog.String("user_id", authInfo.userID), slog.Any("error", err))
		return
	}

	lobbyChatMsg := &pb.Message{Event: &pb.Message_LobbyChat{LobbyChat: &pb.LobbyChat{
		MessageId: msg.MessageID,
		Message:   msg.Message,
		User: &pb.ChatUserSnapshot{
			Id:       msg.UserID.String(),
			Username: msg.Username,
		},
		PostedAt: timestamppb.New(msg.PostedAt),
	}}}

	lobbyChatMsgBytes, err := protojson.Marshal(lobbyChatMsg)
	if err != nil {
		g.log.Error("Message_LobbyChat protojson marshal", slog.String("user_id", authInfo.userID), slog.String("conn_id", authInfo.connID), slog.Any("error", err))
		return
	} else {
		if err := g.bus.Publish(context.Background(), "lobby.chat", lobbyChatMsgBytes); err != nil {
			g.log.Error("LobbyChat publish", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
			return
		}
	}
}

func (g *GameService) handleWSCSendGameChat(authInfo clientAuthInfo, data *pb.SendGameChat) {
	inGame, err := g.pst.ActiveGame.IsUserInActiveGame(context.Background(), uuid.MustParse(authInfo.userID), int64(data.GetGameId()))
	if err != nil {
		fmt.Println("-------------------------- NOT IN GAME ERRRRORRR ---------------------------- ", err)
		return
	}

	if !inGame {
		fmt.Println("-------------------------- NOT IN GAME ---------------------------- ")
		return
	}

	channel := fmt.Sprintf("game.%d.chat", data.GetGameId())
	gameMsg := data.GetMessage()

	username := "guest"

	if authInfo.authState == ws.ClientAuth {
		uname, err := g.usrRdr.GetUsername(context.Background(), uuid.MustParse(authInfo.userID))
		if err != nil {
			g.log.Error("game chat GetUsername", slog.String("user_id", authInfo.userID), slog.Int("game_id", int(data.GetGameId())), slog.Any("error", err))
			return
		}

		username = uname
	}

	msg, err := g.chatSvc.AddChatMessage(context.Background(), channel, chat.ChatMessage{
		Channel:  channel,
		UserID:   uuid.MustParse(authInfo.userID),
		Username: username,
		Message:  gameMsg,
		PostedAt: time.Now(),
	})
	if err != nil {
		g.log.Error("game chat AddChatMessage", slog.String("user_id", authInfo.userID), slog.Any("error", err))
		return
	}

	gameChatMsg := &pb.Message{Event: &pb.Message_GameChat{GameChat: &pb.GameChat{
		GameId:    data.GetGameId(),
		MessageId: msg.MessageID,
		Message:   msg.Message,
		User: &pb.ChatUserSnapshot{
			Id:       msg.UserID.String(),
			Username: msg.Username,
		},
		PostedAt: timestamppb.New(msg.PostedAt),
	}}}

	gameChatMsgBytes, err := protojson.Marshal(gameChatMsg)
	if err != nil {
		g.log.Error("Message_GameChat protojson marshal", slog.String("user_id", authInfo.userID), slog.String("conn_id", authInfo.connID), slog.Any("error", err))
		return
	} else {
		topic := fmt.Sprintf("game.%d.chat", data.GetGameId())
		if err := g.bus.Publish(context.Background(), topic, gameChatMsgBytes); err != nil {
			g.log.Error("GameChat publish", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
			return
		}
	}
}

func (g *GameService) handleWSCListGameChats(authInfo clientAuthInfo, data *pb.ListGameChats) {
	filters := chat.ChatFilters{Cursor: data.Cursor}
	if data.PageSize != nil {
		filters.PageSize = new(int(*data.PageSize))
	}

	gameChatChannel := fmt.Sprintf("game.%d.chat", data.GetGameId())

	msgs, err := g.chatSvc.ListChatMessages(context.Background(), gameChatChannel, filters)
	if err != nil {
		g.log.Error("list game chats ListChatMessages", slog.Any("error", err))
		return
	}

	lobbyChats := make([]*pb.LobbyChat, len(msgs.Data))
	for i, m := range msgs.Data {
		lobbyChats[i] = &pb.LobbyChat{
			MessageId: m.MessageID,
			Message:   m.Message,
			User: &pb.ChatUserSnapshot{
				Id:       m.UserID.String(),
				Username: m.Username,
			},
			PostedAt: timestamppb.New(m.PostedAt),
		}
	}

	lobbyChatList := &pb.Message{Event: &pb.Message_LobbyChats{LobbyChats: &pb.LobbyChatList{
		LobbyChats: lobbyChats,
		HasMore:    msgs.HasMore,
	}}}

	lobbyChatListBytes, err := protojson.Marshal(lobbyChatList)
	if err != nil {
		return
	}

	topic := fmt.Sprintf("user.%s.%s", authInfo.userID, gameChatChannel)
	if err := g.bus.Publish(context.Background(), topic, lobbyChatListBytes); err != nil {
		return
	}
}

func (g *GameService) handleWSCListLobbyChats(authInfo clientAuthInfo, data *pb.ListLobbyChats) {
	filters := chat.ChatFilters{Cursor: data.Cursor}
	if data.PageSize != nil {
		filters.PageSize = new(int(*data.PageSize))
	}

	msgs, err := g.chatSvc.ListChatMessages(context.Background(), "lobby.chat", filters)
	if err != nil {
		g.log.Error("list lobby chats ListChatMessages", slog.Any("error", err))
		return
	}

	lobbyChats := make([]*pb.LobbyChat, len(msgs.Data))
	for i, m := range msgs.Data {
		lobbyChats[i] = &pb.LobbyChat{
			MessageId: m.MessageID,
			Message:   m.Message,
			User: &pb.ChatUserSnapshot{
				Id:       m.UserID.String(),
				Username: m.Username,
			},
			PostedAt: timestamppb.New(m.PostedAt),
		}
	}

	lobbyChatList := &pb.Message{Event: &pb.Message_LobbyChats{LobbyChats: &pb.LobbyChatList{
		LobbyChats: lobbyChats,
		HasMore:    msgs.HasMore,
	}}}

	lobbyChatListBytes, err := protojson.Marshal(lobbyChatList)
	if err != nil {
		return
	}

	topic := fmt.Sprintf("user.%s.lobby.chat", authInfo.userID)
	if err := g.bus.Publish(context.Background(), topic, lobbyChatListBytes); err != nil {
		return
	}
}

type clientAuthInfo struct {
	userID    string
	connID    string
	authState ws.ClientAuthState
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

func (g *GameService) sendLobbyInfo(userID, connID string, guest bool) error {
	return nil
}

func (g *GameService) sendLobbyChatInfo(userID, connID string, guest bool) error {
	msgs, err := g.chatSvc.ListChatMessages(context.Background(), "lobby.chat", chat.ChatFilters{})
	if err != nil {
		g.log.Error("lobby chat ListChatMessages", slog.Any("error", err))
		return err
	}

	lobbyChats := make([]*pb.LobbyChat, len(msgs.Data))
	for i, m := range msgs.Data {
		lobbyChats[i] = &pb.LobbyChat{
			MessageId: m.MessageID,
			Message:   m.Message,
			User: &pb.ChatUserSnapshot{
				Id:       m.UserID.String(),
				Username: m.Username,
			},
			PostedAt: timestamppb.New(m.PostedAt),
		}
	}

	lobbyChatList := &pb.Message{Event: &pb.Message_LobbyChats{LobbyChats: &pb.LobbyChatList{
		LobbyChats: lobbyChats,
		HasMore:    msgs.HasMore,
	}}}

	lobbyChatListBytes, err := protojson.Marshal(lobbyChatList)
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("user.%s.lobby.chat", userID)
	if err := g.bus.Publish(context.Background(), topic, lobbyChatListBytes); err != nil {
		return err
	}

	return nil
}

func (g *GameService) sendGameChatInfo(gameID int64, userID, connID string, guest bool) error {
	gameChatChannel := fmt.Sprintf("game.%d.chat", gameID)

	msgs, err := g.chatSvc.ListChatMessages(context.Background(), gameChatChannel, chat.ChatFilters{})
	if err != nil {
		g.log.Error("game chat ListChatMessages", slog.Any("error", err))
		return err
	}

	gameChats := make([]*pb.GameChat, len(msgs.Data))
	for i, m := range msgs.Data {
		gameChats[i] = &pb.GameChat{
			GameId:    int32(gameID),
			MessageId: m.MessageID,
			Message:   m.Message,
			User: &pb.ChatUserSnapshot{
				Id:       m.UserID.String(),
				Username: m.Username,
			},
			PostedAt: timestamppb.New(m.PostedAt),
		}
	}

	gameChatList := &pb.Message{Event: &pb.Message_GameChats{GameChats: &pb.GameChatList{
		GameId:    int32(gameID),
		GameChats: gameChats,
		HasMore:   msgs.HasMore,
	}}}

	gameChatListBytes, err := protojson.Marshal(gameChatList)
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("user.%s.%s", userID, gameChatChannel)
	if err := g.bus.Publish(context.Background(), topic, gameChatListBytes); err != nil {
		return err
	}

	return nil
}

func (g *GameService) sendGameInfo(gameID int64, userID, connID string, guest bool) error {
	gs, err := g.loadGameState(int64(gameID))
	if err != nil {
		g.log.Error("loadGameState", slog.Any("error", err))
		return err
	}

	uid := uuid.MustParse(userID)
	player := gs.GetPlayerByID(uid)

	whitePlayerInfo := &pb.PlayerInfo{}
	blackPlayerInfo := &pb.PlayerInfo{}

	for _, p := range gs.Players {
		if p.Color == pb.Color_COLOR_WHITE {
			whitePlayerInfo.UserId = p.ID.String()
			whitePlayerInfo.Username = p.Username
			whitePlayerInfo.Guest = p.Guest
			whitePlayerInfo.Rating = 1500
		} else {
			blackPlayerInfo.UserId = p.ID.String()
			blackPlayerInfo.Username = p.Username
			blackPlayerInfo.Guest = p.Guest
			blackPlayerInfo.Rating = 1500
		}
	}

	legalMoves := make([]string, len(gs.Chess.LegalMoves))
	for i, legalMove := range gs.Chess.LegalMoves {
		legalMoves[i] = fmt.Sprint(legalMove.String())
	}

	pendingDrawOffers := make(map[string]*pb.DrawOffer)

	for k, v := range gs.PendingDrawOffers {
		pendingDrawOffers[k.String()] = &pb.DrawOffer{
			GameId:    int32(gs.GameID),
			Ply:       uint32(v.Ply),
			OfferedBy: v.OfferedBy.String(),
			OfferedAt: timestamppb.New(v.OfferedAt),
		}
	}

	gameInfo := &pb.GameInfo{
		GameId:           int32(gameID),
		GameVariant:      gs.GameVariant,
		GameTimeKind:     gs.GameTimeKind,
		GameTimeCategory: gs.GameTimeCategory,
		GameState:        gs.GameState,
		GameTimeControl:  gs.GameTimeControl,
		Color:            player.Color,
		Fen:              gs.Chess.Position.Fen(),
		Ply:              uint32(gs.Chess.Position.Ply),
		Clocks: &pb.Clocks{
			White: durationpb.New(gs.WhiteRemainingGameTime),
			Black: durationpb.New(gs.BlackRemainingGameTime),
		},
		Rated:              gs.Rated,
		LegalMoves:         legalMoves,
		White:              whitePlayerInfo,
		Black:              blackPlayerInfo,
		ReconnectTimeoutMs: int32(gs.ReconnectTimeout.Milliseconds()),
		FirstMoveTimeoutMs: int32(gs.FirstMoveTimeout.Milliseconds()),
		GameMoves:          gs.GameMoves,
		StartTime:          timestamppb.New(*gs.StartTime),
		GameResult:         gs.GameResult,
		GameResultStatus:   gs.GameResultStatus,
		Version:            int32(gs.Version),
		PendingDrawOffers:  pendingDrawOffers,
	}
	if gs.LastMove != nil {
		gameInfo.LastMove = timestamppb.New(*gs.LastMove)
	}

	if gs.EndTime != nil {
		gameInfo.EndTime = timestamppb.New(*gs.EndTime)
	}

	if gs.WhiteDisconnectedAt != nil {
		gameInfo.WhiteDisconnectedAt = timestamppb.New(*gs.WhiteDisconnectedAt)
	}

	if gs.BlackDisconnectedAt != nil {
		gameInfo.BlackDisconnectedAt = timestamppb.New(*gs.BlackDisconnectedAt)
	}

	gameInfoMsg := &pb.Message{Event: &pb.Message_GameInfo{GameInfo: gameInfo}}

	gameInfoMsgBytes, err := protojson.Marshal(gameInfoMsg)
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("user.%s.game.%d", userID, gameID)
	if err := g.bus.Publish(context.Background(), topic, gameInfoMsgBytes); err != nil {
		return err
	}

	return nil
}

func (g *GameService) sendGameTvChatInfo(gameID int64, userID, connID string, guest bool) error {
	return nil
}

func (g *GameService) sendGameTvInfo(gameID int64, userID, connID string, guest bool) error {
	return nil
}

func (g *GameService) publishToUser(ctx context.Context, userID string, msg *pb.Message, channel *string) error {
	bb, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	topic := "user." + userID
	if channel != nil && *channel != "" {
		topic = "user." + userID + "." + *channel
	}

	return g.bus.Publish(ctx, topic, bb)
}

func (g *GameService) publishToConn(ctx context.Context, connID string, msg *pb.Message) error {
	bb, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	topic := "conn." + connID

	return g.bus.Publish(ctx, topic, bb)
}

func (g *GameService) publishToChannel(ctx context.Context, channel string, msg *pb.Message) error {
	bb, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	if channel == "" {
		return errors.New("empty channel")
	}

	return g.bus.Publish(ctx, channel, bb)
}

func (g *GameService) PubsubProcess(ctx context.Context) {
	g.log.Info("game pubsub started")

	for {
		select {
		case event := <-g.gameEvent:
			g.onGameEvent(event)

		case msg := <-g.bus.SubMessages["ipc"]:
			g.onIPCMsg(msg)

		case msg := <-g.bus.SubMessages["wsc.*"]:
			g.onWSCMsg(msg)

		case <-ctx.Done():
			g.log.Debug("game pubsub ctx done")
			return
		}
	}
}

func (g *GameService) loadGameState(gameID int64) (*gameplay.GameState, error) {
	if gs, ok := g.gamestates[gameID]; ok {
		g.log.Debug("loadGameState success", slog.String("from", "memory"))
		return gs, nil
	}

	filters := GetGameByIDFilters{GetGameParams: api.GetGameParams{Embed: &[]api.GetGameParamsEmbed{api.GetGameParamsEmbedMoves}}, WithGameHashes: true}

	game, err := g.pst.Game.GetGameByID(context.Background(), gameID, filters)
	if err != nil {
		g.log.Error("loadGameState GetGameByID", slog.Int64("game_id", gameID), slog.Any("error", err))
		return nil, err
	}

	gs, err := g.gameStateFromPersistence(context.Background(), game.Game, game.GameMoves.Val, game.GameHistoryHashes.Val)
	if err != nil {
		g.log.Error("loadGameState gameStateFromPersistence", slog.Int64("game_id", gameID), slog.Any("error", err))
		return nil, err
	}

	g.log.Debug("loadGameState success", slog.String("from", "persistence"))

	gs.Start(context.Background())

	g.gamestates[gameID] = gs

	return gs, nil
}

func (g *GameService) gameStateFromPersistence(ctx context.Context, game models.Game, moves *[]models.GameMove, hashes *[]models.GameHistoryHash) (*gameplay.GameState, error) {
	var whiteID, blackID uuid.UUID
	if game.GuestBlackID.IsValue() && game.GuestWhiteID.IsValue() {
		whiteID, blackID = game.GuestWhiteID.MustGet(), game.GuestBlackID.MustGet()
	}
	if game.WhiteID.IsValue() && game.BlackID.IsValue() {
		whiteID, blackID = game.WhiteID.MustGet(), game.BlackID.MustGet()
	}

	whiteUsername, blackUsername := "guest", "guest"

	if game.Rated {
		wn, err5 := g.usrRdr.GetUsername(ctx, whiteID)

		bn, err6 := g.usrRdr.GetUsername(ctx, blackID)
		if err5 != nil || err6 != nil {
			g.log.Error("failed to get usernames", slog.Any("error", errors.Join(err5, err6)))
		}

		whiteUsername, blackUsername = wn, bn
	}

	players := [2]gameplay.Player{
		{ID: whiteID, Username: whiteUsername, Color: pb.Color_COLOR_WHITE, Guest: !game.Rated},
		{ID: blackID, Username: blackUsername, Color: pb.Color_COLOR_BLACK, Guest: !game.Rated},
	}

	gtc := &pb.GameTimeControl{ClockMs: game.TimeControlClockMS, IncrementMs: game.TimeControlIncrementMS}

	thresholds := []gameplay.CategoryThreshold{}
	for _, x := range g.categoryThresholds {
		thresholds = append(thresholds, gameplay.CategoryThreshold{
			UpperLimit:   x.upperLimit,
			TimeCategory: x.timeCategory,
		})
	}

	gs, err := gameplay.NewGameState(
		game.ID,
		players,
		gtc,
		thresholds,
		g.gameEvent,
		gameplay.WithFEN(game.Fen),
		gameplay.WithRated(game.Rated),
		gameplay.WithGameVariant(g.gameVariantIDToProto(game.GameVariantID)),
		gameplay.WithGameTimeKind(g.gameTimeKindIDToProto(game.GameTimeKindID)),
		gameplay.WithGameTimeCategory(g.gameTimeCategoryIDToProto(game.GameTimeCategoryID)),
		gameplay.WithGameState(g.gameStateIDToProto(game.GameStateID)),
		gameplay.WithReconnectTimeout(time.Duration(game.ReconnectTimeoutMS)*time.Millisecond),
		gameplay.WithFirstMoveTimeoutOpt(time.Duration(game.FirstMoveTimeoutMS)*time.Millisecond),
		gameplay.WithLastMove(game.LastMove.Ptr()),
		gameplay.WithStartTime(game.StartTime.Ptr()),
		gameplay.WithEndTime(game.EndTime.Ptr()),
		gameplay.WithVersion(int(game.Version)),
	)
	if err != nil {
		return nil, err
	}

	gs.Chess.Repetitions = uint16(game.Repetitions)

	if hashes != nil && len(*hashes) > 0 {
		gs.Chess.HistoryHashes = make([]uint64, len(*hashes))
		for i, hash := range *hashes {
			gs.Chess.HistoryHashes[i] = uint64(hash.Hash)
		}
	}

	var gameMoves []*pb.GameMove

	if moves != nil && len(*moves) > 0 {
		gameMoves = make([]*pb.GameMove, len(*moves))

		for i, m := range *moves {
			move := &pb.GameMove{
				Fen: m.Fen,
			}
			if m.Uci != "" {
				move.Uci = &m.Uci
			}

			if m.San != "" {
				move.San = &m.San
			}

			if m.Lan != "" {
				move.Lan = &m.Lan
			}

			if m.PlayedAt.IsValue() {
				move.PlayedAt = timestamppb.New(m.PlayedAt.MustGet())
			}

			gameMoves[i] = move
		}
	}

	gs.GameMoves = gameMoves

	if game.GameResultID.IsValue() {
		gs.GameResult = g.gameResultIDToProto(game.GameResultID.MustGet())
	}

	if game.GameResultStatusID.IsValue() {
		gs.GameResultStatus = g.gameResultStatusIDToProto(game.GameResultStatusID.MustGet())
	}

	return gs, nil
}
