package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"homelab-dashboard/internal/service"
)

type ClusterHandler struct {
	svc *service.ClusterService
}

func NewClusterHandler(svc *service.ClusterService) *ClusterHandler {
	return &ClusterHandler{svc: svc}
}

func (h *ClusterHandler) GetNodes(c echo.Context) error {
	nodes, err := h.svc.GetNodes(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, nodes)
}

func (h *ClusterHandler) GetNamespaces(c echo.Context) error {
	namespaces, err := h.svc.GetNamespaces(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, namespaces)
}

func (h *ClusterHandler) GetPods(c echo.Context) error {
	ns := c.Param("ns")
	pods, err := h.svc.GetPods(c.Request().Context(), ns)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, pods)
}

func (h *ClusterHandler) GetServices(c echo.Context) error {
	ns := c.Param("ns")
	services, err := h.svc.GetServices(c.Request().Context(), ns)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, services)
}
