package main

import (
	"embed"
	"log"

	"github.com/dankobg/juicer/cmd"
)

//go:embed public/*
var publicFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

func main() {
	if err := cmd.Run(publicFiles, templateFiles); err != nil {
		log.Fatalf("failed to run juicer chess server")
	}
}

// package gameserver

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	myredis "juicer/redis"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"strings"
// 	"syscall"
// 	"time"

// 	"github.com/gobwas/ws"
// 	"github.com/gobwas/ws/wsutil"
// 	"github.com/google/uuid"
// 	"github.com/labstack/echo/v5"
// 	"github.com/labstack/echo/v5/middleware"
// 	kratos "github.com/ory/client-go"
// 	"github.com/redis/go-redis/v9"
// )

// type Message struct {
// 	Type string `json:"type"`
// 	Data string `json:"data"`
// }

// type GameServer struct {
// 	Hub     *Hub
// 	kratosc *kratos.APIClient
// 	rdb     *redis.Client
// }

// type Meta struct {
// 	Session   *kratos.Session
// 	ClientID  string
// 	Anonymous bool
// }

// func (gs *GameServer) handleWS(c echo.Context) error {
// 	ctx := c.Request().Context()
// 	sess := GetSession(ctx)

// 	var clientID string
// 	var anonymous bool

// 	if sess != nil {
// 		clientID = sess.Identity.Id
// 	} else {
// 		anonymous = true
// 		clientID = uuid.NewString()
// 	}

// 	meta := Meta{
// 		Session:   sess,
// 		ClientID:  clientID,
// 		Anonymous: anonymous,
// 	}

// 	cookie := &http.Cookie{
// 		Name:     "juicer_client_id",
// 		Value:    clientID,
// 		Path:     "/",
// 		Expires:  time.Now().AddDate(10, 0, 0),
// 		MaxAge:   -1,
// 		HttpOnly: true,
// 	}
// 	c.SetCookie(cookie)

// 	conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
// 	if err != nil {
// 		return fmt.Errorf("failed to upgrade connection: %w", err)
// 	}

// 	go func() {
// 		defer conn.Close()

// 		gs.Hub.handleNewConnection(conn, clientID)

// 		var (
// 			r       = wsutil.NewReader(conn, ws.StateServerSide)
// 			w       = wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
// 			decoder = json.NewDecoder(r)
// 			encoder = json.NewEncoder(w)
// 		)

// 	readLoop:
// 		for {
// 			hdr, err := r.NextFrame()
// 			if err != nil {
// 				fmt.Printf("r.NextFrame err: %+v\n", err)
// 				break readLoop
// 			}

// 			if hdr.OpCode == ws.OpClose {
// 				fmt.Printf("OpCode == OpClose\n")
// 				break readLoop
// 			}

// 			var req Message
// 			if err := decoder.Decode(&req); err != nil {
// 				fmt.Printf("Decode err: %+v\n", err)
// 			}

// 			if err := gs.handleIncomingWebSocketMessage(conn, req, encoder, meta); err != nil {
// 				fmt.Printf("Handle ws message err: %+v\n", err)
// 				break readLoop
// 			}

// 			if err = w.Flush(); err != nil {
// 				fmt.Printf("Flush err: %+v\n", err)
// 				break readLoop
// 			}
// 		}

// 		gs.Hub.handleDisconnection(clientID)
// 	}()

// 	return nil
// }

// func (gs *GameServer) handleIncomingWebSocketMessage(conn net.Conn, msg Message, encoder *json.Encoder, meta Meta) error {
// 	fmt.Printf("received: %+v\n", msg)

// 	var err error

// 	switch msg.Type {
// 	// Loby related
// 	case "loby:seek_game":
// 		err = gs.onSeekGame(conn, msg.Data, encoder, meta)
// 	case "loby:cancel_seek_game":
// 		err = gs.onCancelSeekGame(conn, msg.Data, encoder, meta)
// 	// Match related
// 	case "match:resign":
// 		err = gs.onResign(conn, msg.Data, encoder, meta)
// 	case "match:offer_draw":
// 		err = gs.onOfferDraw(conn, msg.Data, encoder, meta)
// 	}

