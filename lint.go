package gopkg

import (
	"sort"
	"strings"
	"strconv"
)

func Lint(
	pkg []FileContents,
	extraLintRules ...func([]FileContents)error,
) error {

	defaultRules := []func([]FileContents)error{
		AddRequiredImports,
		AddAliasToAllImports,
	}

	defaultRules = append(defaultRules, extraLintRules...)

	return LintCustom(pkg, defaultRules...)
}

func LintCustom(
	pkg []FileContents,
	lintRules ...func([]FileContents)error,
) error {

	for _, rule := range lintRules {
		err := rule(pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddAliasToAllImports(pkg []FileContents) error {

	for iF := range pkg {

		existingAliases := make(map[string]bool)

		for iI := range pkg[iF].Imports {
			if pkg[iF].Imports[iI].Alias == "" {
				pathElems := strings.Split(pkg[iF].Imports[iI].Import, "/")

				alias := pathElems[len(pathElems)-1]

				if existingAliases[alias] {
					iAlias := 2
					for {
						potentialAlias := alias + strconv.Itoa(iAlias)
						if existingAliases[potentialAlias] {
							iAlias++
						} else {
							alias = potentialAlias
							break
						}
					}
				}

				existingAliases[alias] = true

				pkg[iF].Imports[iI].Alias = alias
			}
		}
	}
	return nil
}

func AddRequiredImports(pkg []FileContents) error {

	for iF, file := range pkg {
		existingImportSet := make(map[string]bool)
		for _, i := range file.Imports {
			existingImportSet[i.Import] = true
		}

		importsToAdd := complement(
			getFileRequiredTypeImports(pkg[iF]),
			existingImportSet,
		)

		importsToAdd[file.PackageImportPath] = false

		for importPath, ok := range importsToAdd {
			if ok {
				pkg[iF].Imports = append(
					pkg[iF].Imports,
					ImportAndAlias{
						Import: importPath,
					},
				)
			}
		}

		sort.Slice(pkg[iF].Imports, func(i, j int) bool {
			return pkg[iF].Imports[i].Import < pkg[iF].Imports[j].Import
		})
	}

	return nil
}

// GroupStdImportsFirst will move all std imports in all files to their own group
// at the start of the import list.
//
// The ordeing of the std imports within a file will remain unchanged.
// Any existing grouping on std imports within a file will be removed.
// The ordering + grouping of all other imports within a file remains unchanged
func GroupStdImportsFirst(pkg []FileContents) error {

	for iF, f := range pkg {
		stdImports := make([]ImportAndAlias, 0, len(f.Imports))
		otherImports := make([]ImportAndAlias, 0, len(f.Imports))
		for _, i := range f.Imports {
			if isStdImport(i) {
				stdImports = append(stdImports, i)
			} else {
				otherImports = append(otherImports, i)
			}
		}

		lowestOtherImportGroup := int64(1)
		for _, i := range otherImports {
			if i.Group < lowestOtherImportGroup {
				lowestOtherImportGroup = i.Group
			}
		}

		for i := range stdImports {
			stdImports[i].Group = lowestOtherImportGroup - 1
		}

		stdImports = append(stdImports, otherImports...)

		pkg[iF].Imports = stdImports
	}

	return nil
}

func getFileRequiredTypeImports(f FileContents) map[string]bool {

	requiredTypeImports := make(map[string]bool)
	for _, c := range f.Consts {
		requiredTypeImports = union(
			requiredTypeImports,
			c.RequiredImports(),
		)
	}
	for _, v := range f.Vars {
		requiredTypeImports = union(
			requiredTypeImports,
			v.RequiredImports(),
		)
	}
	for _, t := range f.Types {
		requiredTypeImports = union(
			requiredTypeImports,
			t.Type.RequiredImports(),
		)
	}
	for _, f := range f.Functions {
		requiredTypeImports = union(
			requiredTypeImports,
			f.RequiredImports(),
		)
	}
	return requiredTypeImports
}

func isStdImport(i ImportAndAlias) bool {

	// List obtained by running `go list std` for go 1.20.4
	stdImports := map[string]bool{
		"archive/tar": true,
		"archive/zip": true,
		"bufio": true,
		"bytes": true,
		"compress/bzip2": true,
		"compress/flate": true,
		"compress/gzip": true,
		"compress/lzw": true,
		"compress/zlib": true,
		"container/heap": true,
		"container/list": true,
		"container/ring": true,
		"context": true,
		"crypto": true,
		"crypto/aes": true,
		"crypto/cipher": true,
		"crypto/des": true,
		"crypto/dsa": true,
		"crypto/ecdh": true,
		"crypto/ecdsa": true,
		"crypto/ed25519": true,
		"crypto/elliptic": true,
		"crypto/hmac": true,
		"crypto/internal/alias": true,
		"crypto/internal/bigmod": true,
		"crypto/internal/boring": true,
		"crypto/internal/boring/bbig": true,
		"crypto/internal/boring/bcache": true,
		"crypto/internal/boring/sig": true,
		"crypto/internal/edwards25519": true,
		"crypto/internal/edwards25519/field": true,
		"crypto/internal/nistec": true,
		"crypto/internal/nistec/fiat": true,
		"crypto/internal/randutil": true,
		"crypto/md5": true,
		"crypto/rand": true,
		"crypto/rc4": true,
		"crypto/rsa": true,
		"crypto/sha1": true,
		"crypto/sha256": true,
		"crypto/sha512": true,
		"crypto/subtle": true,
		"crypto/tls": true,
		"crypto/x509": true,
		"crypto/x509/internal/macos": true,
		"crypto/x509/pkix": true,
		"database/sql": true,
		"database/sql/driver": true,
		"debug/buildinfo": true,
		"debug/dwarf": true,
		"debug/elf": true,
		"debug/gosym": true,
		"debug/macho": true,
		"debug/pe": true,
		"debug/plan9obj": true,
		"embed": true,
		"embed/internal/embedtest": true,
		"encoding": true,
		"encoding/ascii85": true,
		"encoding/asn1": true,
		"encoding/base32": true,
		"encoding/base64": true,
		"encoding/binary": true,
		"encoding/csv": true,
		"encoding/gob": true,
		"encoding/hex": true,
		"encoding/json": true,
		"encoding/pem": true,
		"encoding/xml": true,
		"errors": true,
		"expvar": true,
		"flag": true,
		"fmt": true,
		"go/ast": true,
		"go/build": true,
		"go/build/constraint": true,
		"go/constant": true,
		"go/doc": true,
		"go/doc/comment": true,
		"go/format": true,
		"go/importer": true,
		"go/internal/gccgoimporter": true,
		"go/internal/gcimporter": true,
		"go/internal/srcimporter": true,
		"go/internal/typeparams": true,
		"go/parser": true,
		"go/printer": true,
		"go/scanner": true,
		"go/token": true,
		"go/types": true,
		"hash": true,
		"hash/adler32": true,
		"hash/crc32": true,
		"hash/crc64": true,
		"hash/fnv": true,
		"hash/maphash": true,
		"html": true,
		"html/template": true,
		"image": true,
		"image/color": true,
		"image/color/palette": true,
		"image/draw": true,
		"image/gif": true,
		"image/internal/imageutil": true,
		"image/jpeg": true,
		"image/png": true,
		"index/suffixarray": true,
		"internal/abi": true,
		"internal/buildcfg": true,
		"internal/bytealg": true,
		"internal/cfg": true,
		"internal/coverage": true,
		"internal/coverage/calloc": true,
		"internal/coverage/cformat": true,
		"internal/coverage/cmerge": true,
		"internal/coverage/decodecounter": true,
		"internal/coverage/decodemeta": true,
		"internal/coverage/encodecounter": true,
		"internal/coverage/encodemeta": true,
		"internal/coverage/pods": true,
		"internal/coverage/rtcov": true,
		"internal/coverage/slicereader": true,
		"internal/coverage/slicewriter": true,
		"internal/coverage/stringtab": true,
		"internal/coverage/test": true,
		"internal/coverage/uleb128": true,
		"internal/cpu": true,
		"internal/dag": true,
		"internal/diff": true,
		"internal/fmtsort": true,
		"internal/fuzz": true,
		"internal/goarch": true,
		"internal/godebug": true,
		"internal/goexperiment": true,
		"internal/goos": true,
		"internal/goroot": true,
		"internal/goversion": true,
		"internal/intern": true,
		"internal/itoa": true,
		"internal/lazyregexp": true,
		"internal/lazytemplate": true,
		"internal/nettrace": true,
		"internal/obscuretestdata": true,
		"internal/oserror": true,
		"internal/pkgbits": true,
		"internal/platform": true,
		"internal/poll": true,
		"internal/profile": true,
		"internal/race": true,
		"internal/reflectlite": true,
		"internal/safefilepath": true,
		"internal/saferio": true,
		"internal/singleflight": true,
		"internal/syscall/execenv": true,
		"internal/syscall/unix": true,
		"internal/sysinfo": true,
		"internal/testenv": true,
		"internal/testlog": true,
		"internal/testpty": true,
		"internal/trace": true,
		"internal/txtar": true,
		"internal/types/errors": true,
		"internal/unsafeheader": true,
		"internal/xcoff": true,
		"io": true,
		"io/fs": true,
		"io/ioutil": true,
		"log": true,
		"log/syslog": true,
		"math": true,
		"math/big": true,
		"math/bits": true,
		"math/cmplx": true,
		"math/rand": true,
		"mime": true,
		"mime/multipart": true,
		"mime/quotedprintable": true,
		"net": true,
		"net/http": true,
		"net/http/cgi": true,
		"net/http/cookiejar": true,
		"net/http/fcgi": true,
		"net/http/httptest": true,
		"net/http/httptrace": true,
		"net/http/httputil": true,
		"net/http/internal": true,
		"net/http/internal/ascii": true,
		"net/http/internal/testcert": true,
		"net/http/pprof": true,
		"net/internal/socktest": true,
		"net/mail": true,
		"net/netip": true,
		"net/rpc": true,
		"net/rpc/jsonrpc": true,
		"net/smtp": true,
		"net/textproto": true,
		"net/url": true,
		"os": true,
		"os/exec": true,
		"os/exec/internal/fdtest": true,
		"os/signal": true,
		"os/user": true,
		"path": true,
		"path/filepath": true,
		"plugin": true,
		"reflect": true,
		"reflect/internal/example1": true,
		"reflect/internal/example2": true,
		"regexp": true,
		"regexp/syntax": true,
		"runtime": true,
		"runtime/cgo": true,
		"runtime/coverage": true,
		"runtime/debug": true,
		"runtime/internal/atomic": true,
		"runtime/internal/math": true,
		"runtime/internal/sys": true,
		"runtime/metrics": true,
		"runtime/pprof": true,
		"runtime/race": true,
		"runtime/trace": true,
		"sort": true,
		"strconv": true,
		"strings": true,
		"sync": true,
		"sync/atomic": true,
		"syscall": true,
		"testing": true,
		"testing/fstest": true,
		"testing/internal/testdeps": true,
		"testing/iotest": true,
		"testing/quick": true,
		"text/scanner": true,
		"text/tabwriter": true,
		"text/template": true,
		"text/template/parse": true,
		"time": true,
		"time/tzdata": true,
		"unicode": true,
		"unicode/utf16": true,
		"unicode/utf8": true,
		"unsafe": true,
		"vendor/golang.org/x/crypto/chacha20": true,
		"vendor/golang.org/x/crypto/chacha20poly1305": true,
		"vendor/golang.org/x/crypto/cryptobyte": true,
		"vendor/golang.org/x/crypto/cryptobyte/asn1": true,
		"vendor/golang.org/x/crypto/hkdf": true,
		"vendor/golang.org/x/crypto/internal/alias": true,
		"vendor/golang.org/x/crypto/internal/poly1305": true,
		"vendor/golang.org/x/net/dns/dnsmessage": true,
		"vendor/golang.org/x/net/http/httpguts": true,
		"vendor/golang.org/x/net/http/httpproxy": true,
		"vendor/golang.org/x/net/http2/hpack": true,
		"vendor/golang.org/x/net/idna": true,
		"vendor/golang.org/x/net/nettest": true,
		"vendor/golang.org/x/net/route": true,
		"vendor/golang.org/x/sys/cpu": true,
		"vendor/golang.org/x/text/secure/bidirule": true,
		"vendor/golang.org/x/text/transform": true,
		"vendor/golang.org/x/text/unicode/bidi": true,
		"vendor/golang.org/x/text/unicode/norm": true,
	}
	return stdImports[i.Import]
}
