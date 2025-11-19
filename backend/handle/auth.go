package handle

import (
    "github.com/gin-gonic/gin"
    "github.com/bugoutianzhen123/SoftwareConstructionExp/domain"
    "github.com/bugoutianzhen123/SoftwareConstructionExp/service"
    "github.com/bugoutianzhen123/SoftwareConstructionExp/auth"
)

type AuthHandlers struct { svc *service.Service }

func NewAuthHandlers(s *service.Service) *AuthHandlers { return &AuthHandlers{svc: s} }

type registerBody struct { Name string `json:"name"`; Email string `json:"email"`; Password string `json:"password"`; Role domain.Role `json:"role"`; Skills []string `json:"skills"` }
type loginBody struct { Email string `json:"email"`; Password string `json:"password"` }

func (h *AuthHandlers) Register(c *gin.Context) {
    var b registerBody
    if err := c.ShouldBindJSON(&b); err != nil { c.JSON(400, gin.H{"error":"invalid json"}); return }
    u, err := h.svc.Register(b.Name, b.Email, b.Password, b.Role, b.Skills)
    if err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
    c.JSON(201, u)
}

func (h *AuthHandlers) Login(c *gin.Context) {
    var b loginBody
    if err := c.ShouldBindJSON(&b); err != nil { c.JSON(400, gin.H{"error":"invalid json"}); return }
    u, err := h.svc.Login(b.Email, b.Password)
    if err != nil { c.JSON(401, gin.H{"error": err.Error()}); return }
    tok, err := auth.GenerateToken(u)
    if err != nil { c.JSON(500, gin.H{"error": "token error"}); return }
    c.JSON(200, gin.H{"token": tok})
}