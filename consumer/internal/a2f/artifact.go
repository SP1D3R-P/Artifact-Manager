package a2f

import (
	"bytes"
	"encoding/json"
)

type Artifact struct {
	Config *ArtifactConf
}

func New(cnf *ArtifactConf) *Artifact {
	return &Artifact{
		Config: cnf,
	}
}

func (a *Artifact) GetName() string { return a.Config.BinaryObject.Name }
func (a *Artifact) Process() error  { return nil }
func (a *Artifact) GetData() []byte {
	buff := bytes.Buffer{}

	parser := json.NewEncoder(&buff)
	parser.SetEscapeHTML(false)
	parser.SetIndent("", " ")

	parser.Encode(a.Config)

	return buff.Bytes()
}

func (a *Artifact) GetConfig() ArtifactConf { return *a.Config }

func (a *Artifact) GetHash() string       { return a.Config.Metadata.Checksum }
func (a *Artifact) GetConfigHash() string { return a.Config.Metadata.Checksum }
