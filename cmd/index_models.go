package cmd

type listIndexVersionResp struct {
	IndexSettings `json:"settings"`
}

type IndexSettings struct {
	Index `json:"index"`
}

type Index struct {
	IndexVersion `json:"version"`
}

type IndexVersion struct {
	Created string `json:"created_string"`
}

type indexTemplateLegacyResp map[string]indexTemplateSettings

type indexTemplateSettings struct {
	Patterns []string `json:"index_patterns"`
	Order    int      `json:"order"`
	Version  int      `json:"version"`
}

type CatIndexResp struct {
	Index         string `json:"index"`
	PrimaryShards string `json"pri"`
	ReplicaShards string `json:"rep"`
	Docs          string `json:"docs.count"`
	DeletedDocs   string `json:"docs.deleted"`
	Date          string `json:"creation.date"`
	StoreSize     string `json:"store.size"`
}
