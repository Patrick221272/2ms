package plugins

type Plugin struct {
	name string

	url string

	email string

	token string
}

type Plugins struct {
	plugins map[string]Plugin
}

type Content struct {
	Content     string
	Source      string
	OriginalUrl string
}

func NewPlugins() *Plugins {
	plugins := make(map[string]Plugin)
	return &Plugins{plugins: plugins}
}

func (P *Plugins) AddPlugin(name string, url string, email string, token string) {
	P.plugins["name"] = Plugin{name: name, url: url, email: email, token: token}
}

func (P *Plugins) RunPlugins() []Content {
	contents := []Content{}
	for _, p := range P.plugins {
		switch p.name {
		case "confluence":
			contents = append(contents, p.RunPlugin()...)
		default:
		}
	}
	return contents
}
