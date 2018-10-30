package ctlDoc

import (
    "gitee.com/johng/gf-home/app/lib/doc"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/gfile"
    "net/http"
)

// 文档首页
func Index(r *ghttp.Request) {
    if r.IsAjaxRequest() {
        serveMarkdownAjax(r)
        return
    }
    path := r.Get("path")
    if path == "" {
        if r.URL.RawQuery != "" {
            r.Response.RedirectTo("/index?" + r.URL.RawQuery)
        } else {
            r.Response.RedirectTo("/index")
        }
        return
    }
    config := g.Config()
    // 如果是静态文件请求，那么表示Web Server没有找到该文件，那么直接404，本接口不支持待后缀的静态文件处理。
    // 由于路由规则比较宽，这里也会有未存在的静态文件请求匹配进来。
    if gfile.Ext(path) != "" {
        r.Response.WriteStatus(http.StatusNotFound)
        return
    }
    // 菜单内容
    baseTitle := config.GetString("doc.title")
    title     := libDoc.GetTitleByPath(path)
    if title == "" {
        title = "404 NOT FOUND"
    }
    title += " - " + config.GetString("doc.title")
    // markdown内容
    mdMainContent       := libDoc.GetMarkdown(path)
    mdMainContentParsed := libDoc.ParseMarkdown(mdMainContent)
    r.Response.WriteTpl("doc/index.html", g.Map {
        "title"               : title,
        "baseTitle"           : baseTitle,
        "mdMenuContentParsed" : libDoc.GetParsed("menus"),
        "mdMainContentParsed" : mdMainContentParsed,
        "mdMainContent"       : mdMainContent,
    })
}

// 文档更新hook
func UpdateHook(r *ghttp.Request) {
    raw    := r.GetRaw()
    j, err := gjson.DecodeToJson(raw)
    if err != nil {
        panic(err)
    }
    if j != nil && j.GetString("password") == g.Config().GetString("doc.hook") {
        libDoc.UpdateDocGit()
    }
    r.Response.Write("ok")
}

// 搜索文档
func Search(r *ghttp.Request) {
    r.Response.WriteJson(g.Map{
        "code" : 1,
        "msg"  : "",
        "data" : libDoc.SearchMdByKey(r.GetString("key")),
    })
}

// 处理ajax请求
func serveMarkdownAjax(r *ghttp.Request) {
    r.Response.WriteJson(g.Map{
        "code" : 1,
        "msg"  : "",
        "data" : libDoc.GetMarkdown(r.Get("path", "index")),
    })
}