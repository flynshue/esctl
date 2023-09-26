package cmd

type listIndexSettingsResp struct {
	IndexSettings `json:"settings"`
}

type IndexSettings struct {
	Index `json:"index"`
}

type Index struct {
	IndexVersion     `json:"version"`
	Lifecycle        `json:"lifecycle"`
	NumberOfShards   string `json:"number_of_shards"`
	AutoExpand       string `json:"auto_expand_replicas"`
	NumberOfReplicas string `json:"number_of_replicas"`
}

type Lifecycle struct {
	Name          string `json:"name"`
	RolloverAlias string `json:"rollover_alias"`
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
	PrimaryShards string `json:"pri"`
	ReplicaShards string `json:"rep"`
	Docs          string `json:"docs.count"`
	DeletedDocs   string `json:"docs.deleted"`
	Date          string `json:"creation.date"`
	PriStoreSize  string `json:"pri.store.size"`
	StoreSize     string `json:"store.size"`
}
