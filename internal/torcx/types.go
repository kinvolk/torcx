// Copyright 2017 CoreOS Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package torcx

import (
	"encoding/json"
	"fmt"
)

const (
	// SealUpperProfile is the key label for user profile name
	SealUpperProfile = "TORCX_UPPER_PROFILE"
	// SealLowerProfiles is the key label for vendor profile path
	SealLowerProfiles = "TORCX_LOWER_PROFILES"
	// SealRunProfilePath is the key label for vendor profile path
	SealRunProfilePath = "TORCX_PROFILE_PATH"
	// SealBindir is the key label for seal bindir
	SealBindir = "TORCX_BINDIR"
	// SealUnpackdir is the key label for seal unpackdir
	SealUnpackdir = "TORCX_UNPACKDIR"
	// ImageManifestV0K - image manifest kind, v0
	ImageManifestV0K = "image-manifest-v0"
	// CommonConfigV0K - common torcx config kind, v0
	CommonConfigV0K = "torcx-config-v0"
)

// ConfigV0 holds common torcx configuration in JSON format
type ConfigV0 struct {
	Kind  string       `json:"kind"`
	Value CommonConfig `json:"value"`
}

// CommonConfig contains runtime configuration items common to all
// torcx subcommands
type CommonConfig struct {
	BaseDir    string   `json:"base_dir,omitempty"`
	RunDir     string   `json:"run_dir,omitempty"`
	UsrDir     string   `json:"usr_dir,omitempty"`
	ConfDir    string   `json:"conf_dir,omitempty"`
	StorePaths []string `json:"store_paths,omitempty"`
}

// ApplyConfig contains runtime configuration items specific to
// the `apply` subcommand
type ApplyConfig struct {
	CommonConfig
	LowerProfiles []string
	UpperProfile  string
}

// ProfileConfig contains runtime configuration items specific to
// the `profile` subcommand
type ProfileConfig struct {
	CommonConfig
	LowerProfileNames  []string
	UserProfileName    string
	CurrentProfilePath string
	NextProfile        string
}

// Archive represents a .torcx.squashfs or .torcx.tgz on disk
type Archive struct {
	Image
	Filepath string        `json:"filepath"`
	Format   ArchiveFormat `json:"format"`
}

// Image represents an addon archive within a profile.
type Image struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Remote    string `json:"remote"`
}

// ArchiveFormat is a torcx archive format, either 'tgz' or 'squashfs'
type ArchiveFormat string

const (
	// ArchiveFormatUnknown is the zero value of ArchiveFormat. It indicates the image format is unknown
	ArchiveFormatUnknown ArchiveFormat = ""
	// ArchiveFormatTgz indicates a tar-gzipped image
	ArchiveFormatTgz = "tgz"
	// ArchiveFormatSquashfs indicates a squashfs image archive
	ArchiveFormatSquashfs = "squashfs"
)

// UnmarshalJSON unmarshals an ArchiveFormat
func (arf *ArchiveFormat) UnmarshalJSON(b []byte) error {
	s := ""
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case ArchiveFormatTgz:
		*arf = ArchiveFormatTgz
	case ArchiveFormatSquashfs:
		*arf = ArchiveFormatSquashfs
	default:
		return fmt.Errorf("could not unmarshal into ArchiveFormat: must be one of %q, %q", ArchiveFormatTgz, ArchiveFormatSquashfs)
	}
	return nil
}

// UnmarshalJSON unmarshals an Archive, including defaulting the "format" field
// to "tgz" if it was not set.
func (ar *Archive) UnmarshalJSON(b []byte) error {
	type archiveAlias Archive
	var arAlias archiveAlias
	if err := json.Unmarshal(b, &arAlias); err != nil {
		return err
	}
	if arAlias.Format == ArchiveFormatUnknown {
		// Default to tgz if it wasn't set
		arAlias.Format = ArchiveFormatTgz
	}
	*ar = Archive(arAlias)
	return nil
}

// FileSuffix returns the file extension this archive format must have.
func (arf ArchiveFormat) FileSuffix() string {
	return fmt.Sprintf(".torcx.%s", arf)
}

// ToJSONV0 converts an internal Image into ImageV0.
func (im Image) ToJSONV0() ImageV0 {
	return ImageV0{
		Name:      im.Name,
		Reference: im.Reference,
	}
}

// ImageFromJSONV0 converts an ImageV0 into an internal Image.
func ImageFromJSONV0(j ImageV0) Image {
	return Image{
		Name:      j.Name,
		Reference: j.Reference,
		Remote:    "",
	}
}