// 	errMsg := "unkown message type"
// 	if err != nil {
// 		errMsg = err.Error()
// 	}

// 	resp := Message{
// 		Type: "error",
// 		Data: errMsg,
// 	}

// 	return encoder.Encode(&resp)
// }

// func (gs *GameServer) onSeekGame(conn net.Conn, data string, encoder *json.Encoder, meta Meta) error {
// 	ctx := context.Background()

// 	if err := gs.rdb.Get(ctx, meta.ClientID).Err(); err != nil {
// 		if !errors.Is(err, redis.Nil) {
// 			return fmt.Errorf("player already in queue: %w", err)
// 		}

// 		return err
// 	}

// 	_, err := gs.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
// 		if err := pipe.LPush(ctx, "juicer:seeking", meta.ClientID).Err(); err != nil {
// 			return err
// 		}
// 		if err := pipe.Publish(ctx, "juicer:seeking:joined", meta.ClientID).Err(); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	return err
// }

// func (gs *GameServer) onCancelSeekGame(conn net.Conn, data string, encoder *json.Encoder, meta Meta) error {
// 	ctx := context.Background()

// 	if err := gs.rdb.Get(ctx, meta.ClientID).Err(); err != nil {
// 		if errors.Is(err, redis.Nil) {
// 			return fmt.Errorf("player not in queue: %w", err)
// 		}

// 		return err
// 	}

// 	_, err := gs.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
// 		if err := pipe.LPop(ctx, "juicer:seeking").Err(); err != nil {
// 			return err
// 		}
// 		if err := pipe.Publish(ctx, "juicer:seeking:left", meta.ClientID).Err(); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	return err
// }

// func (gs *GameServer) onOfferDraw(conn net.Conn, data string, encoder *json.Encoder, meta Meta) error {
// 	return nil
// }

// func (gs *GameServer) onResign(conn net.Conn, data string, encoder *json.Encoder, meta Meta) error {
// 	return nil
// }

// func handleHealthReady(c echo.Context) error {
// 	return c.JSON(http.StatusOK, map[string]any{"status": "ready"})
// }

// func handleHealthAlive(c echo.Context) error {
// 	return c.JSON(http.StatusOK, map[string]any{"status": "alive"})
// }

// func RunServer() error {
// 	rdb, err := myredis.New()
// 	if err != nil {
// 		return err
// 	}

// 	gs := &GameServer{
// 		Hub:     NewHub(),
// 		kratosc: newKratosClient(),
// 		rdb:     rdb,
// 	}

// 	e := echo.New()

// 	corsConfig := middleware.CORSConfig{
// 		AllowOrigins:     []string{"https://juicer-dev.xyz:1420"},
// 		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
// 		AllowHeaders:     []string{"Content-Type", "Authorization", "X-CSRF-Token"},
// 		ExposeHeaders:    []string{"Content-Length", "Cache-Control", "Content-Language", "Content-Type", "Content-Range", "Expires", "Last-Modified", "Pragma", "Authorization"},
// 		MaxAge:           96400,
// 		AllowCredentials: false,
// 	}

// 	e.Use(middleware.Recover())
// 	e.Use(middleware.BodyLimit(5 * 1024 * 1024))
// 	e.Use(middleware.CORSWithConfig(corsConfig))
// 	e.Use(gs.AttachSessionData)

// 	e.GET("/api/v1/health/alive", handleHealthAlive)
// 	e.GET("/api/v1/health/ready", handleHealthReady)

// 	e.GET("/ws", gs.handleWS)

// 	go func() {
// 		ctx := context.Background()

// 		pubsub := gs.rdb.PSubscribe(ctx, "juicer:seeking:*")
// 		defer pubsub.Close()

