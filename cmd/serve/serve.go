package serve

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/damongolding/eventsforce/cmd/build"
	"github.com/damongolding/eventsforce/internal/configuration"
	"github.com/damongolding/eventsforce/internal/utils"

	"github.com/jaschaephraim/lrserver"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	config  configuration.Config
	devPort string
)

func init() {
	config = *configuration.NewConfig()
}

// rootCmd represents the base command when called without any subcommands
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start dev server",
	Long:  "Start dev server",

	Run: func(cmd *cobra.Command, args []string) {

		devPort = cmd.Flag("port").Value.String()

		if err := Serve(); err != nil {
			fmt.Print(err)
			return
		}
	}}

func watcher() {

	// Create and start LiveReload server
	lr := lrserver.New("ef", lrserver.DefaultPort)
	lr.SetStatusLog(nil)
	lr.SetErrorLog(nil)
	lr.SetStatusLog(nil)

	fmt.Println(utils.SectionMessage("Watching", config.SrcDir, "for changes", "👀"))

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
					err := build.Build(false)
					if err != nil {
						panic(err)
					}
					fmt.Println(utils.SectionMessage("Rebuilt"))
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

func Serve() error {

	if err := build.Build(false); err != nil {
		return err
	}

	config.PrintConfig()

	serverSectionTitle, err := utils.OutputSectionStyling("S", "E", "R", "V", "E")
	if err != nil {
		return err
	}

	fmt.Println(serverSectionTitle)

	go watcher()

	url := fmt.Sprintf("http://localhost:%s", devPort)

	http.Handle("GET /", http.FileServer(http.Dir(config.BuildDir)))
	http.Handle("GET /_assets/", http.FileServer(http.Dir(config.SrcDir)))

	fmt.Println(utils.SectionMessage("Serving on", url))

	err = utils.Openbrowser(url)
	if err != nil {
		return err
	}

	err = http.ListenAndServe(fmt.Sprintf(":%s", devPort), nil)
	if err != nil {
		return err
	}

	return nil
}
