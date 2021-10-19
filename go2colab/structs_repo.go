package go2colab

type Repo struct {
	Name      string   `json:"name"`
	Owner     string   `json:"owner"`
	Host      string   `json:"host"`
	Url       string   `json:"url"`
	GoVersion string   `json:"go_version"`
	Paths     []string `json:"paths"`
	Commits   []Commit `json:"commits"`
	Tags      []Tag    `json:"tags"`
	Release   string   `json:"release"`
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
