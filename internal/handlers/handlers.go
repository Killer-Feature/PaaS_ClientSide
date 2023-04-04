package handlers

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"net/netip"

	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
)

type Handler struct {
	logger *zap.Logger
	u      internal.Usecase
}

func NewHandler(logger *zap.Logger, u internal.Usecase) *Handler {
	return &Handler{logger: logger, u: u}
}

//go:embed dist
var ui embed.FS

func (h *Handler) Register(s *echo.Echo) {
	// Register http handlers
	s.GET("/hello", h.GetHello)
	s.GET("/api/getClusterNodes", h.GetClusterNodes)
	s.POST("/api/addNode", h.AddNode)
	s.POST("/api/addNodeToCluster", h.AddNodeToCluster)
	s.POST("/api/removeNode", h.RemoveNode)
	s.POST("/api/removeNodeFromCluster", h.RemoveNodeFromCluster)
	s.POST("/api/addResource", h.RemoveNodeFromCluster)

	fsys, err := fs.Sub(ui, "dist")
	if err != nil {
		log.Fatal(err)
	}

	s.GET("/*", echo.WrapHandler(http.FileServer(http.FS(fsys))))
}

func (h *Handler) GetHello(c echo.Context) error {

	h.logger.Info("request received", zap.String("host", c.Request().RemoteAddr), zap.String("command", c.QueryParam("command")))

	command, err := h.u.ExecCommand(c.QueryParam("command"))
	if err != nil {
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	return c.HTML(http.StatusOK, string(command))
}

func (h *Handler) GetClusterNodes(c echo.Context) error {
	nodes, err := h.u.GetClusterNodes(c.Request().Context())
	if err != nil {
		return c.HTML(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, nodes)
}

type InputNode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

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
	return ctx.HTML(http.StatusInternalServerError, "not implemented")
}

type ResourceType struct {
	Type string `json:"type"`
}

func (h *Handler) AddResource(ctx echo.Context) error {
	rType := ResourceType{}
	if err := ctx.Bind(&rType); err != nil {
		h.logger.Error("error occurred during parsing ResourceType", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err := h.u.AddResource(ctx.Request().Context(), ConvertResourceString(rType.Type))
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func ConvertResourceString(resource string) internal.ResourceType {
	switch resource {
	case "postgres":
		return internal.Postgres
	case "redis":
		return internal.Redis
	default:
	}
	return internal.Undefined
}
