package go2colab

type Repo struct {
	Name                string   `json:"name"`
	Owner               string   `json:"owner"`
	Url                 string   `json:"url"`
	Host                string   `json:"host"`
	UrlPath             string   `json:"url_path"`
	SystemPath          string   `json:"system_path"`
	GoVersion           string   `json:"go_version"`
	UseLatestReleaseTag bool     `json:"use_latest_release_tag"`
	ReleaseTag          TagInfo  `json:"release_tag"`
	TutorialPaths       []string `json:"tutorial_paths"`
	Head                Commit   `json:"head"`
}

type Tutorial struct {
	Name      string `json:"name"`
	Output    string `json:"output"`
	Docstring string `json:"docstring"`
	Path      string `json:"path"`
	Source    string `json:"source"`
}

type Commit struct {
	Hash    string `json:"hash"`
	Date    string `json:"date"`
	Message string `json:"message"`
	Author  string `json:"author"`
}

type Tag struct {
	Name    string `json:"name"`
	Date    string `json:"date"`
	Message string `json:"message"`
	Commit  Commit `json:"commit"`
}
