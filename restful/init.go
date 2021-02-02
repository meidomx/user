package restful

import scaffold "github.com/moetang/webapp-scaffold"

var webscaff *scaffold.WebappScaffold

func InitRestful(w *scaffold.WebappScaffold) {
	webscaff = w

	initApi(webscaff.GetGin())
}
