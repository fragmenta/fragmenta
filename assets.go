package main

// Perhaps the way to go here is to check when serving -

// DEV - if in dev link to full set of js assets path (separate paths for each)
// when generating script and style tags
// PRODUCTION -  just serve from assets paths (app.js being an amalgam)

const ignore = `
// Handles assets in a web application under /app/assets. In development they are served as individual files to allow quick reloads, in production they can be compiled into the assets folder in public. This does not require manifest files or configuration - assets are bundled by folder and taken in file order. 
package assets

import (
	"bitbucket.org/maxhauser/jsmin" // vendor this dependency?
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Styles and Js are included with these default templates 
const styleTagTemplate = "<link rel='stylesheet' href='/%s' type='text/css' media='all' title='%s' charset='utf-8'>\n\t"
const jsTagTemplate = "<script src='/%s' type='text/javascript' charset='utf-8'></script>\n\t"

// To avoid using os.Glob for every call,
//  we should cache the value of global js and css present in public/assets, and reuse it
// If you don't want to serve compressed assets, remove them from the public/assets dir
type AssetPaths struct {
	js  map[string]string // map of group/path
	css map[string]string // map of group/path
}

var Paths AssetPaths

// Convert a path to a compressed one (at present just adds .gz)
//
func compressedPath(in string) string {
	return in + ".gz"
}

func publicPathForStyleSheetGroup(group string) string {
	// Need slightly more complex hashing function here to generate name
	cssPath := fmt.Sprintf("/assets/styles/%s", group)
	return cssPath
}
func appStyleSheetAssetTags(assetGroups string) string {
	tags := ""

	// Don't compress so write out all global paths	
	cssAssets, _ := filepath.Glob("./app/assets/styles/**/*.css")
	for _, file := range cssAssets {
		group := filepath.Base(filepath.Dir(file))
		tags = tags + fmt.Sprintf(styleTagTemplate, file, group)
	}

	return tags
}

// Expects an asset group like 'global' which corresponds to a folder in app/assets/styles/
func StyleSheetAssetTags(assetGroups string, compiled bool) string {
	tags := ""

	// Get the groups and process each one individually
	groups := strings.Split(assetGroups, ",")
	if compiled {
		// If compressed path exists, use that
		for _, group := range groups {

			// TODO - compile assets into public path so that they are served by web server, not go
			// store asset paths in Paths.css so that we don't use filepath.Glob in production more than once
			// that or we could store the more recent hash for each group on disk somewhere
			// but it seems easier just to check the filesystem as that requires less config

			// Find the gz file in ./public/assets/
			// .gz is not working, css combined is - to investigate how we serve gz files
			// at present this is disabled, add .gz to find gz files instead
			searchPath := fmt.Sprintf("./public/assets/styles/%s*.css", group)
			cssPaths, err := filepath.Glob(searchPath)
			if err == nil && len(cssPaths) > 0 {
				tagPath := strings.Replace(cssPaths[0], "public/", "", 1)
				tags = tags + fmt.Sprintf(styleTagTemplate, tagPath, group)
			} else {
				//	log.Printf("No css found for group %s %s %s",group,searchPath,cssPaths)
				tags = appStyleSheetAssetTags(assetGroups)

			}

		}

	} else {
		tags = appStyleSheetAssetTags(assetGroups)

	}

	return tags
}

func publicPathForJavaScriptGroup(group string) string {
	// Need slightly more complex hashing function here to generate name
	cssPath := fmt.Sprintf("/assets/scripts/%s", group)
	return cssPath
}

func appJavaScriptAssetTags(assetGroups string) string {
	tags := ""

	// Don't compress so write out all global paths	
	jsAssets, _ := filepath.Glob("./app/assets/scripts/**/*.js")
	for _, file := range jsAssets {
		tags = tags + fmt.Sprintf(jsTagTemplate, file)
	}

	return tags
}

// Expects an asset group like 'global' which corresponds to a folder in app/assets/scripts/
func JavaScriptAssetTags(assetGroups string, compiled bool) string {
	tags := ""

	groups := strings.Split(assetGroups, ",")

	if compiled {

		// If compressed path exists, use that
		for _, group := range groups {
			// Find the gz file in ./public/assets/
			searchPath := fmt.Sprintf("./public/assets/scripts/%s*.js", group)
			jsPaths, err := filepath.Glob(searchPath)
			if err == nil && len(jsPaths) > 0 {
				tagPath := strings.Replace(jsPaths[0], "public/", "", 1)
				tags = tags + fmt.Sprintf(jsTagTemplate, tagPath)
			} else {
				//	log.Printf("No css found for group %s %s %s",group,searchPath,jsPaths)
				tags = appJavaScriptAssetTags(assetGroups)
			}

		}

	} else {
		tags = appJavaScriptAssetTags(assetGroups)
	}

	return tags
}

var compress = true

// Compile js and css assets from app/assets/
// and place them in public/assets
// We should only do this on deploy really
// At other times we should use full app/assets paths
// Call this from gopher with gopher generate assets
func Compile() (err error) {

	// For each folder within assets, search for js and css
	// and attempt to compile any assets found within into
	// one big file, and gzip in output)
	appAssets := "./app/assets/"
	jsAssets, _ := filepath.Glob(appAssets + "scripts/**/*.js")
	cssAssets, _ := filepath.Glob(appAssets + "styles/**/*.css")

	// Walk through assets, concatenating them into groups
	var css = map[string]string{}
	for _, file := range cssAssets {
		group := filepath.Base(filepath.Dir(file))

		// Read the file
		in, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal("Error reading file ", file)
		}

		// Append to group
		css[group] = css[group] + "\n" + string(in)

	}

	var js = map[string]string{}
	for _, file := range jsAssets {
		group := filepath.Base(filepath.Dir(file))

		if compress {
			out := new(bytes.Buffer)
			in, err := os.Open(file) // For read access.
			if err != nil {
				log.Fatal(err)
			} else {
				jsmin.Run(in, out)
			}
			js[group] = js[group] + "\n" + out.String()
		} else {
			// Read the file
			in, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal("Error reading file ", file)
			}

			js[group] = js[group] + "\n" + string(in)
		}

	}

	for group, value := range js {
		// Write out the concatenated group file for js
		md5 := md5.New()
		io.WriteString(md5, fmt.Sprintf("%s", value)) // Peek simpler?
		hash := fmt.Sprintf("%x", md5.Sum(nil))

		file := fmt.Sprintf("./public%s-%s.js", publicPathForJavaScriptGroup(group), hash)
		os.MkdirAll(filepath.Dir(file), 0774)
		err := ioutil.WriteFile(file, []byte(value), 0774)
		if err != nil {
			log.Fatal("Error writing file ", file)
		}
		// Now compress
		if compress {
			compressFile(file)
		}
	}

	for group, value := range css {
		// Write out the concatenated group file for css
		md5 := md5.New()
		io.WriteString(md5, fmt.Sprintf("%s", value)) // Peek simpler?
		hash := fmt.Sprintf("%x", md5.Sum(nil))
		file := fmt.Sprintf("./public%s-%s.css", publicPathForStyleSheetGroup(group), hash)
		os.MkdirAll(filepath.Dir(file), 0774)
		err := ioutil.WriteFile(file, []byte(value), 0774)
		if err != nil {
			log.Fatal("Error writing file ", file)
		}

		// Grab md5 rename and compress
		if compress {
			compressFile(file)
		}
	}

	return err

}

// Compress a file with gz, in preparation for serving with another server
func compressFile(file string) error {

	contents, err := os.Open(file)
	if err != nil {
		return err
	}

	newName := strings.Replace(file, ".js", ".js.gz", 1)
	newName = strings.Replace(newName, ".css", ".css.gz", 1)

	fileOS, err := os.Create(newName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open %s: error: %s\n", newName, err)
		os.Exit(1)
	}
	defer fileOS.Close()
	fileGzip := gzip.NewWriter(fileOS)
	if err != nil {
		fmt.Printf("The file %v is not in gzip format.\n", newName)
		return err
	}
	defer fileGzip.Close()

	contents.Seek(0, 0)
	io.Copy(fileGzip, contents)

	return err

}
`
