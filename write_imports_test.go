package gopkg_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func TestWriteImports(t *testing.T) {

	testCases := []struct {
		Name        string
		Imports     []gopkg.ImportAndAlias
		ExpectedErr error
	}{
		{
			Name: "empty imports writes nothing",
		},
		{
			Name:    "single import without alias",
			Imports: tmpl.UnnamedImports("my/import"),
		},
		{
			Name: "single import with alias",
			Imports: []gopkg.ImportAndAlias{
				{
					Import: "some/import",
					Alias:  "my_alias",
				},
			},
		},
		{
			Name: "multiple imports without alias or group",
			Imports: tmpl.UnnamedImports(
				"some",
				"my/import",
				"another/import",
			),
		},
		{
			Name: "multiple imports with alias and no group",
			Imports: []gopkg.ImportAndAlias{
				{
					Import: "aimport",
					Alias:  "a_alias",
				},
				{
					Import: "some/import",
					Alias:  "some_alias",
				},
				{
					Import: "another/import",
					Alias:  "another_alias",
				},
			},
		},
		{
			Name: "multiple imports with mix of alias and no group",
			Imports: []gopkg.ImportAndAlias{
				{
					Import: "animport",
				},
				{
					Import: "some/path/to/import",
					Alias:  "some_alias",
				},
				{
					Import: "another/import",
					Alias:  "another_alias",
				},
				{
					Import: "pkg",
				},
			},
		},
		{
			Name: "multiple imports in several groups",
			Imports: []gopkg.ImportAndAlias{
				{
					Import: "group0/import1",
					Group:  0,
				},
				{
					Import: "group0/import2",
					Group:  0,
				},
				{
					Import: "group1/import1",
					Group:  1,
				},
				{
					Import: "group1/import2",
					Group:  1,
				},
				{
					Import: "group2/import1",
					Group:  2,
				},
				{
					Import: "group2/import2",
					Group:  2,
				},
			},
		},
		{
			Name: "several groups starting from 5",
			Imports: []gopkg.ImportAndAlias{
				{
					Import: "group5/import1",
					Group:  5,
				},
				{
					Import: "group5/import2",
					Group:  5,
				},
				{
					Import: "group6/import1",
					Group:  6,
				},
				{
					Import: "group6/import2",
					Group:  6,
				},
				{
					Import: "group7/import1",
					Group:  7,
				},
				{
					Import: "group7/import2",
					Group:  7,
				},
			},
		},
		{
			Name: "several groups not monotonically increasing but in order",
			Imports: []gopkg.ImportAndAlias{
				{
					Import: "group-1/import1",
					Group:  -1,
				},
				{
					Import: "group-1/import2",
					Group:  -1,
				},
				{
					Import: "group6/import1",
					Group:  6,
				},
				{
					Import: "group6/import2",
					Group:  6,
				},
				{
					Import: "group10/import1",
					Group:  10,
				},
				{
					Import: "group10/import2",
					Group:  10,
				},
			},
		},
		{
			Name: "several groups which are not in order returns error",
			Imports: []gopkg.ImportAndAlias{
				{
					Import: "group2/import1",
					Group:  2,
				},
				{
					Import: "group10/import1",
					Group:  10,
				},
				{
					Import: "group6/import1",
					Group:  6,
				},
				{
					Import: "group10/import2",
					Group:  10,
				},
				{
					Import: "group6/import2",
					Group:  6,
				},
				{
					Import: "group2/import2",
					Group:  2,
				},
			},
			ExpectedErr: errors.New("WriteImports: import groups are not in order"),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			buffer := bytes.NewBuffer(nil)

			err := gopkg.WriteImports(
				buffer,
				test.Imports,
			)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), buffer.Bytes())
		})
	}
}
