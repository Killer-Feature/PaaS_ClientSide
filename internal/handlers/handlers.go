package handlers

import (
	"context"
	"embed"
	"github.com/Killer-Feature/PaaS_ClientSide/internal"
	"github.com/gorilla/websocket"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/netip"
	"time"
)

const (
	READ_BUFSIZE        = 1024
	WRITE_BUFSIZE       = 1024
	SESSION_COOKIE_NAME = "session"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  READ_BUFSIZE,
	WriteBufferSize: WRITE_BUFSIZE,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	logger *zap.Logger
	u      internal.Usecase
}

func NewHandler(logger *zap.Logger, u internal.Usecase) *Handler {
	return &Handler{logger: logger, u: u}
}

// Embedding static files to add it in binary
//
//go:embed dist/assets
var ui embed.FS

//go:embed dist/index.html
var htmlPage string

// Register func receives echo server and register all http handlers
func (h *Handler) Register(s *echo.Echo) {
	// Register http handlers
	s.GET("/hello", h.GetHello)

	s.POST("/api/login", h.Login)

	s.GET("/api/logout", h.Logout, h.AuthMW)

	s.GET("/api/getClusterNodes", h.GetClusterNodes, h.AuthMW)

	s.POST("/api/addNode", h.AddNode, h.AuthMW)
	s.POST("/api/addNodeToCluster", h.AddNodeToCluster, h.AuthMW)

	s.POST("/api/removeNode", h.RemoveNode, h.AuthMW)
	s.POST("/api/removeNodeFromCluster", h.RemoveNodeFromCluster, h.AuthMW)

	s.POST("/api/addResource", h.AddResource, h.AuthMW)
	s.POST("/api/removeResource", h.RemoveResource, h.AuthMW)
	s.GET("/api/getResources", h.GetResources, h.AuthMW)

	s.GET("/api/getAdminConfig", h.GetAdminConfig, h.AuthMW)

	s.GET("/api/getServices", h.GetServices, h.AuthMW)

	s.GET("/api/getProgress", h.GetProgress)

	fsys, err := fs.Sub(ui, "dist")
	if err != nil {
		log.Fatal(err)
	}

	s.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(fsys))))
	s.GET("/*", func(c echo.Context) error {
		return c.HTML(http.StatusOK, htmlPage)
	})
}

