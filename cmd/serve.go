package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/damongolding/eventsforce/internal/utils"
	"github.com/jaschaephraim/lrserver"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

// rootCmd represents the base command when called without any subcommands
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start dev server",
	Long:  "Start dev server",

	Run: func(cmd *cobra.Command, args []string) {
		if err := serve(); err != nil {
			fmt.Print(err)
			return
		}
	}}

func watcher() {

	// Create and start LiveReload server
	lr := lrserver.New("ef", lrserver.DefaultPort)

	print("Watching", config.SrcDir, "👀", "\n")

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		// ticker for debouncing
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) {
					err := build(false)
					if err != nil {
						panic(err)
					}
					print("Rebuilt")
					lr.Reload(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}

			// Wait for 1 sec to pass (if it hasnt already)
			select {
			case <-ticker.C:
				continue
			}
		}
	}()

	// Watch all subdirs
	filepath.Walk(config.SrcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Add a path.to watch
			err = watcher.Add(path)
			if err != nil {
				return err
			}
		}

		return nil

	})

	go lr.ListenAndServe()

	select {}

}

func serve() error {

	if err := build(false); err != nil {
		return err
	}

	go watcher()

	url := fmt.Sprintf("http://localhost:%d", devPort)

	http.Handle("GET /", http.FileServer(http.Dir(config.BuildDir)))
	http.Handle("GET /_assets/", http.FileServer(http.Dir(config.SrcDir)))

	print("Serving on", url)

	err := utils.Openbrowser(url)
	if err != nil {
		return err
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", devPort), nil)
	if err != nil {
		return err
	}

	return nil
}
