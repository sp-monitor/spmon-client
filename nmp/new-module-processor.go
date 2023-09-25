package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
)

func main() {
	logger, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_DEBUG, "SPMON-CLIENT")
	if err != nil {
		log.Fatal(err)
	}
	//look for files in the module directory that end in .tar.gz
	filepath.WalkDir("./modules/", func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			logger.Err(fmt.Sprintf("Error Reading Modules Directory: %+v", err))
		}

		//If the file is not a directory
		if !file.IsDir() {
			//Create new directory with the file name minus the .tar.gz extension
			os.Mkdir(fmt.Sprintf("./modules/%s", file.Name()[:len(file.Name())-7]), 0755)

			//Move the .tar.gz file into the new directory
			os.Rename(fmt.Sprintf("./modules/%s", file.Name()), fmt.Sprintf("./modules/%s/%s", file.Name()[:len(file.Name())-7], file.Name()))

			//Unzip the file into the new directory
			extractModule(fmt.Sprintf("./modules/%s", file.Name()[:len(file.Name())-7]),fmt.Sprintf("./modules/%s/%s", file.Name()[:len(file.Name())-7], file.Name()))

			//Remove the .tar.gz file
			os.Remove(fmt.Sprintf("./modules/%s/%s", file.Name()[:len(file.Name())-7], file.Name()))
		}

		return nil
	})

}

func extractModule(folder string, file string) {
	r, err := os.Open(file)
	if err != nil {
		fmt.Println("error")
	}

	uncompressedStream, err := gzip.NewReader(r)
    if err != nil {
        log.Fatal("ExtractTarGz: NewReader failed")
    }

    tarReader := tar.NewReader(uncompressedStream)

    for {
        header, err := tarReader.Next()

        if err == io.EOF {
            break
        }

        if err != nil {
            log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
        }

        switch header.Typeflag {
        case tar.TypeDir:
						//if directory is ./ or ../, skip it
						if header.Name == "./" || header.Name == "../" {
							continue
						}
						fmt.Println(header.Name)
            if err := os.Mkdir(fmt.Sprintf("%s/%s", folder, header.Name), 0755); err != nil {
                log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
            }
        case tar.TypeReg:
            outFile, err := os.Create(fmt.Sprintf("%s/%s", folder, header.Name))
            if err != nil {
                log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
            }
            if _, err := io.Copy(outFile, tarReader); err != nil {
                log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
            }
            outFile.Close()

        default:
            log.Fatal("ExtractTarGz: uknown type:")
      	}

    }
}