// GetHello func is deprecated
func (h *Handler) GetHello(c echo.Context) error {

	h.logger.Info("request received", zap.String("host", c.Request().RemoteAddr), zap.String("command", c.QueryParam("command")))

	command, err := h.u.ExecCommand(c.QueryParam("command"))
	if err != nil {
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	return c.HTML(http.StatusOK, string(command))
}

// GetClusterNodes returns all cluster nodes in huginn database
func (h *Handler) GetClusterNodes(c echo.Context) error {
	nodes, err := h.u.GetClusterNodes(c.Request().Context())
	if err != nil {
		return c.HTML(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, nodes)
}

// InputNode is struct for casting node data in json format
type InputNode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// AddNodeToCluster adds new node to huginn database and checks is it available via ssh connecting
func (h *Handler) AddNodeToCluster(ctx echo.Context) error {
	nodeData := NodeID{}
	if err := ctx.Bind(&nodeData); err != nil {
		h.logger.Error("error occurred during parsing nodeData", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	nodeID, err := h.u.AddNodeToCurrentCluster(ctx.Request().Context(), nodeData.ID)

	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, nodeID)
}

func (h *Handler) AddNode(ctx echo.Context) error {
	nodeData := InputNode{}
	if err := ctx.Bind(&nodeData); err != nil {
		h.logger.Error("error occurred during parsing nodeData", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	parsedIP, err := netip.ParseAddrPort(nodeData.IP)
	if err != nil {
		return ctx.HTML(http.StatusBadRequest, err.Error())
	}
	nodeID, err := h.u.AddNode(ctx.Request().Context(), internal.FullNode{
		Name:     nodeData.Name,
		IP:       parsedIP,
		Login:    nodeData.Login,
		Password: nodeData.Password,
	})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, nodeID)
}

type NodeID struct {
	ID int `json:"id"`
}

func (h *Handler) RemoveNode(ctx echo.Context) error {
	nodeData := NodeID{}
	if err := ctx.Bind(&nodeData); err != nil {
		h.logger.Error("error occurred during parsing nodeData", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	err := h.u.RemoveNode(ctx.Request().Context(), nodeData.ID)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) RemoveNodeFromCluster(ctx echo.Context) error {
	nodeData := NodeID{}
	if err := ctx.Bind(&nodeData); err != nil {
		h.logger.Error("error occurred during parsing nodeData", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	nodeID, err := h.u.RemoveNodeFromCurrentCluster(ctx.Request().Context(), nodeData.ID)

	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, nodeID)
}

type ResourceData struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func (h *Handler) AddResource(ctx echo.Context) error {
	rData := ResourceData{}
	if err := ctx.Bind(&rData); err != nil {
		h.logger.Error("error occurred during parsing ResourceData", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err := h.u.AddResource(ctx.Request().Context(), ConvertResourceTypeToString(rData.Type), rData.Name)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) RemoveResource(ctx echo.Context) error {
	rData := ResourceData{}
	if err := ctx.Bind(&rData); err != nil {
		h.logger.Error("error occurred during parsing ResourceData", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err := h.u.RemoveResource(ctx.Request().Context(), ConvertResourceTypeToString(rData.Type), rData.Name)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func ConvertResourceTypeToString(resource string) internal.ResourceType {
	switch resource {
	case "postgres":
		return internal.Postgres
	case "redis":
		return internal.Redis
	case "prometheus":
		return internal.Prometheus
	case "grafana":
		return internal.Grafana
	case "nginx-ingress-controller":
		return internal.NginxIngressController
	case "metallb":
		return internal.MetalLB
	default:
	}
	return internal.Undefined
}

func (h *Handler) GetAdminConfig(ctx echo.Context) error {
	conf, err := h.u.GetAdminConfig(ctx.Request().Context(), 1)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, conf)
}

func (h *Handler) GetResources(ctx echo.Context) error {
	resources, err := h.u.GetResources(ctx.Request().Context())
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}
	if resources == nil {
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.JSON(http.StatusOK, resources)
}

func (h *Handler) GetProgress(ctx echo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		h.logger.Error("get progress request error", zap.String("error", err.Error()))
		return nil
	}

	_ = h.u.GetProgress(ctx.Request().Context(), ws)
	return nil
}

func (h *Handler) GetServices(ctx echo.Context) error {
	resources, err := h.u.GetServices(ctx.Request().Context())
	if err != nil {
		h.logger.Error("error getting services", zap.Error(err))
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}
	if resources == nil {
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.JSON(http.StatusOK, resources)
}

func (h *Handler) Login(ctx echo.Context) error {
	if IsAdmin(ctx) {
		return ctx.HTML(http.StatusConflict, "already authorized")
	}

	var loginData internal.LoginData
	if err := ctx.Bind(&loginData); err != nil {
		h.logger.Error("error occurred during parsing LoginData", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	sessionKey, err := h.u.Login(context.Background(), loginData)
	if err != nil {
		h.logger.Error("error during login request", zap.Error(err))
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}

	host, _, _ := net.SplitHostPort(ctx.Request().Host)
	ctx.SetCookie(&http.Cookie{
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   false,
		HttpOnly: true,
		Name:     SESSION_COOKIE_NAME,
		Value:    sessionKey,
		Domain:   host,
		Path:     "/",
	})

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) Logout(ctx echo.Context) error {
	sessionCookie, err := ctx.Request().Cookie(SESSION_COOKIE_NAME)

	if err != nil {
		return ctx.HTML(http.StatusUnauthorized, "auth required")
	}

	err = h.u.Logout(ctx.Request().Context(), sessionCookie.Value)

	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}

	host, _, _ := net.SplitHostPort(ctx.Request().Host)
	ctx.SetCookie(&http.Cookie{
		Expires: time.Now().Add(-time.Hour),
		Name:    SESSION_COOKIE_NAME,
		Domain:  host,
		Path:    "/",
	})

	return ctx.NoContent(http.StatusOK)
}
