// Copyright 2020 xiashuangxi<xiashuangxi@hotmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mingrammer/cfmt"
)

const (
	APP_VERSION = "0.0.1"

	EXT              = ".ada"
	SOURCE_DIRECTORY = "./"
	TEMPLATE_FILE    = "./LICENSETEMPLATE"
)

var (
	source_directory string
	ext              string
	template         string
)

///////////////////////// utils /////////////////////////
func is_empty(str string) bool { return str == "" }

func is_exists(path string) bool {
	if !is_empty(path) {
		_, err := os.Stat(path)
		if err == nil || os.IsExist(err) {
			return true
		}
	}
	return false
}

func is_directory(path string) bool {
	if !is_empty(path) {
		if is_exists(path) {
			fi, err := os.Stat(path)
			if err != nil {
				return false
			}
			return fi.IsDir()
		}
	}
	return false
}

///////////////////// end utils /////////////////////////

func file_list(ext string, directory string) (files []string, err error) {
	if is_empty(ext) {
		log.Fatal("err")
		return nil, err
	}
	if is_empty(directory) {
		log.Fatal("err")
		return nil, err
	}

	if is_directory(directory) {
		files, _ := ioutil.ReadDir(directory)
		var files_string []string
		var sp string = "/"
		if string(directory[len(directory)-1]) == "/" {
			sp = ""
		}

		for _, f := range files {
			if strings.EqualFold(ext, path.Ext(f.Name())) {
				files_string = append(files_string, directory+sp+f.Name())
			}
		}
		return files_string, nil
	}
	return nil, err
}

func license_exit_sp(data []byte) int {
	offset := 0
	ss := strings.Split(string(data), "\n")

	for _, s := range ss {
		offset += len([]byte(s)) + 1
		if len(strings.TrimSpace(s)) == 0 {
			break
		}
	}
	return offset
}

func execute(ext string, source_directory string, template string) (count int, err error) {
	if is_empty(ext) {
		ext = EXT
	}
	if is_empty(source_directory) {
		source_directory = SOURCE_DIRECTORY
	}
	if is_empty(template) {
		template = TEMPLATE_FILE
	}

	if !is_exists(template) {
		log.Fatal("不存在模板文件。")
		return 0, nil
	}
	var result_exec_count int

	license, err := ioutil.ReadFile(template)
	if err != nil {
		log.Fatal(err)
	} else {

		fs, err := file_list(ext, source_directory)
		if err != nil {
			log.Fatal(err)
		}

		for _, fp := range fs {
			source_data, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintf(os.Stdout, "%s \t执行\r", fp)

			sp := []byte{'\n', '\n'}
			sp_offset := license_exit_sp(source_data)
			// println(sp_offset)
			var new_d []byte
			if len(source_data) == 0 {
				new_d = make([]byte, 0)
			} else {
				new_d = source_data[sp_offset:len(source_data)]
			}
			source := make([]byte, len(license)+len(new_d)+len(sp))

			copy(source, license)
			copy(source[len(license):], sp)
			copy(source[len(license)+len(sp):], new_d)
			ioutil.WriteFile(fp, source, os.ModePerm)
			result_exec_count += 1
			time.Sleep(time.Microsecond * 30)
			fmt.Fprintf(os.Stdout, "%s \t", fp)
			cfmt.Success("完成\r\n")
		}

	}
	return result_exec_count, nil
}

func init() {
	flag.StringVar(&ext, "e", EXT, "源文件扩展名")
	flag.StringVar(&source_directory, "s", SOURCE_DIRECTORY, "源代码文件目录")
	flag.StringVar(&template, "t", TEMPLATE_FILE, "License 模板文件")
}

func main() {

	cfmt.Success("LicenseTemp:", "v"+APP_VERSION)
	fmt.Println(`
 _    _                    _____               
| |  (_)__ ___ _ _  ___ __|_   _|__ _ __  _ __ 
| |__| / _/ -_) ' \(_-</ -_)| |/ -_) '  \| '_ \
|____|_\__\___|_||_/__/\___||_|\___|_|_|_| .__/
                                    	 |_|    
	`)

	flag.Parse()
	count, _ := execute(ext, source_directory, template)
	fmt.Printf("执行完成，实现对 %d 个文件的 License 信息的插入。", count)
}