// ToJSONV1 converts an internal Image into ImageV1.
func (im Image) ToJSONV1() ImageV1 {
	return ImageV1{
		Name:      im.Name,
		Reference: im.Reference,
		Remote:    "",
	}
}

// ImageFromJSONV1 converts an ImageV1 into an internal Image.
func ImageFromJSONV1(j ImageV1) Image {
	entry := Image{
		Name:      j.Name,
		Reference: j.Reference,
		Remote:    j.Remote,
	}
	return entry
}

// ImagesToJSONV0 converts an internal Image list into ImagesV0.
func ImagesToJSONV0(ims []Image) ImagesV0 {
	j := ImagesV0{}
	for _, im := range ims {
		entry := im.ToJSONV0()
		j.Images = append(j.Images, entry)
	}
	return j
}

// ImagesFromJSONV0 converts an ImagesV0 into an internal Image list.
func ImagesFromJSONV0(j ImagesV0) []Image {
	result := []Image{}
	for _, im := range j.Images {
		entry := ImageFromJSONV0(im)
		result = append(result, entry)
	}
	return result
}

// ImagesToJSONV1 converts an internal Image list into ImagesV1.
func ImagesToJSONV1(ims []Image) ImagesV1 {
	j := ImagesV1{}
	for _, im := range ims {
		entry := im.ToJSONV1()
		j.Images = append(j.Images, entry)
	}
	return j
}

// ImagesFromJSONV1 converts an ImagesV1 into an internal Image list.
func ImagesFromJSONV1(j ImagesV1) []Image {
	result := []Image{}
	for _, im := range j.Images {
		entry := ImageFromJSONV1(im)
		result = append(result, entry)
	}
	return result
}

// ImageManifestV0 holds JSON image manifest
type ImageManifestV0 struct {
	Kind  string `json:"kind"`
	Value Assets `json:"value"`
}

// Assets holds lists of assets propagated from an image to the system
type Assets struct {
	Binaries  []string `json:"bin,omitempty"`
	Network   []string `json:"network,omitempty"`
	Units     []string `json:"units,omitempty"`
	Sysusers  []string `json:"sysusers,omitempty"`
	Tmpfiles  []string `json:"tmpfiles,omitempty"`
	UdevRules []string `json:"udev_rules,omitempty"`
}

type Remote struct {
	TemplateURL string
	ArmoredKeys []string
}

// RemoteFromJSONV0 translates a RemoteKeyV0 to an internal Remote.
func RemoteFromJSONV0(j RemoteV0) Remote {
	res := Remote{
		TemplateURL: j.BaseURL,
	}
	for _, key := range j.Keys {
		res.ArmoredKeys = append(res.ArmoredKeys, key.ArmoredKeyring)
	}
	return res
}

// RemoteContents holds contents metadata for a remote manifest.
type RemoteContents struct {
	Images map[string]RemoteImage
}

// RemoteContentsFromJSONV1 translates a RemoteImagesV1 to an internal Remote.
func RemoteContentsFromJSONV1(j RemoteImagesV1) RemoteContents {
	var res RemoteContents

	images := map[string]RemoteImage{}
	for _, im := range j.Images {
		if im.Name == "" {
			continue
		}
		tmpVersions := []RemoteVersion{}
		for _, v := range im.Versions {
			tmpVersions = append(tmpVersions, RemoteVersionFromJSONV1(v))
		}
		tmpImage := RemoteImage{
			name:           im.Name,
			defaultVersion: im.DefaultVersion,
			versions:       tmpVersions,
		}
		images[im.Name] = tmpImage
	}
	res.Images = images

	return res
}

// RemoteImage list remote versions of an image.
type RemoteImage struct {
	defaultVersion string
	name           string
	versions       []RemoteVersion
}

// RemoteVersion describes a remote image archive.
type RemoteVersion struct {
	format   string
	version  string
	hash     string
	location string
}

// RemoteVersionFromJSONV1 translates a RemoteVersionV1 to an internal RemoteVersion.
func RemoteVersionFromJSONV1(j RemoteVersionV1) RemoteVersion {
	remoteVer := RemoteVersion{
		format:   j.Format,
		hash:     j.Hash,
		location: j.Location,
		version:  j.Version,
	}
	return remoteVer
}

// kindValueJSON holds a generic, typed, kind-value JSON manifest.
type kindValueJSON struct {
	Kind  string          `json:"kind"`
	Value json.RawMessage `json:"value"`
}
