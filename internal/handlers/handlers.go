package handlers

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"net/netip"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
)

const (
	CLUSTER_ID_PARAM_NAME = "clusterId"
)

type Handler struct {
	logger *zap.Logger
	u      internal.Usecase
}

func NewHandler(logger *zap.Logger, u internal.Usecase) *Handler {
	return &Handler{logger: logger, u: u}
}

//go:embed dist/assets
var ui embed.FS

//go:embed dist/index.html
var htmlPage string

func (h *Handler) Register(s *echo.Echo) {
	// Register http handlers
	s.GET("/hello", h.GetHello)
	s.GET("/api/getClusterNodes", h.GetClusterNodes)

	s.POST("/api/addNode", h.AddNode)
	s.POST("/api/addNodeToCluster", h.AddNodeToCluster)

	s.POST("/api/removeNode", h.RemoveNode)
	s.POST("/api/removeNodeFromCluster", h.RemoveNodeFromCluster)

	s.POST("/api/addResource", h.AddResource)
	s.POST("/api/removeResource", h.RemoveResource)

	s.GET("/api/getAdminConfig", h.GetAdminConfig)
	s.GET("/api/getResources", h.GetResources)

	fsys, err := fs.Sub(ui, "dist")
	if err != nil {
		log.Fatal(err)
	}

	s.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(fsys))))
	s.GET("/*", func(c echo.Context) error {
		return c.HTML(http.StatusOK, htmlPage)
	})
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
	clusterIdStr := ctx.QueryParam(CLUSTER_ID_PARAM_NAME)
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		clusterId = 1
		//return ctx.HTML(http.StatusBadRequest, err.Error())
	}
	conf, err := h.u.GetAdminConfig(ctx.Request().Context(), clusterId)
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
