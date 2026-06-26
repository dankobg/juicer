package presence

import (
	"context"
	"log/slog"
	"time"

	"github.com/dankobg/juicer/bus"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserPresenceInfo struct {
	ID       string
	Username string
	Guest    bool
}

type PresenceChannelsDiff struct {
	ConnJoined []string
	ConnLeft   []string
	UserJoined []string
	UserLeft   []string
}

type PresenceService struct {
	bus *bus.Bus
	pst PresencePersistor
	log *slog.Logger
}

func NewPresenceService(bus *bus.Bus, pst PresencePersistor, l *slog.Logger) *PresenceService {
	return &PresenceService{
		bus: bus,
		pst: pst,
		log: l,
	}
}

func (p *PresenceService) PublishPresenceDiff(ctx context.Context, channelsDiff PresenceChannelsDiff, userID, connID, username string, guest bool) error {
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
		p.log.Error("Message_PresenceDiff protojson marshal", slog.String("user_id", userID), slog.String("conn_id", connID), slog.Any("error", err))
		return err
	}

	if err := p.bus.Publish(ctx, "presence.diff."+userID, presenceDiffMsgBytes); err != nil {
		p.log.Error("publish Message_PresenceDiff", slog.String("user_id", userID), slog.String("conn_id", connID), slog.Any("error", err))
		return err
	}

	return nil
}

func (p *PresenceService) SendUserPresenceDiffToChannel(ctx context.Context, channelsDiff PresenceChannelsDiff, channel, userID, connID, username string, guest bool) error {
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

	if err := p.bus.Publish(ctx, channel, presenceDiffMsgBytes); err != nil {
		return err
	}

	return nil
}

func (p *PresenceService) SendChannelPresenceStateToConn(ctx context.Context, channel, connID string) error {
	users, err := p.ListUsersInChannel(ctx, channel)
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

	if err := p.bus.Publish(ctx, "conn."+connID, presenceStateMsgBytes); err != nil {
		return err
	}

	return nil
}

func (s *PresenceService) SetPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool, channels []string) (PresenceChannelsDiff, error) {
	return s.pst.SetPresence(ctx, userID, connID, username, guest, channels)
}

func (s *PresenceService) ClearPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) (PresenceChannelsDiff, error) {
	return s.pst.ClearPresence(ctx, userID, connID, username, guest)
}

func (s *PresenceService) RefreshPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) error {
	return s.pst.RefreshPresence(ctx, userID, connID, username, guest)
}

func (s *PresenceService) UserLastSeen(ctx context.Context, userID uuid.UUID) (time.Time, error) {
	return s.pst.UserLastSeen(ctx, userID)
}

func (s *PresenceService) ListUsersInChannel(ctx context.Context, channel string) ([]UserPresenceInfo, error) {
	return s.pst.ListUsersInChannel(ctx, channel)
}

func (s *PresenceService) ListChannelsForUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.pst.ListChannelsForUser(ctx, userID)
}

func (s *PresenceService) UsersCountInChannel(ctx context.Context, channel string) (int64, error) {
	return s.pst.UsersCountInChannel(ctx, channel)
}

func (s *PresenceService) TotalActiveConnsCount(ctx context.Context) (int64, error) {
	return s.pst.TotalActiveConnsCount(ctx)
}
