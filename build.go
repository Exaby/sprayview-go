package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func checkFileExists(files []string) (missingFiles []string) {
	fmt.Println("Checking for required files...");
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
		fmt.Println("success: ", file);
	}
	return missingFiles
}

func checkDirectoryExists(directory string) error {
	fmt.Println("Checking for required directory...");
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return fmt.Errorf("%s folder doesn't exist", directory)
	}
	fmt.Println("success: ", directory);
	return nil
}

func convertVTFtoJPG(srcDirectory, dstDirectory string, vtfCmdPath string) error {
	command := fmt.Sprintf("%s -folder \"%s*.vtf\" -output \"%s\" -exportformat jpg", vtfCmdPath, srcDirectory, dstDirectory)
	fmt.Println("Running command:", command)

	cmd := exec.Command(vtfCmdPath, "-folder", fmt.Sprintf("%s*.vtf", srcDirectory), "-output", dstDirectory, "-exportformat", "jpg")
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("error converting VTF files: %v", err)
	}
	// wait for command to finish before returning
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error converting VTF files: %v", err)
	}
	fmt.Println("Conversion success");
	return nil
}


func processSprayFiles(srcDirectory, dstDirectory string) (map[string][]string, error) {
	fmt.Println("Processing spray files into JSON...");
	spraysBySteamID := make(map[string][]string)

	err := os.MkdirAll(dstDirectory, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating spray files directory: %v", err)
	}
	d, err := os.Open(dstDirectory)
	if err != nil {
		return nil, fmt.Errorf("error reading spray files: %v", err)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("error reading spray files: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		parts := strings.Split(file.Name(), "_")
		if len(parts) != 2 {
			continue
		}
		steamID, sprayID := parts[0], parts[1]
		sprays, ok := spraysBySteamID[steamID]
		if !ok {
			sprays = []string{}
		}
		sprays = append(sprays, sprayID)
		spraysBySteamID[steamID] = sprays
	}
	fmt.Println("JSON success");
	return spraysBySteamID, nil
}

func writeSpraysToJSON(sprays map[string][]string, dstFile string) error {
	fmt.Println("Writing spray files to JSON...");
	f, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	err = json.NewEncoder(w).Encode(sprays)
	if err != nil {
		return fmt.Errorf("error encoding sprays to JSON: %v", err)
	}
	w.Flush()
	fmt.Println("JSON write success");
	return nil
}

func main() {
	requiredFiles := []string{"VTFLib.dll", "VTFCmd.exe", "HLLib.dll", "DevIL.dll"}
	missingFiles := checkFileExists(requiredFiles)
	if len(missingFiles) > 0 {
		fmt.Println("The following required files are missing:")
		for _, file := range missingFiles {
			fmt.Println("-", file)
		}
		return
	}

	srcDirectory := "sr_sprays"
	err := checkDirectoryExists(srcDirectory)
	if err != nil {
		fmt.Println(err)
		return
	}
	currentDirectory, _ := filepath.Abs(".")
	srcPath, _ := filepath.Abs(srcDirectory)
	dstDirectory := filepath.Dir(srcPath) + "\\sr_sprays_jpg"
	dstPath, _ := filepath.Abs(dstDirectory)
	vtfCmdPath, _ := filepath.Abs("VTFCmd.exe")

	if srcInfo, err := os.Stat(srcPath); err == nil {
		if dstInfo, err := os.Stat(dstPath); err == nil {
			if srcInfo.Size() > dstInfo.Size() {
				os.RemoveAll(dstPath)
			}
		}
	}
	fmt.Println();
	err = convertVTFtoJPG(srcPath+"\\", dstPath+"\\", vtfCmdPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println();
	sprays, err := processSprayFiles(srcPath, dstPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println();
	err = writeSpraysToJSON(sprays, currentDirectory+"\\sr_sprays_api.json")
	if err != nil {
		fmt.Println(err)
		return
	}
}