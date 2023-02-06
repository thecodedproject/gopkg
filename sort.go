package gopkg

import (
	"sort"
)

// TODO Implement general sorting func of whole Contents type

func SortFuncs(f []DeclFunc) {

	sort.Slice(f, func(i, j int) bool {

		iF := f[i]
		jF := f[j]

		if iF.Receiver.TypeName != jF.Receiver.TypeName {

			if iF.Receiver.TypeName == "" {
				return false
			}
			if jF.Receiver.TypeName == "" {
				return true
			}

			return iF.Receiver.TypeName < jF.Receiver.TypeName
		}

		return iF.Name < jF.Name
	})

	return

}
