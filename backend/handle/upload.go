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
    cu := currentUser(c)
    if cu == nil { c.JSON(401, gin.H{"error":"未认证"}); return }
    app := h.svc.Repo().GetApplication(appID)
    if app == nil { c.JSON(404, gin.H{"error":"申请不存在"}); return }
    if app.Status != "approved" { c.JSON(403, gin.H{"error":"仅已通过的申请可上传文档"}); return }
    if cu.Role != domain.RoleStudent || cu.ID != app.StudentID { c.JSON(403, gin.H{"error":"无权上传该申请的文档"}); return }
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
    app := h.svc.Repo().GetApplication(appID)
    if app == nil { c.JSON(404, gin.H{"error":"申请不存在"}); return }
    cu := currentUser(c)
    if cu == nil { c.JSON(401, gin.H{"error":"未认证"}); return }
    if cu.Role == domain.RoleStudent && cu.ID != app.StudentID { c.JSON(403, gin.H{"error":"无权查看该申请文档"}); return }
    if cu.Role == domain.RoleTeacher {
        proj := h.svc.Repo().GetProject(app.ProjectID)
        if proj == nil || proj.TeacherID != cu.ID { c.JSON(403, gin.H{"error":"无权查看该申请文档"}); return }
    }
    docs := h.svc.ListDocumentsByApplication(appID)
    c.JSON(200, docs)
}

func (h *UploadHandlers) Download(c *gin.Context) {
    idStr := c.Query("id")
    if idStr == "" { c.JSON(400, gin.H{"error":"缺少id"}); return }
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil { c.JSON(400, gin.H{"error":"id格式错误"}); return }
    doc := h.svc.Repo().GetDocument(id)
    if doc == nil { c.JSON(404, gin.H{"error":"文档不存在"}); return }
    app := h.svc.Repo().GetApplication(doc.ApplicationID)
    if app == nil { c.JSON(404, gin.H{"error":"申请不存在"}); return }
    cu := currentUser(c)
    if cu == nil { c.JSON(401, gin.H{"error":"未认证"}); return }
    if cu.Role == domain.RoleStudent && cu.ID != app.StudentID { c.JSON(403, gin.H{"error":"无权下载该文档"}); return }
    if cu.Role == domain.RoleTeacher {
        proj := h.svc.Repo().GetProject(app.ProjectID)
        if proj == nil || proj.TeacherID != cu.ID { c.JSON(403, gin.H{"error":"无权下载该文档"}); return }
    }
    if _, err := os.Stat(doc.Path); err != nil { c.JSON(404, gin.H{"error":"文件不存在"}); return }
    c.FileAttachment(doc.Path, doc.Name)
}
