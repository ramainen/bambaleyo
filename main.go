package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/manifoldco/promptui"
)

//bam
//bam artklen.loc.artklen.ru
//bam artklen.loc.artklen.ru ubuntu64

func main() {

	searchTemplatePaths := []string{"/etc/bambaleyo/templates", "/mnt/c/goprojects/bamtemplates", "/var/www/templates"}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	searchTemplatePaths = append(searchTemplatePaths, filepath.Dir(ex))

	fmt.Println(`
     ____  ___    __  _______  ___    __    ________  ______  __
    / __ )/   |  /  |/  / __ )/   |  / /   / ____/\ \/ / __ \/ /
   / __  / /| | / /|_/ / __  / /| | / /   / __/    \  / / / / /
  / /_/ / ___ |/ /  / / /_/ / ___ |/ /___/ /___    / / /_/ /_/
 /_____/_/  |_/_/  /_/_____/_/  |_/_____/_____/   /_/\____(_)
 `)

	template := "default"
	domain := "default.loc"

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		//

		prompt := promptui.Prompt{
			Label: "Введи домен",
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		domain = result
		fmt.Printf("Your username is %q\n", result)

	} else {
		domain = os.Args[1]
	}
	if len(argsWithoutProg) == 1 || len(argsWithoutProg) == 0 {
		//

		itemstoask := []string{}
		for _, element := range searchTemplatePaths {
			files, err := ioutil.ReadDir(element)
			if err != nil {
				//	log.Fatal(err)
				continue
			}

			for _, f := range files {
				if f.IsDir() {
					//fmt.Println()
					itemstoask = append(itemstoask, element+"/"+f.Name())
				}
			}
		}

		prompt := promptui.Select{
			Label: "Выбери шаблон",
			Items: itemstoask,
		}
		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		template = result

	} else {
		template = os.Args[2]
	}

	fmt.Println("Домен: ", domain)
	fmt.Println("Шаблон: ", template)

	targetPath := "/var/www/domains/" + domain
	fmt.Println("Копирую файлы из ", template, " в "+targetPath)

	cmd := exec.Command("mkdir", targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Не удалось создать папку " + targetPath + ". Она уже есть?")
	}

	cmd = exec.Command("mkdir", targetPath+"/www")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("mkdir", targetPath+"/tmp")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("mkdir", targetPath+"/logs")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("chmod", "-R", "777", targetPath+"/logs")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Не удалось chmod на logs")
	}

	cmd = exec.Command("chmod", "-R", "777", targetPath+"/tmp")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Не удалось chmod на tmp")
	}

	files, err := ioutil.ReadDir(template)
	if err != nil {
		//	log.Fatal(err)
		fmt.Println("Не удалось прочитать файлы из директории " + template)
	}

	for _, f := range files {
		if !f.IsDir() {
			//fmt.Println()
			fileContents, err := ioutil.ReadFile(template + "/" + f.Name())
			if err != nil {
				//	log.Fatal(err)
				fmt.Println("Не удалось прочитать файл " + template + "/" + f.Name())
				continue
			}
			fileContents = bytes.ReplaceAll(fileContents, []byte("###DOMAIN###"), []byte(domain))
			err = ioutil.WriteFile(targetPath+"/"+f.Name(), fileContents, 0777)
			if err != nil {
				//	log.Fatal(err)
				fmt.Println("Не удалось записать файл " + targetPath + "/" + f.Name())
				continue
			}
		}
	}
	fmt.Println(">docker container restart core_nginx_1")
	cmd = exec.Command("docker", "container", "restart", "core_nginx_1")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Не удалось запустить докер nginx")
	}
	fmt.Println(">docker-compose up -d")
	cmd = exec.Command("docker-compose", "up", "-d")
	cmd.Dir = targetPath
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Не удалось запустить докер проекта")
	}

}
