package service

import (
    "strings"
    "github.com/bugoutianzhen123/SoftwareConstructionExp/domain"
)

type SimpleMatcher struct{}

func normalize(words []string) []string {
    r := make([]string, 0, len(words))
    for _, w := range words {
        w = strings.TrimSpace(strings.ToLower(w))
        if w != "" { r = append(r, w) }
    }
    return r
}

func set(items []string) map[string]struct{} {
    m := map[string]struct{}{}
    for _, it := range items { m[it] = struct{}{} }
    return m
}

func (SimpleMatcher) Match(student *domain.User, projects []*domain.Project) []domain.MatchResult {
    ss := set(normalize(student.Skills))
    var out []domain.MatchResult
    for _, p := range projects {
        req := set(normalize(p.Requirements))
        var matched []string
        for k := range req { if _, ok := ss[k]; ok { matched = append(matched, k) } }
        score := 0.0
        if len(req) > 0 { score = float64(len(matched)) / float64(len(req)) }
        reason := ""
        if len(matched) > 0 { reason = "匹配技能: " + strings.Join(matched, ", ") } else { reason = "未找到直接匹配技能，建议补充相关能力" }
        out = append(out, domain.MatchResult{Project: p, Score: score, Reason: reason})
    }
    return out
}
