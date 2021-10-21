package go2colab

// Notebook is a struct that represents an entire iPython notebook
type Notebook struct {
	Metadata      Metadata `json:"metadata"`
	Nbformat      int      `json:"nbformat"`
	NbformatMinor int      `json:"nbformat_minor"`
	Cells         []Cell   `json:"cells"`
	RepoPath      string   `json:"-"`
}

type Metadata struct {
	Name       string     `json:"name"`
	Colab      Colab      `json:"colab"`
	KernelSpec KernelSpec `json:"kernelspec"`
}

// Cell is a struct that represents iPython notebook code cells
type Cell struct {
	CellType       string       `json:"cell_type"`
	Source         string       `json:"source"`
	Metadata       CellMetadata `json:"metadata"`
	ExecutionCount int          `json:"execution_count"`
	Outputs        []Output     `json:"outputs"`
}

type Colab struct {
	Name              string       `json:"name"`
	Provenance        []Provenance `json:"provenance"`
	CollapsedSections []string     `json:"collapsed_section"`
}

type KernelSpec struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type Provenance struct {
	Name string `json:"name"`
}

type CellMetadata struct {
	Id        string   `json:"id"`
	Collapsed bool     `json:"collapsed"`
	Deletable bool     `json:"deletable"`
	Format    string   `json:"format"` //mime type
	Name      string   `json:"name"`
	Tags      []string `json:"tags"`
}

type Output struct {
	OutputType string   `json:"output_type"`
	Name       string   `json:"name"`
	Text       []string `json:"text"`
}
