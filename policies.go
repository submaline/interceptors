package interceptors

// AuthPolicy 明示的に"funcFullPath": falseとしない限りデフォルトでtrue(認証が必要)になります。
type AuthPolicy map[string]bool
