package cmd

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/summuss/anki-bridge/internal/parser"
	"golang.design/x/clipboard"
	"log"
	"path/filepath"
	"strings"
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
		go listenClipboard()
		folder := ""
		if len(args) > 0 {
			folder = args[0]
		} else {
			folder = "D:\\Documents\\voice-records"
		}
		go listen(folder)
		<-make(chan struct{})
		return nil
	},
}

func listenClipboard() {
	watch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	pre := ""
	for {
		data := <-watch
		txt := string(data)
		if txt != pre && !parser.JPWordPattern.MatchString(txt) {
			txt = strings.Replace(txt, "\n", "", -1)
			txt = strings.Replace(txt, "\r", "", -1)
			clipboard.Write(clipboard.FmtText, []byte(txt))
			pre = txt
		}
	}
}

func listen(folder string) {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()
	err = watcher.Add(folder)
	if err != nil {
		panic(err)
	}
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
}
