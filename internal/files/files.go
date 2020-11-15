package files

import (
	"github.com/goreleaser/fileglob"
	"sort"

	"github.com/goreleaser/nfpm/internal/glob"
)

// FileToCopy describes the source and destination
// of one file to copy into a package.
type FileToCopy struct {
	Source string `yaml:"source,omitempty"`
	Destination string `yaml:"destination,omitempty"`
	Type string `yaml:"type,omitempty"`
	Packager string `yaml:"packager,omitempty`
	Mode uint `yaml:"mode,omitempty"`
}

//type tmpFruitBasket []map[string]yaml.Node
//func (f *FileToCopy) UnmarshalYAML(value *yaml.Node) error {
//	fmt.Println(value.Kind)
//	if value.Kind == yaml.MappingNode {
//		return nil
//	}
//	if err := value.Decode(f); err != nil {
//		return err
//	}
//	return nil
//}

// Expand gathers all of the real files to be copied into the package.
func Expand(filesSrcDstMap map[string]string, disableGlobbing bool) (files []*FileToCopy, err error) {
	var globbed map[string]string

	for srcGlob, dstRoot := range filesSrcDstMap {
		if disableGlobbing {
			srcGlob = fileglob.QuoteMeta(srcGlob)
		}

		globbed, err = glob.Glob(srcGlob, dstRoot)
		if err != nil {
			return nil, err
		}
		appendAndSort(globbed, files)
	}

	return files, nil
}
// Expand gathers all of the real files to be copied into the package.
func ExpandFiles(filesSrcDstMap []*FileToCopy, disableGlobbing bool) (files []*FileToCopy, err error) {
	var globbed map[string]string
	for _, f := range filesSrcDstMap {
		if disableGlobbing {
			f.Source = fileglob.QuoteMeta(f.Source)
		}

		globbed, err = glob.Glob(f.Source, f.Destination)
		if err != nil {
			return nil, err
		}
		appendAndSort(globbed, files)
	}

	return files, nil
}

func appendAndSort(globbed map[string]string, files []*FileToCopy) {
	for src, dst := range globbed {
		files = append(files, &FileToCopy{
			Destination: dst,
			Source: src,
		})

	}
	// sort the files for reproducibility and general cleanliness
	sort.Slice(files, func(i, j int) bool {
		a, b := files[i], files[j]
		if a.Source != b.Source {
			return a.Source < b.Source
		}
		return a.Destination < b.Destination
	})
}