// 		ch := pubsub.Channel()

// 		for msg := range ch {
// 			clientID := msg.Payload
// 			parts := strings.Split(msg.Channel, ":")
// 			action := parts[2]

// 			client, ok := gs.Hub.Cons[clientID]

// 			if !ok {
// 				fmt.Println("NOT OK")
// 				return
// 			}

// 			fmt.Println("CIENT", client)

// 			w := wsutil.NewWriter(client.Conn, ws.StateServerSide, ws.OpText)
// 			encoder := json.NewEncoder(w)

// 			if action == "joined" {
// 				fmt.Println("JOINED")

// 				resp := Message{
// 					Type: "juicer:seeking:joined",
// 					Data: "success",
// 				}
// 				if err := encoder.Encode(&resp); err != nil {
// 					fmt.Printf("failed to encode msg: %+v\n", err)
// 				}
// 			} else if action == "left" {
// 				fmt.Println("LEFT")

// 				resp := Message{
// 					Type: "juicer:seeking:left",
// 					Data: "success",
// 				}
// 				if err := encoder.Encode(&resp); err != nil {
// 					fmt.Printf("failed to encode msg: %+v\n", err)
// 				}
// 			}

// 			length, err := gs.rdb.LLen(ctx, "juicer:seeking").Result()
// 			if err != nil {
// 				fmt.Println(err)

// 				resp := Message{
// 					Type: "juicer:seeking:error",
// 					Data: err.Error(),
// 				}
// 				if err := encoder.Encode(&resp); err != nil {
// 					fmt.Printf("failed to encode msg: %+v\n", err)
// 				}
// 			}

// 			if length > 1 && length%2 == 0 {
// 				ids, err := gs.rdb.LPopCount(ctx, "juicer:seeking", 2).Result()
// 				if err != nil {
// 					fmt.Println(err)

// 					resp := Message{
// 						Type: "juicer:seeking:error",
// 						Data: err.Error(),
// 					}
// 					if err := encoder.Encode(&resp); err != nil {
// 						fmt.Printf("failed to encode msg: %+v\n", err)
// 					}
// 				}

// 				fmt.Println("POPPED IDS", ids)

// 				id1, id2 := ids[0], ids[1]
// 				c1, c2 := gs.Hub.Cons[id1], gs.Hub.Cons[id2]

// 				if err := gs.Hub.SetupGame(c1, c2); err != nil {
// 					fmt.Printf("StartGame err: %+v", err)
// 					return
// 				}

// 				resp := Message{
// 					Type: "juicer:game:started",
// 					Data: "success",
// 				}

// 				bb, err := json.Marshal(&resp)
// 				if err != nil {
// 					panic(err)
// 				}

// 				if err := wsutil.WriteServerText(c1.Conn, bb); err != nil {
// 					fmt.Printf("failed to write ws srv msg: %+v\n", err)
// 				}
// 				if err := wsutil.WriteServerText(c2.Conn, bb); err != nil {
// 					fmt.Printf("failed to write ws srv msg: %+v\n", err)
// 				}
// 			}

// 			if err = w.Flush(); err != nil {
// 				fmt.Printf("Write Flush err: %+v\n", err)
// 			}
// 		}
// 	}()

// 	srv := &http.Server{
// 		Addr:              ":1337",
// 		ReadTimeout:       30 * time.Second,
// 		WriteTimeout:      30 * time.Second,
// 		IdleTimeout:       45 * time.Second,
// 		ReadHeaderTimeout: 45 * time.Second,
// 		Handler:           e,
// 	}

// 	go func() {
// 		fmt.Printf("server running on: http://localhost:1337\n")

// 		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
// 			log.Fatalf("failed to run server: %v\n", err)
// 		}
// 	}()

// 	interrupt := make(chan os.Signal, 1)
// 	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

// 	<-interrupt
// 	fmt.Printf("interrupt received\n")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer func() {
// 		// cleanup
// 		cancel()
// 	}()

