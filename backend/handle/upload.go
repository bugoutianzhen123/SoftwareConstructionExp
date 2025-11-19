package handle

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/bugoutianzhen123/SoftwareConstructionExp/domain"
	"github.com/bugoutianzhen123/SoftwareConstructionExp/service"
	"github.com/gin-gonic/gin"
)

type UploadHandlers struct { svc *service.Service }

func NewUploadHandlers(s *service.Service) *UploadHandlers { return &UploadHandlers{svc: s} }

func (h *UploadHandlers) Upload(c *gin.Context) {
    appIDStr := c.PostForm("application_id")
    if appIDStr == "" { c.JSON(400, gin.H{"error":"缺少application_id"}); return }
    appID, err := strconv.ParseInt(appIDStr, 10, 64)
    if err != nil { c.JSON(400, gin.H{"error":"application_id格式错误"}); return }
    _, err = c.FormFile("file")
    if err != nil { c.JSON(400, gin.H{"error":"缺少file"}); return }
    f, _ := c.FormFile("file")
    dir := filepath.Join("uploads", strconv.FormatInt(appID, 10))
    os.MkdirAll(dir, 0755)
    path := filepath.Join(dir, filepath.Base(f.Filename))
    if err := c.SaveUploadedFile(f, path); err != nil { c.JSON(500, gin.H{"error":"保存失败"}); return }
    d := &domain.Document{ApplicationID: appID, Name: f.Filename, Path: path}
    created, err := h.svc.AddDocument(d)
    if err != nil { c.JSON(500, gin.H{"error":"保存失败"}); return }
    c.JSON(201, created)
}

func (h *UploadHandlers) List(c *gin.Context) {
    appIDStr := c.Query("application_id")
    if appIDStr == "" { c.JSON(400, gin.H{"error":"缺少application_id"}); return }
    appID, err := strconv.ParseInt(appIDStr, 10, 64)
    if err != nil { c.JSON(400, gin.H{"error":"application_id格式错误"}); return }
    docs := h.svc.ListDocumentsByApplication(appID)
    c.JSON(200, docs)
}