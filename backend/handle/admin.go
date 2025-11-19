package handle

import (
    "github.com/gin-gonic/gin"
    "github.com/bugoutianzhen123/SoftwareConstructionExp/service"
)

type AdminHandlers struct { svc *service.Service }

func NewAdminHandlers(s *service.Service) *AdminHandlers { return &AdminHandlers{svc: s} }

func (h *AdminHandlers) Stats(c *gin.Context) {
    c.JSON(200, h.svc.Stats())
}