// 	if err := srv.Shutdown(ctx); err != nil {
// 		return fmt.Errorf("server shutdown failed: %w", err)
// 	}

// 	fmt.Printf("server shut down\n")

// 	return nil
// }

// // ##################################################################

// func newKratosClient() *kratos.APIClient {
// 	publicURL := os.Getenv("KRATOS_PUBLIC_URL")

// 	c := kratos.NewConfiguration()
// 	c.Servers = kratos.ServerConfigurations{{URL: publicURL}}

// 	return kratos.NewAPIClient(c)
// }

// const (
// 	oryKratosCsrfCookieSuffix  = "csrf_token"
// 	oryKratosSessionCookieName = "ory_kratos_session"

// 	prefixBearer = "Bearer"
// )

// func ExtractAuthHeadersFromRequest(r *http.Request) (csrf *http.Cookie, session *http.Cookie, authHeader, cookieHeader string) {
// 	cookieHeader = r.Header.Get(echo.HeaderCookie)
// 	authHeader = r.Header.Get(echo.HeaderAuthorization)

// 	cookies := r.Cookies()
// 	for _, c := range cookies {
// 		if c != nil {
// 			if ok := strings.HasSuffix(c.Name, string(oryKratosCsrfCookieSuffix)); ok {
// 				csrf = c
// 			}
// 		}
// 	}

// 	sessionCookie, _ := r.Cookie(string(oryKratosSessionCookieName))
// 	if sessionCookie != nil {
// 		session = sessionCookie
// 	}

// 	return
// }

// type contextKey string

// const (
// 	oryKratosSessionCtxKey contextKey = "ory_kratos_session_key"
// )

// func WithSession(ctx context.Context, sess *kratos.Session) context.Context {
// 	return context.WithValue(ctx, oryKratosSessionCtxKey, sess)
// }

// func GetSession(ctx context.Context) *kratos.Session {
// 	sess, ok := ctx.Value(oryKratosSessionCtxKey).(*kratos.Session)
// 	if !ok {
// 		return nil
// 	}

// 	return sess
// }

// func (gs *GameServer) RequireSession(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		sess := GetSession(c.Request().Context())

// 		if sess == nil {
// 			return &echo.HTTPError{Code: 400, Message: "session is invalid or has expired already"}
// 		}

// 		if sess != nil && sess.Active != nil && !*sess.Active {
// 			return &echo.HTTPError{Code: 400, Message: "session is invalid or has expired already"}
// 		}

// 		return next(c)
// 	}
// }

// func (gs *GameServer) RequireAnonymous(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		sess := GetSession(c.Request().Context())

// 		if sess != nil {
// 			return &echo.HTTPError{Code: 400, Message: "must have no session"}
// 		}

// 		return next(c)
// 	}
// }

// func (gs *GameServer) AttachSessionData(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		ctx := c.Request().Context()
// 		_, _, authHeader, cookieHeader := ExtractAuthHeadersFromRequest(c.Request())

// 		var hasAuthHeader bool
// 		var hasCookieHeader bool

// 		if authHeader != "" && strings.HasPrefix(authHeader, prefixBearer) {
// 			hasAuthHeader = true
// 		}

// 		if cookieHeader != "" {
// 			hasCookieHeader = true
// 		}

// 		if !hasAuthHeader && !hasCookieHeader {
// 			return next(c)
// 		}

// 		req := gs.kratosc.FrontendApi.ToSession(ctx).XSessionToken(authHeader).Cookie(cookieHeader)
// 		session, _, err := req.Execute()
// 		if err != nil {
// 			return next(c)
// 		}

// 		if session != nil && session.Active != nil && !*session.Active {
// 			return next(c)
// 		}

// 		if session != nil {
// 			ctxWithSession := WithSession(ctx, session)
// 			c.SetRequest(c.Request().WithContext(ctxWithSession))
// 		}

// 		return next(c)
// 	}
// }
