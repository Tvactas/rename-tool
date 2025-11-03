package antisamename

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"

	"rename-tool/common/dialogcustomize"
	"rename-tool/common/pathgen"
	"rename-tool/setting/model"
)

// CheckConflicts computes target paths for the given files and config,
// and returns a list of conflicting target paths (case-insensitive):
//  1. duplicates within the batch; 2) paths that already exist on disk
//     and are not the same file (ignoring case-only change).
func CheckConflicts(files []string, config model.RenameConfig) ([]string, error) {
	seen := make(map[string]string) // lower(target) -> firstOriginal
	conflictsSet := make(map[string]struct{})

	// state used to mirror batch naming logic
	globalCounter := 0
	perExtCounters := make(map[string]int)

	addConflict := func(path string) {
		conflictsSet[path] = struct{}{}
	}

	for _, file := range files {
		var target string
		var err error

		switch config.Type {
		case model.RenameTypeBatch:
			current := globalCounter
			globalCounter++
			ext := filepath.Ext(file)
			if _, ok := perExtCounters[ext]; !ok {
				perExtCounters[ext] = 0
			}
			target, err = pathgen.GenerateBatchRenamePath(file, config, current, perExtCounters)
		case model.RenameTypeExtension:
			target, err = pathgen.GenerateExtensionRenamePath(file, config)
		case model.RenameTypeCase:
			target, err = pathgen.GenerateCaseRenamePath(file, config)
		case model.RenameTypeInsertChar:
			target, err = pathgen.GenerateInsertCharRenamePath(file, config)
		case model.RenameTypeReplace:
			target, err = pathgen.GenerateReplaceRenamePath(file, config)
		case model.RenameTypeDeleteChar:
			target, err = pathgen.GenerateDeleteCharRenamePath(file, config)
		default:
			// unknown type: skip
			continue
		}
		if err != nil {
			// treat as conflict source information
			addConflict(file)
			continue
		}

		lower := strings.ToLower(target)
		if first, exists := seen[lower]; exists {
			addConflict(first)
			addConflict(target)
		} else {
			seen[lower] = file
		}

		// filesystem existence (ignore case-only self-change)
		if info, err := os.Stat(target); err == nil && info != nil {
			if !strings.EqualFold(file, target) {
				addConflict(target)
			}
		}
	}

	// collect set to slice
	out := make([]string, 0, len(conflictsSet))
	for p := range conflictsSet {
		out = append(out, p)
	}
	return out, nil
}

// CheckAndShowConflicts runs CheckConflicts and shows a dialog if any conflicts found.
// Returns true if a dialog was shown (caller should abort execution).
func CheckAndShowConflicts(window fyne.Window, files []string, config model.RenameConfig) (bool, error) {
	conflicts, err := CheckConflicts(files, config)
	if err != nil {
		return false, err
	}
	if len(conflicts) > 0 {
		dialogcustomize.ShowMultiLineCopyDialog("error", dialogTr("duplicateNames"), conflicts, window)
		return true, nil
	}
	return false, nil
}

// GenerateUniquePath returns a non-conflicting file path by appending
// an incremental suffix like _1, _2 before the extension when needed.
func GenerateUniquePath(desiredPath string) string {
	base := desiredPath
	counter := 1
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	path := desiredPath
	for {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path
		}
		path = fmt.Sprintf("%s_%d%s", name, counter, ext)
		counter++
	}
}
