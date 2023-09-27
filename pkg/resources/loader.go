package resources

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kanengo/egoist/utils"
	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/validation/path"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type kubernetesManifest interface {
	Kind() string
}

type Loader[T kubernetesManifest] interface {
	Load() ([]T, error)
}

type LocalResourcesLoader[T kubernetesManifest] struct {
	paths []string
	kind  string
}

const yamlSeparator = "\n---"

func NewLocalResourcesLoader[T kubernetesManifest](paths ...string) *LocalResourcesLoader[T] {
	var zero T
	return &LocalResourcesLoader[T]{paths: paths, kind: zero.Kind()}
}

func (lr *LocalResourcesLoader[T]) Load() ([]T, error) {
	var resources []T
	log.Debug("LocalResourcesLoader", zap.Strings("paths", lr.paths))
	for _, p := range lr.paths {
		loaded, err := lr.loadFromPath(p)
		if err != nil {
			return nil, err
		}
		if len(loaded) > 0 {
			resources = append(resources, loaded...)
		}
	}

	return resources, nil
}

func (lr *LocalResourcesLoader[T]) loadFromPath(path string) ([]T, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	comps := make([]T, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		log.Debug("LocalResourcesLoader file", zap.String("name", fileName))
		if !utils.IsYaml(fileName) {
			log.Warn("A non-YAML file was detected, it will not be loaded", zap.String("fileName", fileName),
				zap.String("kind", lr.kind))
			continue
		}

		comps = append(comps, lr.loadFromFile(filepath.Join(path, fileName))...)
	}

	return comps, nil
}

func (lr *LocalResourcesLoader[T]) loadFromFile(path string) []T {
	var comps []T

	b, err := os.ReadFile(path)
	if err != nil {
		log.Warn("egoist load component error when reading file", zap.String("filePath", path), zap.String("kind", lr.kind))
		return comps
	}
	comps, errs := lr.decodeYaml(b)
	for _, err := range errs {
		log.Warn("egoist load error when parsing yaml resource", zap.Error(err), zap.String("filePath", path),
			zap.String("kind", lr.kind))
	}
	return comps
}

type typeInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// decodeYaml decodes the yaml document.
func (lr *LocalResourcesLoader[T]) decodeYaml(b []byte) ([]T, []error) {
	list := make([]T, 0)
	var errors []error
	scanner := bufio.NewScanner(bytes.NewReader(b))
	scanner.Split(splitYamlDoc)

	for {
		if !scanner.Scan() {
			err := scanner.Err()
			if err != nil {
				errors = append(errors, err)
				continue
			}

			break
		}

		scannerBytes := scanner.Bytes()
		var ti typeInfo
		if err := yaml.Unmarshal(scannerBytes, &ti); err != nil {
			errors = append(errors, err)
			continue
		}

		if ti.Kind != lr.kind {
			continue
		}

		if errs := path.IsValidPathSegmentName(ti.Name); len(errs) > 0 {
			errors = append(errors, fmt.Errorf("invalid name %q for %q: %s", ti.Name, lr.kind, strings.Join(errs, "; ")))
			continue
		}

		var manifest T
		//if lr.zvFn != nil {
		//	manifest = lr.zvFn()
		//}
		if err := yaml.Unmarshal(scannerBytes, &manifest); err != nil {
			errors = append(errors, err)
			continue
		}
		list = append(list, manifest)
	}

	return list, errors
}

// splitYamlDoc - splits the yaml docs.
func splitYamlDoc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		i += sep
		after := data[i:]

		if len(after) == 0 {
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
