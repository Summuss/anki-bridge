package cmd

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
	"log"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(listenFolderCMD)
}

var listenFolderCMD = &cobra.Command{
	Use:  "listen_folder",
	Args: cobra.MaximumNArgs(0),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return listen(args[0])
		} else {
			return listen("D:\\Documents\\voice-records")
		}
	},
}

func listen(folder string) (merr error) {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if merr != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) {
					filename := filepath.Base(event.Name)
					clipboard.Write(clipboard.FmtText, []byte("#FILENAME "+filename))
				}
				log.Println(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)

			}
		}
	}()
	err = watcher.Add(folder)
	if err != nil {
		return err
	}
	<-make(chan struct{})

	return nil
}
