package auth


// 路由权限标识表
var PathTag = map[string]string{
	// // 超管
	// "admin:/api/dataform/v1/summary":                "150200",
	// "admin:/api/dataform/v1/trend":                  "150200",
	// "admin:/api/dataform/v1/dealAnalysisForDate":    "150300",
	// "admin:/api/dataform/v1/dealAnalysisForShop":    "150300",
	// "admin:/api/dataform/v1/viewAnalysisForDate":    "150400",
	// "admin:/api/dataform/v1/viewAnalysisForShop":    "150400",
	// "admin:/api/dataform/v1/productAnalysisForDate": "150500",
	// "admin:/api/dataform/v1/productAnalysisForShop": "150500",
	// "admin:/api/dataform/v1/sourceAnalysisForDate":  "150600",
	// "admin:/api/dataform/v1/sourceAnalysisForShop":  "150600",
	// "admin:/api/dataform/v1/sourceList":             "150600",
	
	// 普通
	"/api/dataform/v1/summary":                           "270600",
	"/api/dataform/v1/trend":                             "270600",
	"/api/dataform/v1/dealAnalysisForDate":               "270700",
	"/api/dataform/v1/dealAnalysisForShop":               "270700",
	"/api/dataform/v1/dealDetailByShop":                  "270700",
	"/api/dataform/v1/dealDetailByDate":                  "270700",
	"/api/dataform/v1/viewAnalysisForDate":               "270800",
	"/api/dataform/v1/viewAnalysisForShop":               "270800",
	"/api/dataform/v1/viewDetailByShop":                  "270800",
	"/api/dataform/v1/viewDetailByDate":                  "270800",
	"/api/dataform/v1/pageList":                          "270800",
	"/api/dataform/v1/productAnalysisForDate":            "270900",
	"/api/dataform/v1/productAnalysisForShop":            "270900",
	"/api/dataform/v1/productDetailByShop":               "270900",
	"/api/dataform/v1/productDetailByDate":               "270900",
	"/api/dataform/v1/category":                          "270900",
	"/api/dataform/v1/sourceAnalysisForDate":             "271000",
	"/api/dataform/v1/sourceAnalysisForShop":             "271000",
	"/api/dataform/v1/sourceDetailByShop":                "271000",
	"/api/dataform/v1/sourceDetailByDate":                "271000",
	"/api/dataform/v1/sourceList":                        "271000",
	"/api/dataform/v1/sellTrend":                         "270500",
	"/api/dataform/v1/goodsTypeSell/money":               "270500",
	"/api/dataform/v1/goodsTypeSell/amount":              "270500",
	"/api/dataform/v1/goodsTypeSell/productRank":         "270500",
	"/api/dataform/v1/goodsTypeSell/profit":              "270500",
	"/api/dataform/v1/asset":                             "260500",
	"/api/dataform/v1/asset/detail":                      "260501",
	"/api/dataform/v1/exportList/productAnalysisForDate": "270900",
	"/api/dataform/v1/exportList/productAnalysisForShop": "270900",
}

func authCheck(path string, permissions []string) bool {
	tag, ok := PathTag[path]
	if !ok {
		return false
	}
	for _, p := range permissions {
		if p == tag {
			return true
		}
	}
	return false
}
