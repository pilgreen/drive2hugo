package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "time"

  "github.com/pilgreen/drive2hugo/auth"
  "google.golang.org/api/drive/v3"
)

type Config struct {
  Folders []Folder
}

type Folder struct {
  Id string
  Path string
}

func loadConfig(file string) (*Config, error) {
  f, err := os.Open(file)
  if err != nil {
    return nil, err
  }
  c := &Config{}
  err = json.NewDecoder(f).Decode(c)
  defer f.Close()
  return c, err
}

func modTime() (string, error) {
  b, err := ioutil.ReadFile("d2h.modified.txt")
  if err != nil {
    return "", err
  }

  return string(b), err
}

func saveModFile() {
  f, err := os.Create("d2h.modified.txt")
  if err != nil {
    log.Fatalf("Unable to create modified file: %v\n", err)
  }

  defer f.Close()
  t, _ := time.Now().MarshalText()
  f.Write(t)
}

func main() {
  cPtr := flag.String("config", "d2h.config.json", "Path to the config file")
  mPtr := flag.Bool("modified", false, "Only pulls files modified since the last run")
  flag.Parse()

  // Get the Service
  srv, err := auth.GetService(*cPtr)
  if err != nil {
    log.Fatalf("Unable to get the Drive Service: %v", err)
  }

  // Load the config
  config, err := loadConfig(*cPtr);
  if err != nil {
    log.Fatalf("Unable to read a config file: %v", err)
  }

  // Main loop
  for _, folder := range config.Folders {
    var q string

    mod, err := modTime()
    if *mPtr == true && err == nil {
      q = fmt.Sprintf("mimeType = 'application/vnd.google-apps.document' and '%s' in parents and modifiedTime > '%s'", folder.Id, mod)
    } else {
      q = fmt.Sprintf("mimeType = 'application/vnd.google-apps.document' and '%s' in parents", folder.Id)
    }

    r, err := srv.Files.List().Q(q).Fields("files(name, id)").Do()
    if err != nil {
      log.Fatalf("Unable to retrieve files: %v", err)
    }

    log.Printf("%d files found in %s\n", len(r.Files), folder.Id)
    if len(r.Files) > 0 {
      os.MkdirAll(folder.Path, 0755)
      queue := make(chan *drive.File, len(r.Files))
      done := make(chan bool)

      // Set up the goroutine to handle the individual files
      go func() {
        for {
          file, more := <-queue
          if more {
            req, err := srv.Files.Export(file.Id, "text/plain").Download()
            if err != nil {
              log.Fatalf("Unable to download file: %v", err)
            }

            filename := file.Id + ".md"
            path := filepath.Join(folder.Path, filename)
            body, _ := ioutil.ReadAll(req.Body)

            err = ioutil.WriteFile(path, body, 0644)
            if err != nil {
              log.Printf("Problem writing %s: %v\n", file.Name, err)
            } else {
              log.Printf("%s written as %s\n", file.Name, path)
            }
          } else {
            log.Printf("all done with folder: %s\n", folder.Id)
            done <- true
            return
          }
        }
      }()


      for _, file := range r.Files {
        queue <- file
      }

      close(queue)
      <- done
    }
  }

  // Save the modified file for the next run
  if *mPtr == true {
    saveModFile()
  }
}
