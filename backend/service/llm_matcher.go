package service

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "github.com/tmc/langchaingo/llms"
    "github.com/tmc/langchaingo/llms/openai"
    "github.com/bugoutianzhen123/SoftwareConstructionExp/domain"
)

type LLMMatcher struct{ llm llms.Model }

func NewLLMMatcher() (domain.Matcher, error) {
    model := os.Getenv("SC_LLM_MODEL")
    if model == "" { model = "gpt-3.5-turbo" }
    baseURL := os.Getenv("SC_LLM_BASE_URL")
    var opts []openai.Option
    opts = append(opts, openai.WithModel(model))
    if baseURL != "" { opts = append(opts, openai.WithBaseURL(baseURL)) }
    m, err := openai.New(opts...)
    if err != nil { return nil, err }
    return &LLMMatcher{llm: m}, nil
}

func (l *LLMMatcher) Match(student *domain.User, projects []*domain.Project) []domain.MatchResult {
    var out []domain.MatchResult
    for _, p := range projects {
        prompt := buildPrompt(student, p)
        ctx := context.Background()
        resp, err := l.llm.Call(ctx, prompt, llms.WithTemperature(0))
        if err != nil || resp == "" {
            out = append(out, domain.MatchResult{Project: p, Score: 0, Reason: "LLM错误或无响应"})
            continue
        }
        var r struct{ Score float64 `json:"score"`; Reason string `json:"reason"` }
        txt := resp
        if e := json.Unmarshal([]byte(txt), &r); e != nil {
            r.Score = 0
            r.Reason = strings.TrimSpace(txt)
        }
        if r.Score < 0 { r.Score = 0 } else if r.Score > 1 { r.Score = 1 }
        out = append(out, domain.MatchResult{Project: p, Score: r.Score, Reason: r.Reason})
    }
    return out
}

func buildPrompt(stu *domain.User, p *domain.Project) string {
    return fmt.Sprintf("请基于以下信息进行科研需求与学生能力的语义匹配，并返回JSON：{\"score\":0..1浮点, \"reason\":中文理由}. 学生技能: %s; 项目标题: %s; 项目描述: %s; 项目要求: %s", strings.Join(stu.Skills, ", "), p.Title, p.Description, strings.Join(p.Requirements, ", "))
}