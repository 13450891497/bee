package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var cmdFix = &Command{
	UsageLine: "fix",
	Short:     "fix the beego application to compatibel with beego 1.6",
	Long: `
As from beego1.6, there's some incompatible code with the old version.

bee fix help to upgrade the application to beego 1.6
`,
}

func init() {
	cmdFix.Run = runFix
}

func runFix(cmd *Command, args []string) int {
	dir, err := os.Getwd()
	if err != nil {
		ColorLog("GetCurrent Path:%s\n", err)
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if err != nil {
			return err
		}
		ColorLog("%s\n", path)
		err = fixFile(path)
		if err != nil {
			ColorLog("fixFile:%s\n", err)
		}
		return err
	})
	return 0
}

var rules = []string{
	"beego.AppName", "beego.BConfig.AppName",
	"beego.RunMode", "beego.BConfig.RunMode",
	"beego.RecoverPanic", "beego.BConfig.RecoverPanic",
	"beego.RouterCaseSensitive", "beego.BConfig.RouterCaseSensitive",
	"beego.BeegoServerName", "beego.BConfig.ServerName",
	"beego.EnableGzip", "beego.BConfig.EnableGzip",
	"beego.ErrorsShow", "beego.BConfig.EnableErrorsShow",
	"beego.CopyRequestBody", "beego.BConfig.CopyRequestBody",
	"beego.MaxMemory", "beego.BConfig.MaxMemory",
	"beego.Graceful", "beego.BConfig.Listen.Graceful",
	"beego.HttpAddr", "beego.BConfig.Listen.HTTPAddr",
	"beego.HttpPort", "beego.BConfig.Listen.HTTPPort",
	"beego.ListenTCP4", "beego.BConfig.Listen.ListenTCP4",
	"beego.EnableHttpListen", "beego.BConfig.Listen.HTTPEnable",
	"beego.EnableHttpTLS", "beego.BConfig.Listen.HTTPSEnable",
	"beego.HttpsAddr", "beego.BConfig.Listen.HTTPSAddr",
	"beego.HttpsPort", "beego.BConfig.Listen.HTTPSPort",
	"beego.HttpCertFile", "beego.BConfig.Listen.HTTPSCertFile",
	"beego.HttpKeyFile", "beego.BConfig.Listen.HTTPSKeyFile",
	"beego.EnableAdmin", "beego.BConfig.Listen.AdminEnable",
	"beego.AdminHttpAddr", "beego.BConfig.Listen.AdminAddr",
	"beego.AdminHttpPort", "beego.BConfig.Listen.AdminPort",
	"beego.UseFcgi", "beego.BConfig.Listen.EnableFcgi",
	"beego.HttpServerTimeOut", "beego.BConfig.Listen.ServerTimeOut",
	"beego.AutoRender", "beego.BConfig.WebConfig.AutoRender",
	"beego.ViewsPath", "beego.BConfig.WebConfig.ViewsPath",
	"beego.StaticDir", "beego.BConfig.WebConfig.StaticDir",
	"beego.StaticExtensionsToGzip", "beego.BConfig.WebConfig.StaticExtensionsToGzip",
	"beego.DirectoryIndex", "beego.BConfig.WebConfig.DirectoryIndex",
	"beego.FlashName", "beego.BConfig.WebConfig.FlashName",
	"beego.FlashSeperator", "beego.BConfig.WebConfig.FlashSeperator",
	"beego.EnableDocs", "beego.BConfig.WebConfig.EnableDocs",
	"beego.XSRFKEY", "beego.BConfig.WebConfig.XSRFKEY",
	"beego.EnableXSRF", "beego.BConfig.WebConfig.EnableXSRF",
	"beego.XSRFExpire", "beego.BConfig.WebConfig.XSRFExpire",
	"beego.TemplateLeft", "beego.BConfig.WebConfig.TemplateLeft",
	"beego.TemplateRight", "beego.BConfig.WebConfig.TemplateRight",
	"beego.SessionOn", "beego.BConfig.WebConfig.Session.SessionOn",
	"beego.SessionProvider", "beego.BConfig.WebConfig.Session.SessionProvider",
	"beego.SessionName", "beego.BConfig.WebConfig.Session.SessionName",
	"beego.SessionGCMaxLifetime", "beego.BConfig.WebConfig.Session.SessionGCMaxLifetime",
	"beego.SessionSavePath", "beego.BConfig.WebConfig.Session.SessionProviderConfig",
	"beego.SessionCookieLifeTime", "beego.BConfig.WebConfig.Session.SessionCookieLifeTime",
	"beego.SessionAutoSetCookie", "beego.BConfig.WebConfig.Session.SessionAutoSetCookie",
	"beego.SessionDomain", "beego.BConfig.WebConfig.Session.SessionDomain",
	"Ctx.Input.CopyBody(", "Ctx.Input.CopyBody(beego.BConfig.MaxMemory",
	".UrlFor(", ".URLFor(",
	".ServeJson(", ".ServeJSON(",
	".ServeXml(", ".ServeXML(",
	".XsrfToken(", ".XSRFToken(",
	".CheckXsrfCookie(", ".CheckXSRFCookie(",
	".XsrfFormHtml(", ".XSRFFormHTML(",
	"beego.UrlFor(", "beego.URLFor(",
	"beego.GlobalDocApi", "beego.GlobalDocAPI",
	"beego.Errorhandler", "beego.ErrorHandler",
	"Output.Jsonp(", "Output.JSONP(",
	"Output.Json(", "Output.JSON(",
	"Output.Xml(", "Output.XML(",
	"Input.Uri()", "Input.URI()",
	"Input.Url()", "Input.URL()",
	"Input.AcceptsHtml()", "Input.AcceptsHTML()",
	"Input.AcceptsXml()", "Input.AcceptsXML()",
	"Input.AcceptsJson()", "Input.AcceptsJSON()",
	"Ctx.XsrfToken()", "Ctx.XSRFToken()",
	"Ctx.CheckXsrfCookie()", "Ctx.CheckXSRFCookie()",
	"session.SessionStore", "session.Store",
	".TplNames", ".TplName",
	"swagger.ApiRef", "swagger.APIRef",
	"swagger.ApiDeclaration", "swagger.APIDeclaration",
	"swagger.Api", "swagger.API",
	"swagger.ApiRef", "swagger.APIRef",
	"toolbox.UrlMap", "toolbox.URLMap",
	"logs.LoggerInterface", "logs.Logger",
	"Input.Request", "Input.Context.Request",
	"Input.Params)", "Input.Params())",
	"httplib.BeegoHttpSettings", "httplib.BeegoHTTPSettings",
	"httplib.BeegoHttpRequest", "httplib.BeegoHTTPRequest",
	".TlsClientConfig", ".TLSClientConfig",
	".JsonBody", ".JSONBody",
	".ToJson", ".ToJSON",
	".ToXml", ".ToXML",
	"beego.Html2str", "beego.HTML2str",
	"beego.AssetsCss", "beego.AssetsCSS",
	"orm.DR_Sqlite", "orm.DRSqlite",
	"orm.DR_Postgres", "orm.DRPostgres",
	"orm.DR_MySQL", "orm.DRMySQL",
	"orm.DR_Oracle", "orm.DROracle",
	"orm.Col_Add", "orm.ColAdd",
	"orm.Col_Minus", "orm.ColMinus",
	"orm.Col_Multiply", "orm.ColMultiply",
	"orm.Col_Except", "orm.ColExcept",
	"GenerateOperatorSql", "GenerateOperatorSQL",
	"OperatorSql", "OperatorSQL",
	"orm.Debug_Queries", "orm.DebugQueries",
	"orm.COMMA_SPACE", "orm.CommaSpace",
}

func fixFile(file string) error {
	rp := strings.NewReplacer(rules...)
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	fixed := rp.Replace(string(content))

	// forword the RequestBody from the replace
	// "Input.Request", "Input.Context.Request",
	fixed = strings.Replace(fixed, "Input.Context.RequestBody", "Input.RequestBody", -1)

	// regexp replace
	pareg := regexp.MustCompile(`(Input.Params\[")(.*)("])`)
	fixed = pareg.ReplaceAllString(fixed, "Input.Param(\"$2\")")
	pareg = regexp.MustCompile(`(Input.Data\[\")(.*)(\"\])(\s)(=)(\s)(.*)`)
	fixed = pareg.ReplaceAllString(fixed, "Input.SetData(\"$2\", $7)")
	// replace the v.Apis in docs.go
	if strings.Contains(file, "docs.go") {
		fixed = strings.Replace(fixed, "v.Apis", "v.APIs", -1)
	}

	err = os.Truncate(file, 0)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(fixed), 0666)
}
